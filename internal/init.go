package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller"
	v1 "sorkin_bot/internal/controller/v1"
	"sorkin_bot/internal/domain/service/admin"
	"sorkin_bot/internal/domain/service/referal"
	"sorkin_bot/internal/domain/service/user"
	"sorkin_bot/internal/domain/usercases/create_admin"
	"sorkin_bot/internal/domain/usercases/create_referal"
	"sorkin_bot/internal/domain/usercases/create_user"
	"sorkin_bot/pkg/client/postgres"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService    v1.UserService
	referalService v1.ReferalService
	adminService   admin.AdminService
}

type useCases struct {
	createReferalUseCase referal.CreateReferalUseCase
	createUserUseCase    user.CreateUserUseCase
	createAdminUseCase   admin.CreateAdminUseCase
}

type storages struct {
	adminReadStorage    admin.ReadAdminStorage
	adminWriteStorage   create_admin.WriteRepo
	referalReadStorage  referal.ReadReferalStorage
	referalWriteStorage create_referal.WriteRepo
	userReadStorage     user.ReadUserStorage
	userWriteStorage    create_user.WriteRepo
}

type App struct {
	logger     *slog.Logger
	config     *config.Config
	controller controllers
	services   services
	storages   storages
	useCases   useCases
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
		InitUseCases(ctx).
		InitServices(ctx).
		InitControllers(ctx)

	return app
}

func (app *App) Run() error {
	app.logger.Info("start server")
	return app.server.ListenAndServe()
}
