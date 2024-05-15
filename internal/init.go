package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/clients/gateways/mis_reno"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/adapter"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/bot"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/internal/domain/usecases/bot/save_message_log"
	"sorkin_bot/internal/domain/usecases/user/change_language"
	"sorkin_bot/internal/domain/usecases/user/change_user_status"
	"sorkin_bot/internal/domain/usecases/user/create_user"
	"sorkin_bot/internal/domain/usecases/user/update_user_birth_date"
	"sorkin_bot/internal/domain/usecases/user/update_user_patient_id"
	"sorkin_bot/internal/domain/usecases/user/update_user_phone"
	"sorkin_bot/internal/domain/usecases/user/update_user_third_name"
	"sorkin_bot/internal/storage/read_repo"
	"sorkin_bot/internal/storage/write_repo"
	"sorkin_bot/internal/worker_pool"
	"sorkin_bot/internal/worker_pool/tasks"
	"sorkin_bot/pkg/client/postgres"
	"sorkin_bot/pkg/client/telegram"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService        user.UserService
	appointmentService appointment.AppointmentService
	messageService     message.MessageService
	botService         bot.BotService
}

type useCases struct {
	createUserUserCase         create_user.CreateUserUseCase
	changeLanguageUseCase      change_language.ChangeLanguageUseCase
	changeStatusUseCase        change_user_status.ChangeStatusUseCase
	saveMessageUseCase         save_message_log.SaveMessageLogUseCase
	updateUserPatientIdUseCase update_user_patient_id.UpdateUserPatientIdUseCase
	updateUserPhoneUseCase     update_user_phone.UpdateUserPhoneUseCase
	updateUserThirdNameUseCase update_user_third_name.UpdateUserThirdNameUseCase
	updateUserBirthDateUseCase update_user_birth_date.UpdateUserBirthDateUseCase
}

type storages struct {
	readUserStorage        read_repo.UserStorage
	readTranslationStorage read_repo.TranslationStorage
	readMessageStorage     read_repo.MessageStorage
	writeUserStorage       write_repo.UserStorage
	writeTelegramStorage   write_repo.TelegramMessageStorage
}

type workers struct {
	everyDayWorker worker_pool.Worker
}

type periodicalTasks struct {
	getTranslatedSpeciality tasks.GetTranslatedSpecialityTask
}

type adapters struct {
	appointmentServiceAdapter *adapter.AppointmentServiceAdapter
}

type gateways struct {
	MisRenoGateway mis_reno.MisRenoGateway
}

type App struct {
	logger          *slog.Logger
	config          *config.Config
	controller      controllers
	machine         *state_machine.UserStateMachine
	services        services
	storages        storages
	useCases        useCases
	gateways        gateways
	workers         workers
	adapters        adapters
	periodicalTasks periodicalTasks
	bot             telegram.Bot
	pgxClient       postgres.Client
	server          *http.Server
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
		InitGateways(ctx).
		InitAdapters(ctx).
		InitUseCases(ctx).
		InitServices(ctx).
		InitTelegram(ctx).
		InitTasks(ctx).
		InitWorkers(ctx).
		InitControllers(ctx)

	return app
}

func (app *App) Run() error {
	app.logger.Info("start server")
	go app.workers.everyDayWorker.Run()
	return app.server.ListenAndServe()
}
