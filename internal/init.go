package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/internal/domain/usecases/changeLanguage"
	"sorkin_bot/internal/domain/usecases/create_user"
	"sorkin_bot/internal/storage/read_repo"
	"sorkin_bot/internal/storage/write_repo"
	"sorkin_bot/pkg/client/postgres"
	"sorkin_bot/pkg/client/telegram"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService user.UserService
}

type useCases struct {
	createUserUserCase    create_user.CreateUserUseCase
	changeLanguageUseCase changeLanguage.ChangeLanguageUseCase
}

type storages struct {
	readUserStorage  read_repo.UserStorage
	writeUserStorage write_repo.UserStorage
}

type App struct {
	logger     *slog.Logger
	config     *config.Config
	controller controllers
	ufsm       *state_machine.UserStateMachine
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
		InitFSM(ctx).
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
