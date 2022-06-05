package dialog

import "errors"

type Conversation struct {
	Steps            []Step
	CurrentStepIndex int
	Message          string
	Data             map[string]string
}

func (c *Conversation) GetSteps() []Step {
	return c.Steps
}

func (c *Conversation) GetCurrentStepIndex() int {
	return c.CurrentStepIndex
}

func (c *Conversation) GetMessage() string {
	return c.Message
}

func (c *Conversation) GetData() map[string]string {
	return c.Data
}

func (c *Conversation) AddData(key, value string) {
	c.Data[key] = value
}

func (c *Conversation) GetCurrentStep() Step {
	return c.Steps[c.CurrentStepIndex]
}

func (c *Conversation) Next() error {
	if c.CurrentStepIndex+1 == len(c.Steps) {
		return errors.New("max steps reached")
	}

	c.CurrentStepIndex++

	return nil
}

func (c *Conversation) GetName() string {
	return "task"
}
