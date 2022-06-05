package handlers

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gobuffalo/envy"
	"jira-bot/internal"
	"jira-bot/internal/bot"
	"jira-bot/internal/jira"
	"strconv"
)

type Payload struct {
	Project     string `json:"project,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IssueType   string `json:"issue_type,omitempty"`
	ChatID      string `json:"chat_id,omitempty"`
}

func NewTelegramTaskHandler(bot *tgbotapi.BotAPI) *TelegramTaskHandler {
	jiraLogin := envy.Get("JIRA_LOGIN", "")
	jiraPassword := envy.Get("JIRA_PASSWORD", "")
	jiraHost := envy.Get("JIRA_HOST", "")

	return &TelegramTaskHandler{bot: bot, jira: jira.NewJira(jiraLogin, jiraPassword, jiraHost)}
}

type TelegramTaskHandler struct {
	bot  *tgbotapi.BotAPI
	jira *jira.Jira
}

func (h *TelegramTaskHandler) Handle(message *internal.Message) error {
	var payload *Payload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return err
	}

	chatID, err := strconv.Atoi(payload.ChatID)
	if err != nil {
		return err
	}

	chatIDConverted := int64(chatID)
	jiraAuth, err := h.jira.Authenticate()
	if err != nil {
		if _, err = h.bot.Send(bot.NewErrorMessage(chatIDConverted)); err != nil {
			return err
		}

		return err
	}

	issue, err := h.jira.CreateIssue(jiraAuth, &jira.Issue{Fields: jira.Fields{
		Project: jira.Project{
			ID: payload.Project,
		},
		IssueType: jira.IssueType{
			ID: payload.IssueType,
		},
		Summary:     payload.Title,
		Description: payload.Description,
	}})
	if err != nil {
		if _, err = h.bot.Send(bot.NewErrorMessage(chatIDConverted)); err != nil {
			return err
		}

		return err
	}

	text := fmt.Sprintf("Задача создана: %s/browse/%s", envy.Get("JIRA_HOST", ""), issue.Key)
	msg := tgbotapi.NewMessage(chatIDConverted, text)

	if _, err = h.bot.Send(msg); err != nil {
		return err
	}

	return nil
}
