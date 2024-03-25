package bot

import "log/slog"

type BotService struct {
	readRepo                 ReadMessagesRepo
	administratorHelpUseCase AdministratorHelpUseCase
	// MORE usecases
	logger *slog.Logger
}

func NewBotService(readRepo ReadMessagesRepo, administratorHelpUseCase AdministratorHelpUseCase, logger *slog.Logger) BotService {
	return BotService{
		readRepo:                 readRepo,
		administratorHelpUseCase: administratorHelpUseCase,
		// MORE usecases
		logger: logger,
	}
}

func (bs BotService) AdministratorHelp() {
	// 	TODO create request message_log - controller or adapter mb
	//  TODO get language
	//	TODO get admin message by language
	//  return message
	//	TODO create response message_log may be go send_message() {} - controller or adapter mb
}

func (bs BotService) CancelAppointment() {
	// 	TODO create request message_log - controller or adapter mb

	//  TODO POST TO cancel_appointment
	//  TODO get language
	//	TODO get cancel_appointment message by language and by status
	//  return message

	//	TODO create response message_log may be go send_message() {} - controller or adapter mb
}
