package change_language

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

type ChangeLanguageCommand struct {
	logger      *slog.Logger
	botGateway  botGateway
	tgUser      tg.TgUserDTO
	machine     *state_machine.UserStateMachine
	userService userService
}

func NewChangeLanguageCommand(logger *slog.Logger, botGateway botGateway, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService) ChangeLanguageCommand {
	return ChangeLanguageCommand{
		logger:      logger,
		botGateway:  botGateway,
		tgUser:      tgUser,
		machine:     machine,
		userService: userService,
	}
}

func (c ChangeLanguageCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	userEntity, err := c.userService.GetUser(ctx, c.tgUser.TgID)
	if err != nil {
		return
	}

	c.botGateway.SendChangeLanguageMessage(ctx, userEntity, messageDTO)
	c.machine.SetState(userEntity, state_machine.ChooseLanguage)
}
