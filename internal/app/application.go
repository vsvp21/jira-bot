package app

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gobuffalo/envy"
	"jira-bot/internal"
	"jira-bot/internal/handlers"
	"log"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func NewApplication() *Application {
	chatBot := internal.NewChatBot()
	bot, err := tgbotapi.NewBotAPI(envy.Get("BOT_TOKEN", ""))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = debugMode()

	return &Application{
		bot:            bot,
		chatBot:        chatBot,
		messageHandler: handlers.NewMessageHandler(bot, chatBot),
	}
}

type Application struct {
	bot            *tgbotapi.BotAPI
	chatBot        *internal.ChatBot
	messageHandler *handlers.MessageHandler
}

func (a *Application) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	go a.chatBot.TaskQueue.Listen(handlers.NewTelegramTaskHandler(a.bot))
	go a.runUpdates(cancel)
	go a.clearDialogs()

	select {
	case <-ctx.Done():
		a.shutDown()
	}

	return nil
}

func (a *Application) runUpdates(cancelFunc context.CancelFunc) {
	updateOffset, err := strconv.Atoi(envy.Get("UPDATE_OFFSET", "0"))
	if err != nil {
		log.Println(err)
		cancelFunc()
	}

	updateTimeout, err := strconv.Atoi(envy.Get("UPDATE_TIMEOUT", "60"))
	if err != nil {
		log.Println(err)
		cancelFunc()
	}

	u := tgbotapi.NewUpdate(updateOffset)
	u.Timeout = updateTimeout
	for update := range a.bot.GetUpdatesChan(u) {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				a.messageHandler.HandleCommand(cancelFunc, update)
			} else {
				a.messageHandler.Handle(cancelFunc, update)
			}
		}
	}
}

func (a *Application) clearDialogs() {
	for {
		a.chatBot.ClearDialogs()
		fmt.Println("Dialogs cleared")
		time.Sleep(time.Hour)
	}
}

func (a *Application) shutDown() {
	shutDownTimeout, err := strconv.Atoi(envy.Get("SHUTDOWN_TIMEOUT", "10"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Stopping Application gracefully")

	a.chatBot.ClearDialogs()
	if !a.chatBot.TaskQueue.IsEmpty() {
		<-time.After(time.Second * time.Duration(shutDownTimeout))
	}

	a.chatBot.TaskQueue.Close()

	fmt.Println("Exited")
}

func debugMode() bool {
	if envy.Get("APP_DEBUG", "true") == "true" {
		return true
	}

	return false
}
