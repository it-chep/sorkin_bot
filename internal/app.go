package internal

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sorkin_bot/internal/clients/gateways/mis_reno"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/internal/domain/usecases/bot/changeLanguage"
	"sorkin_bot/internal/domain/usecases/bot/save_message_log"
	"sorkin_bot/internal/domain/usecases/user/create_user"
	"sorkin_bot/internal/storage/read_repo"
	"sorkin_bot/internal/storage/write_repo"
	"sorkin_bot/pkg/client/postgres"
	"sorkin_bot/pkg/client/telegram"
	"time"
)

func (app *App) InitLogger(ctx context.Context) *App {
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return app
}

func (app *App) InitPgxConn(ctx context.Context) *App {
	client, err := postgres.NewClient(ctx, app.config.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}
	app.pgxClient = client
	app.logger.Info("init pgxclient", app.pgxClient)
	return app
}

func (app *App) InitStorage(ctx context.Context) *App {
	app.storages.writeUserStorage = write_repo.NewUserStorage(app.pgxClient, app.logger)
	app.storages.readUserStorage = read_repo.NewUserStorage(app.pgxClient, app.logger)
	app.storages.readTranslationStorage = read_repo.NewTranslationRepo(app.pgxClient, app.logger)
	app.storages.readMessageStorage = read_repo.NewReadMessageStorage(app.pgxClient, app.logger)
	app.storages.writeTelegramStorage = write_repo.NewTelegramMessageStorage(app.pgxClient, app.logger)
	return app
}

func (app *App) InitGateways(ctx context.Context) *App {
	app.gateways.MisRenoGateway = mis_reno.NewMisRenoGateway(app.logger, http.Client{})
	return app
}

func (app *App) InitUseCases(ctx context.Context) *App {
	app.useCases.createUserUserCase = create_user.NewCreateUserUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.changeLanguageUseCase = changeLanguage.NewChangeLanguageUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.saveMessageUseCase = save_message_log.NewSaveMessageLogUseCase(app.storages.writeTelegramStorage, app.logger)
	return app
}

func (app *App) InitServices(ctx context.Context) *App {
	app.services.userService = user.NewUserService(
		app.useCases.createUserUserCase,
		app.useCases.changeLanguageUseCase,
		app.useCases.changeStatusUseCase,
		app.storages.readUserStorage,
		app.logger,
	)
	// todo исправить
	app.services.appointmentService = appointment.NewAppointmentService(
		&app.gateways.MisRenoGateway,
		app.storages.readTranslationStorage,
		app.logger,
	)
	app.services.messageService = message.NewMessageService(
		app.useCases.saveMessageUseCase,
		app.storages.readMessageStorage,
		app.logger,
	)
	return app

}

func (app *App) InitMachine(ctx context.Context) *App {
	app.machine = state_machine.NewUserStateMachine(app.services.userService)
	return app
}

func (app *App) InitTelegram(ctx context.Context) *App {
	app.bot = *telegram.NewTelegramBot(*app.config)
	return app
}

func (app *App) InitControllers(ctx context.Context) *App {
	app.controller.telegramWebhookController = controller.NewRestController(*app.config, app.logger, app.bot, app.machine, app.services.userService, app.services.appointmentService, app.services.messageService)
	app.controller.telegramWebhookController.InitController(ctx)

	app.server = &http.Server{
		Addr:         app.config.HTTPServer.Address,
		Handler:      app.controller.telegramWebhookController,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return app
}
