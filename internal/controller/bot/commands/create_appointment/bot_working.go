package create_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/bot/bot_interfaces"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type CreateAppointmentCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        bot_interfaces.UserService
	machine            *state_machine.UserStateMachine
	appointmentService bot_interfaces.AppointmentService
	messageService     bot_interfaces.MessageService
}

func NewCreateAppointmentCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService bot_interfaces.UserService, machine *state_machine.UserStateMachine, appointmentService bot_interfaces.AppointmentService, messageService bot_interfaces.MessageService,
) CreateAppointmentCommand {
	return CreateAppointmentCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
	}
}

func (c CreateAppointmentCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	c.appointmentService.GetSchedules(ctx, 0)

	//todo возможно добавить сообщение, что я загружаю ваши записи, пожалуйста подождите
	//msg = tgbotapi.NewMessage(c.tgUser.TgID)
	//_, _ = c.bot.Bot.Send(msg)

	// todo это для создания записи на прием
	//if c.appointmentService.GetPatient(ctx, userEntity) {
	//} else {
	//	c.appointmentService.CreatePatient(ctx, userEntity)
	//}
	// todo это для создания записи на прием

	////todo докрутить логику со специальностями
	//switch userEntity.GetState() {
	//case "":
	//	err, specialities := c.appointmentService.mis.GetSpecialities(ctx)
	//	if err != nil {
	//		return
	//	}
	//	messageText, err := c.messageService.GetMessage(ctx, userEntity, "Choose speciality")
	//	msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
	//	if err != nil {
	//		_, _ = c.bot.Bot.Send(msg)
	//		return
	//	}
	//	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	//	translatedSpecialities, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities)
	//	if err != nil {
	//		return
	//	}
	//	for specialityId, translatedSpeciality := range translatedSpecialities {
	//		btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", specialityId))
	//		row := tgbotapi.NewInlineKeyboardRow(btn)
	//		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	//	}
	//	msg.ReplyMarkup = keyboard
	//}
	//todo докрутить логику со специальностями

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseSpeciality)

	sentMessage, err := c.bot.Bot.Send(msg)
	// todo мб вынести в отдельный метод
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
	messageDTO.MessageID = int64(sentMessage.MessageID)
	messageDTO.Text = sentMessage.Text

	go func() {
		err := c.messageService.SaveMessageLog(context.TODO(), messageDTO)
		if err != nil {
			return
		}
	}()
}
