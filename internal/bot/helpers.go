package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewErrorMessage(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Произошла ошибка, попробуйте позже")
}

func NewProjectReplyMessage(chatID int64, message string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("express"),
			tgbotapi.NewKeyboardButton("cube"),
		),
	)

	return msg
}

func NewIssueTypeReplyMessage(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Укажите тип задачи")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("task"),
			tgbotapi.NewKeyboardButton("bug"),
		),
	)

	return msg
}

func NewTitleReplyMessage(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Укажите заголовок задачи")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	return msg
}

func NewDescriptionReplyMessage(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Укажите описание задачи")
}
