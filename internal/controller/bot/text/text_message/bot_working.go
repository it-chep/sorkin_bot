package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/pkg/client/telegram"
)

type TextBotMessage struct {
	logger             *slog.Logger
	bot                telegram.Bot
	botGateway         botGateway
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        userService
	messageService     messageService
	appointmentService appointmentService
}

func NewTextBotMessage(logger *slog.Logger, bot telegram.Bot, botGateway botGateway, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService, messageService messageService, appointmentService appointmentService) TextBotMessage {
	return TextBotMessage{
		logger:             logger,
		bot:                bot,
		botGateway:         botGateway,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		messageService:     messageService,
		appointmentService: appointmentService,
	}
}

func (c TextBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	var _ tgbotapi.MessageConfig
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)

	if userEntity.GetPatientId() == nil {
		switch *userEntity.GetState() {
		case state_machine.GetPhone:
			c.getPhone(ctx, userEntity, messageDTO)
		case state_machine.GetName:
			c.getName(ctx, userEntity, messageDTO)
		case state_machine.GetBirthDate:
			c.getBirthDate(ctx, userEntity, messageDTO)
		}
	}
}

func (c TextBotMessage) getBirthDate(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig

	user, result, err := c.userService.UpdateBirthDate(ctx, c.tgUser, messageDTO.Text)
	if err != nil {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
		c.bot.SendMessage(msg, messageDTO)
		return
	}
	if !result {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid birth date")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	user.SetBirthDate(messageDTO.Text)
	draftAppointment, err := c.appointmentService.GetDraftAppointment(ctx, c.tgUser.TgID)
	if err != nil {
		return
	}

	c.botGateway.SendConfirmAppointmentMessage(ctx, user, messageDTO, *draftAppointment.GetDoctorId())
	c.appointmentService.CreatePatient(ctx, user)
	go c.machine.SetState(user, state_machine.CreateAppointment)
}

func (c TextBotMessage) getPhone(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	var phone string
	if messageDTO.Contact != nil {
		phone = messageDTO.Contact.PhoneNumber
	} else {
		phone = messageDTO.Text
	}

	_, result, err := c.userService.UpdatePhone(ctx, c.tgUser, phone)
	if err != nil {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	if !result {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid phone")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	messageText, _ := c.messageService.GetMessage(ctx, user, "enter name")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
	c.bot.SendMessage(msg, messageDTO)
	go c.machine.SetState(user, state_machine.GetName)
}

func (c TextBotMessage) getName(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig

	_, result, err := c.userService.UpdateFullName(ctx, c.tgUser, messageDTO.Text)
	if err != nil {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	if !result {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid name")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	messageText, _ := c.messageService.GetMessage(ctx, user, "enter birthdate")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
	c.bot.SendMessage(msg, messageDTO)
	go c.machine.SetState(user, state_machine.GetBirthDate)
}
