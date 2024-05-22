package bot_gateway

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/clients/bot_gateway/keyboards"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/pkg/client/telegram"
)

type BotGateway struct {
	logger             *slog.Logger
	bot                telegram.Bot
	messageService     messageService
	appointmentService appointmentService
	keyboard           keyboardsInterface
}

func NewBotGateway(
	logger *slog.Logger,
	bot telegram.Bot,
	messageService messageService,
	appointmentService appointmentService,
) BotGateway {
	keyboard := keyboards.NewKeyboards(
		logger, messageService, appointmentService,
	)
	return BotGateway{
		logger:             logger,
		bot:                bot,
		messageService:     messageService,
		appointmentService: appointmentService,
		keyboard:           keyboard,
	}
}

func (bg BotGateway) CreateMessageLog(sentMessage tgbotapi.Message, messageDTO tg.MessageDTO) {
	messageDTO.MessageID = int64(sentMessage.MessageID)
	messageDTO.Text = sentMessage.Text
	go func() {
		err := bg.messageService.SaveMessageLog(context.Background(), messageDTO)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendChooseSpecialityMessage(
	ctx context.Context,
	idToDelete int,
	translatedSpecialities map[int]string,
	user entity.User,
	messageDTO tg.MessageDTO,
) {

	bg.bot.RemoveMessage(user.GetTgId(), idToDelete)
	msgText, keyboard := bg.keyboard.ConfigureGetSpecialityMessage(ctx, user, translatedSpecialities, 0)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendGetDoctorsMessage(
	ctx context.Context,
	user entity.User,
	messageDTO tg.MessageDTO,
	doctors map[int]string,
	offset int,
) {

	msgText, keyboard := bg.keyboard.ConfigureGetDoctorMessage(ctx, user, doctors, offset)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)

	if keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	bg.bot.RemoveMessage(user.GetTgId(), int(messageDTO.MessageID))
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendSchedulesMessage(
	ctx context.Context,
	userEntity entity.User,
	messageDTO tg.MessageDTO,
	schedules []appointment.Schedule,
	offset int,
) {
	msgText, keyboard := bg.keyboard.ConfigureGetScheduleMessage(ctx, userEntity, schedules, offset)
	msg := tgbotapi.NewMessage(userEntity.GetTgId(), msgText)
	if keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendSpecialityMessage(
	ctx context.Context,
	userEntity entity.User,
	messageDTO tg.MessageDTO,
	specialities map[int]string,
	offset int,
) {
	msgText, keyboard := bg.keyboard.ConfigureGetSpecialityMessage(ctx, userEntity, specialities, offset)
	msg := tgbotapi.NewMessage(userEntity.GetTgId(), msgText)

	if keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}
