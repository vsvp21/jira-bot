package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gobuffalo/envy"
	"jira-bot/internal"
	"jira-bot/internal/dialog"
	"log"
)

const (
	StartCommand     = "start"
	PingCommand      = "ping"
	TaskCommand      = "task"
	TaskConversation = "task"
)

var projects = map[string]string{
	"express": envy.Get("EXPRESS_ID", ""),
	"cube":    envy.Get("CUBE_ID", ""),
}

var issueTypes = map[string]string{
	"task": envy.Get("TASK_ID", ""),
	"bug":  envy.Get("BUG_ID", ""),
}

func NewMessageHandler(bot *tgbotapi.BotAPI, chatBot *internal.ChatBot) *MessageHandler {
	return &MessageHandler{Bot: bot, ChatBot: chatBot}
}

type MessageHandler struct {
	Bot     *tgbotapi.BotAPI
	ChatBot *internal.ChatBot
}

func (h *MessageHandler) HandleCommand(cancelFunc context.CancelFunc, update tgbotapi.Update) {
	switch update.Message.Command() {
	case StartCommand:
		if err := handleStart(update, h.Bot); err != nil {
			log.Println(err)
			cancelFunc()
		}
	case PingCommand:
		if err := handlePing(update, h.Bot); err != nil {
			log.Println(err)
			cancelFunc()
		}
	case TaskCommand:
		d := h.ChatBot.NewDialog(update.FromChat().ID, dialog.NewTask())
		if err := handleTask(update, h.Bot, d); err != nil {
			log.Println(err)
			cancelFunc()
		}
	default:
		if err := handleUnknown(update, h.Bot); err != nil {
			log.Println(err)
			cancelFunc()
		}
	}
}

func (h *MessageHandler) Handle(cancelFunc context.CancelFunc, update tgbotapi.Update) {
	d, _ := h.ChatBot.GetActiveDialog(update.FromChat().ID)
	step := d.GetEntryPointStep()

	switch step.GetName() {
	case TaskConversation:
		if err := handleTaskConversation(update, h.Bot, h.ChatBot, step); err != nil {
			log.Println(err)
			cancelFunc()
		}
	default:
		if err := handleUnknownConversation(update, h.Bot); err != nil {
			log.Println(err)
			cancelFunc()
		}
	}
}
