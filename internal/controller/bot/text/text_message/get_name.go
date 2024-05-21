package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"strings"
)

func (c TextBotMessage) getName(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig

	if c.validateNameMessage(messageDTO.Text) {
		_, err := c.userService.UpdateFullName(ctx, c.tgUser, messageDTO.Text)
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
	} else {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid name")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}
	messageText, _ := c.messageService.GetMessage(ctx, user, "enter birthdate")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
	c.bot.SendMessage(msg, messageDTO)
	c.machine.SetState(user, *user.GetState(), state_machine.GetBirthDate)
}

func (c TextBotMessage) validateNameMessage(name string) (valid bool) {
	nameItems := strings.Split(name, " ")

	validNameItemsLength := len(nameItems) == 3
	if validNameItemsLength {
		return true
	}
	return false
}
