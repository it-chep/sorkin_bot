package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (c *CallbackBotMessage) mainMenu(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {

}
