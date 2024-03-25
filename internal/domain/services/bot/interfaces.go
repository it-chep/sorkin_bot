package bot

type ReadMessagesRepo interface {
	GetMessageByCondition()
}

type AdministratorHelpUseCase interface {
	Execute()
}

type CancelAppointmentUseCase interface {
	Execute()
}
