package internal

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
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
	return app
}

func (app *App) InitUseCases(ctx context.Context) *App {
	return app

}

func (app *App) InitServices(ctx context.Context) *App {
	return app

}

func (app *App) InitFSM(ctx context.Context) *App {
	app.ufsm = state_machine.NewUserStateMachine("start")
	return app
}

func (app *App) InitTelegram(ctx context.Context) *App {
	app.bot = *telegram.NewTelegramBot(*app.config)
	return app
}

func (app *App) InitControllers(ctx context.Context) *App {
	app.controller.telegramWebhookController = controller.NewRestController(*app.config, app.logger, app.bot, app.ufsm)
	app.controller.telegramWebhookController.InitController(ctx)

	app.server = &http.Server{
		Addr:         app.config.HTTPServer.Address,
		Handler:      app.controller.telegramWebhookController,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return app
}
