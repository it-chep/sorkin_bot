package bot_gateway

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (bg BotGateway) SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureMainMenuMessage(ctx, user)

	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)

	err := bg.messageService.SaveMessageLog(ctx, sentMessage)
	bg.logger.Error(fmt.Sprintf("%s", err))
}

func (bg BotGateway) SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureChangeLanguageMessage(ctx, user)

	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendGetPhoneMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureGetPhoneMessage(ctx, user)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}
