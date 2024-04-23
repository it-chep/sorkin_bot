package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/internal/domain/usecases/bot /changeLanguage"
	"sorkin_bot/internal/domain/usecases/bot /save_message_log"
	"sorkin_bot/internal/domain/usecases/user/create_user"
	"sorkin_bot/internal/storage/read_repo"
	"sorkin_bot/internal/storage/write_repo"
	"sorkin_bot/pkg/client/postgres"
	"sorkin_bot/pkg/client/telegram"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService    user.UserService
	messageService message.MessageService
}

type useCases struct {
	createUserUserCase    create_user.CreateUserUseCase
	changeLanguageUseCase changeLanguage.ChangeLanguageUseCase
	saveMessageUseCase    save_message_log.SaveMessageLogUseCase
}

type storages struct {
	readUserStorage      read_repo.UserStorage
	writeUserStorage     write_repo.UserStorage
	writeTelegramStorage write_repo.TelegramMessageStorage
}

type App struct {
	logger     *slog.Logger
	config     *config.Config
	controller controllers
	machine    *state_machine.UserStateMachine
	services   services
	storages   storages
	useCases   useCases
	bot        telegram.Bot
	pgxClient  postgres.Client
	server     *http.Server
}

func NewApp(ctx context.Context) *App {
	cfg := config.NewConfig()

	app := &App{
		config: cfg,
	}

	app.InitLogger(ctx).
		InitPgxConn(ctx).
		InitStorage(ctx).
		InitMachine(ctx).
		InitUseCases(ctx).
		InitServices(ctx).
		InitTelegram(ctx).
		InitControllers(ctx)

	return app
}

func (app *App) Run() error {
	app.logger.Info("start server")
	return app.server.ListenAndServe()
}
