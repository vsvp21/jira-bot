package dialog

func NewTask() *Conversation {
	return &Conversation{
		Steps: []Step{
			&Conversation{
				Message: "Укажите название проекта",
			},
			&Conversation{
				Message: "Укажите тип задачи",
			},
			&Conversation{
				Message: "Укажите заголовок задачи",
			},
			&Conversation{
				Message: "Укажите описание задачи",
			},
		},
		CurrentStepIndex: 0,
		Data:             make(map[string]string),
	}
}
