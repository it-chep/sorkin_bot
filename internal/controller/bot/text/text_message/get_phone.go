package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
)

func (c TextBotMessage) GetPhone(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	var phone string
	if messageDTO.Contact.PhoneNumber != "" {
		phone = messageDTO.Contact.PhoneNumber
	} else {
		phone = messageDTO.Text
	}
	if c.validatePhoneMessage(phone) {
		_, err := c.userService.UpdatePhone(ctx, c.tgUser, messageDTO.Text)
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
	} else {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid phone")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	c.machine.SetState(user, state_machine.GetPhone, state_machine.GetName)

}

func (c TextBotMessage) validatePhoneMessage(phone string) (valid bool) {
	pattern := ""
	//todo add phone pattern
	matchString, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	if matchString {
		return true
	}
	return false
}
