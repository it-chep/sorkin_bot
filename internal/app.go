package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/clients/bot_gateway"
	"sorkin_bot/internal/clients/gateways/mis_reno"
	"sorkin_bot/internal/clients/sms_gateway/wau_sms"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/adapter"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/notification"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/internal/domain/usecases/appointment/clean_draft_appointment"
	"sorkin_bot/internal/domain/usecases/appointment/create_draft_appointment"
	"sorkin_bot/internal/domain/usecases/appointment/fast_update_draft_appointment_use_case"
	"sorkin_bot/internal/domain/usecases/appointment/update_appointment_date"
	"sorkin_bot/internal/domain/usecases/appointment/update_appointment_status"
	"sorkin_bot/internal/domain/usecases/appointment/update_int_appointment_field"
	"sorkin_bot/internal/domain/usecases/appointment/update_str_appointment_field"
	"sorkin_bot/internal/domain/usecases/bot/save_message_log"
	"sorkin_bot/internal/domain/usecases/user/change_language"
	"sorkin_bot/internal/domain/usecases/user/change_user_status"
	"sorkin_bot/internal/domain/usecases/user/create_user"
	"sorkin_bot/internal/domain/usecases/user/update_home_address"
	"sorkin_bot/internal/domain/usecases/user/update_user_birth_date"
	"sorkin_bot/internal/domain/usecases/user/update_user_full_name"
	"sorkin_bot/internal/domain/usecases/user/update_user_patient_id"
	"sorkin_bot/internal/domain/usecases/user/update_user_phone"
	"sorkin_bot/internal/storage/read_repo"
	"sorkin_bot/internal/storage/write_repo"
	"sorkin_bot/internal/worker_pool"
	"sorkin_bot/internal/worker_pool/tasks/check_speciallity_translation_task"
	"sorkin_bot/internal/worker_pool/tasks/check_support_calls"
	notify_appointment "sorkin_bot/internal/worker_pool/tasks/notify_appointment"
	"sorkin_bot/pkg/client/postgres"
	"sorkin_bot/pkg/client/telegram"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService         user.UserService
	appointmentService  appointment.AppointmentService
	messageService      message.MessageService
	notificationService *notification.Service
}

type useCases struct {
	createUserUserCase                    create_user.CreateUserUseCase
	changeLanguageUseCase                 change_language.ChangeLanguageUseCase
	changeStatusUseCase                   change_user_status.ChangeStatusUseCase
	saveMessageUseCase                    save_message_log.SaveMessageLogUseCase
	updateUserPatientIdUseCase            update_user_patient_id.UpdateUserPatientIdUseCase
	updateUserPhoneUseCase                update_user_phone.UpdateUserPhoneUseCase
	updateUserHomeAddressUseCase          update_home_address.UpdateUserHomeAddressUseCase
	updateUserThirdNameUseCase            update_user_full_name.UpdateUpdateFullNameUseCase
	updateUserBirthDateUseCase            update_user_birth_date.UpdateUserBirthDateUseCase
	createDraftAppointmentUseCase         create_draft_appointment.CreateDraftAppointmentUseCase
	updateDraftAppointmentStatusUseCase   update_appointment_status.UpdateAppointmentStatusUseCase
	updateDraftAppointmentIntFieldUseCase update_int_appointment_field.UpdateIntAppointmentFieldUseCase
	updateDraftAppointmentStrFieldUseCase update_str_appointment_field.UpdateStrAppointmentFieldUseCase
	updateDraftAppointmentDateUseCase     update_appointment_date.UpdateAppointmentDate
	cleanDraftAppointmentUseCase          clean_draft_appointment.CleanDraftAppointmentUseCase
	fastUpdateDraftAppointmentUseCase     fast_update_draft_appointment_use_case.FastUpdateDraftAppointmentUseCase
}

type storages struct {
	readUserStorage              read_repo.UserStorage
	readTranslationStorage       read_repo.TranslationStorage
	readMessageStorage           read_repo.MessageStorage
	readDraftAppointmentStorage  read_repo.AppointmentStorage
	readLogsStorage              read_repo.TelegramMessageStorage
	writeUserStorage             write_repo.UserStorage
	writeTelegramStorage         write_repo.TelegramMessageStorage
	writeDraftAppointmentStorage write_repo.AppointmentStorage
}

type periodicalTasks struct {
	getTranslatedSpeciality check_speciallity_translation_task.Task
	checkSupportTask        check_support_calls.Task
	notifyAppointmentTask   notify_appointment.Task
}

type adapters struct {
	appointmentServiceAdapter *adapter.AppointmentServiceAdapter
}

type gateways struct {
	MisRenoGateway mis_reno.MisRenoGateway
	WAUSMSGateway  *wau_sms.Sender
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
	adapters        adapters
	botGateway      bot_gateway.BotGateway
	periodicalTasks periodicalTasks
	workerPool      worker_pool.WorkerPool
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
		InitGateways(ctx).
		InitAdapters(ctx).
		InitUseCases(ctx).
		InitServices(ctx).
		InitMachine(ctx).
		InitTelegram(ctx).
		InitBotGateway(ctx).
		InitTasks(ctx).
		InitWorkers(ctx).
		InitControllers(ctx)

	return app
}

func (app *App) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			app.logger.Error("application recovered from panic", slog.Any("error", r))
		}
	}()

	app.logger.Info("start server")
	go app.workerPool.Run(ctx)
	return app.server.ListenAndServe()
}
