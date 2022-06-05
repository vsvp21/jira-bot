package dialog

type Step interface {
	GetSteps() []Step
	GetCurrentStepIndex() int
	GetMessage() string
	AddData(string, string)
	GetData() map[string]string
	Next() error
	GetCurrentStep() Step
	GetName() string
}

type Dialog struct {
	EntryPointStep Step
}

func (s *Dialog) GetEntryPointStep() Step {
	return s.EntryPointStep
}
