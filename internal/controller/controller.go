package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	botapi "sorkin_bot/internal/controller/bot"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"

	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
)

type RestController struct {
	router           *gin.Engine
	cfg              config.Config
	logger           *slog.Logger
	botApiController botapi.TelegramWebhookController
	userService      user.UserService
	messageService   message.MessageService
}

func NewRestController(cfg config.Config, logger *slog.Logger, bot telegram.Bot, machine *state_machine.UserStateMachine, userService user.UserService, messageService message.MessageService) *RestController {
	router := gin.New()
	router.Use(gin.Recovery())

	botApiController := botapi.NewTelegramWebhookController(cfg, logger, bot, machine, userService, messageService)

	return &RestController{
		router:           router,
		cfg:              cfg,
		logger:           logger,
		botApiController: botApiController,
	}
}

func (r RestController) InitController(ctx context.Context) {
	r.router.POST("/"+r.cfg.Bot.Token+"/", r.botApiController.BotWebhookHandler)
}

func (r RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
