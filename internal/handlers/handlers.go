package handlers

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"jira-bot/internal"
	bot2 "jira-bot/internal/bot"
	"jira-bot/internal/dialog"
	"log"
	"strconv"
	"sync"
)

func handleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Привет! "+update.SentFrom().FirstName)

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func handlePing(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Pong")
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func handleTask(update tgbotapi.Update, bot *tgbotapi.BotAPI, dialog *dialog.Dialog) error {
	msg := bot2.NewProjectReplyMessage(update.FromChat().ID, dialog.GetEntryPointStep().GetCurrentStep().GetMessage())
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func handleUnknown(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	_, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "К сожалению я не знаю такой команды("))
	if err != nil {
		return err
	}

	return nil
}

func handleTaskConversation(update tgbotapi.Update, bot *tgbotapi.BotAPI, chatBot *internal.ChatBot, step dialog.Step) error {
	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()

	if step.GetCurrentStepIndex() == 0 {
		step.AddData("project", projects[update.Message.Text])

		if err := step.Next(); err != nil {
			log.Fatal(err)
		}

		if _, err := bot.Send(bot2.NewIssueTypeReplyMessage(update.FromChat().ID)); err != nil {
			return err
		}

		return nil
	}

	if step.GetCurrentStepIndex() == 1 {
		step.AddData("issue_type", issueTypes[update.Message.Text])

		if _, err := bot.Send(bot2.NewTitleReplyMessage(update.FromChat().ID)); err != nil {
			return err
		}

		if err := step.Next(); err != nil {
			log.Fatal(err)
		}

		return nil
	}

	if step.GetCurrentStepIndex() == 2 {
		step.AddData("title", update.Message.Text)

		_, err := bot.Send(bot2.NewDescriptionReplyMessage(update.FromChat().ID))
		if err != nil {
			return err
		}

		if err = step.Next(); err != nil {
			log.Fatal(err)
		}

		return nil
	}

	if step.GetCurrentStepIndex() == 3 {
		step.AddData("description", update.Message.Text)
		step.AddData("chat_id", strconv.Itoa(int(update.FromChat().ID)))

		data, err := json.Marshal(step.GetData())
		if err != nil {
			log.Fatal(err)
		}

		chatBot.TaskQueue.Send(&internal.Message{Payload: data})

		_, err = bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Задача создается"))
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func handleUnknownConversation(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "К сожалению я не знаю ответа(")); err != nil {
		return err
	}

	return nil
}
