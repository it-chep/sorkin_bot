package bot_gateway

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/message"
)

func (bg BotGateway) SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureMainMenuMessage(ctx, user)

	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureChangeLanguageMessage(ctx, user)

	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendGetPhoneMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureGetPhoneMessage(ctx, user)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendError(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msg := tgbotapi.NewMessage(user.GetTgId(), message.ServerError)
	bg.bot.SendMessage(msg, messageDTO)
}
