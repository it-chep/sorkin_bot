package controller

import (
	"github.com/gin-gonic/gin"
	botapi "sorkin_bot/internal/controller/bot"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/middleware"
	"sorkin_bot/pkg/client/telegram"

	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
)

type RestController struct {
	router             *gin.Engine
	cfg                config.Config
	logger             *slog.Logger
	botApiController   botapi.TelegramWebhookController
	userService        UserService
	appointmentService AppointmentService
	messageService     MessageService
}

func NewRestController(cfg config.Config, logger *slog.Logger, bot telegram.Bot, machine *state_machine.UserStateMachine, userService UserService, appointmentService AppointmentService, messageService MessageService, botService BotService) *RestController {
	router := gin.New()
	botMiddleware := middleware.NewMessageLogMiddleware(messageService)
	sentryMiddleware := middleware.NewSentryMiddleware()
	tgAdminMiddleware := middleware.NewTgAdminWarningMiddleware()
	router.Use(gin.Recovery(), botMiddleware.ProcessRequest, sentryMiddleware.ProcessRequest, tgAdminMiddleware.ProcessRequest)

	botApiController := botapi.NewTelegramWebhookController(cfg, logger, bot, machine, userService, appointmentService, messageService, botService)

	return &RestController{
		router:           router,
		cfg:              cfg,
		logger:           logger,
		botApiController: botApiController,
	}
}

func (r RestController) InitController() {
	r.router.POST("/"+r.cfg.Bot.Token+"/", r.botApiController.BotWebhookHandler)
}

func (r RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
