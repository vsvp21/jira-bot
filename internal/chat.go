package internal

import (
	"errors"
	"jira-bot/internal/dialog"
	"sync"
)

var ErrDialogNotFound = errors.New("dialog not found")

func NewChatBot() *ChatBot {
	return &ChatBot{
		Dialogs:   make(map[int64]*dialog.Dialog),
		TaskQueue: NewQueue(),
	}
}

type ChatBot struct {
	Dialogs   map[int64]*dialog.Dialog
	TaskQueue *Queue
	mu        sync.Mutex
}

func (c *ChatBot) NewDialog(userID int64, step dialog.Step) *dialog.Dialog {
	c.mu.Lock()
	defer c.mu.Unlock()

	d := &dialog.Dialog{EntryPointStep: step}
	c.Dialogs[userID] = d

	return d
}

func (c *ChatBot) GetActiveDialog(userID int64) (*dialog.Dialog, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.Dialogs[userID]
	if !ok {
		return nil, ErrDialogNotFound
	}

	return d, nil
}

func (c *ChatBot) ClearDialogs() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Dialogs = make(map[int64]*dialog.Dialog)
}
