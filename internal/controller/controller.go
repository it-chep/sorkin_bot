package controller

import (
	"github.com/gin-gonic/gin"
	"sorkin_bot/internal/controller/api"
	botapi "sorkin_bot/internal/controller/bot"
	"sorkin_bot/internal/domain/entity/user/state_machine"
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
	userService        userService
	appointmentService appointmentService
	messageService     messageService
	botGateway         botGateway
	apiController      apiController
}

func NewRestController(
	cfg config.Config,
	logger *slog.Logger,
	bot telegram.Bot,
	machine *state_machine.UserStateMachine,
	userService userService,
	appointmentService appointmentService,
	messageService messageService,
	notificationService notificationService,
	botGateway botGateway,
) *RestController {
	router := gin.New()
	//languageMiddleware
	//botMiddleware := middleware.NewMessageLogMiddleware(messageService)
	//sentryMiddleware := middleware.NewSentryMiddleware()
	//tgAdminMiddleware := middleware.NewTgAdminWarningMiddleware()
	router.Use(gin.Recovery())

	botApiController := botapi.NewTelegramWebhookController(
		cfg, logger, bot, machine, userService, appointmentService, messageService, botGateway,
	)

	apiWebhookController := api.NewController(
		logger,
		notificationService,
	)

	return &RestController{
		router:           router,
		cfg:              cfg,
		logger:           logger,
		botApiController: botApiController,
		apiController:    apiWebhookController,
	}
}

func (r RestController) InitController() {
	webhookGroup := r.router.Group("/webhook_event/")
	webhookGroup.POST("cancel_appointment/", r.apiController.CancelAppointmentWebhook)
	webhookGroup.POST("create_appointment/", r.apiController.CreateAppointmentWebhook)

	r.router.POST("/"+r.cfg.Bot.Token+"/", r.botApiController.BotWebhookHandler)
}

func (r RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
