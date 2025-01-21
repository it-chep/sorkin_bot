package internal

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sorkin_bot/internal/clients/bot_gateway"
	"sorkin_bot/internal/clients/gateways/mis_reno"
	"sorkin_bot/internal/clients/sms_gateway/wau_sms"
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
	"time"
)

//const (
//	RU    = "RU"
//	PT_BR = "PT_BR"
//	EN    = "EN"
//)

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
	app.storages.readLogsStorage = read_repo.NewTelegramMessageStorage(app.pgxClient, app.logger)
	app.storages.writeTelegramStorage = write_repo.NewTelegramMessageStorage(app.pgxClient, app.logger)
	app.storages.readDraftAppointmentStorage = read_repo.NewAppointmentStorage(app.pgxClient, app.logger)
	app.storages.writeDraftAppointmentStorage = write_repo.NewAppointmentStorage(app.pgxClient, app.logger)
	return app
}

func (app *App) InitGateways(ctx context.Context) *App {
	app.gateways.MisRenoGateway = mis_reno.NewMisRenoGateway(app.logger, http.Client{Timeout: time.Second * 10})
	app.gateways.WAUSMSGateway = wau_sms.NewSender(app.logger, *app.config)
	return app
}

func (app *App) InitTasks(ctx context.Context) *App {
	app.periodicalTasks.getTranslatedSpeciality = check_speciallity_translation_task.NewTask(
		&app.services.appointmentService,
		app.services.userService,
		app.logger,
		app.bot,
	)

	app.periodicalTasks.checkSupportTask = check_support_calls.NewTask(
		app.logger,
		app.bot,
		app.services.messageService,
		app.services.userService,
	)

	app.periodicalTasks.notifyAppointmentTask = notify_appointment.NewTask(
		&app.services.appointmentService,
		app.services.notificationService,
		app.logger,
		app.bot,
	)
	return app
}

func (app *App) InitWorkers(ctx context.Context) *App {
	workers := []worker_pool.Worker{
		worker_pool.NewWorker(app.periodicalTasks.getTranslatedSpeciality),
		worker_pool.NewWorker(app.periodicalTasks.checkSupportTask),
		worker_pool.NewWorker(app.periodicalTasks.notifyAppointmentTask),
	}
	app.workerPool = worker_pool.NewWorkerPool(workers)
	return app
}

func (app *App) InitUseCases(ctx context.Context) *App {
	app.useCases.createUserUserCase = create_user.NewCreateUserUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.changeLanguageUseCase = change_language.NewChangeLanguageUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.changeStatusUseCase = change_user_status.NewChangeStatusUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.saveMessageUseCase = save_message_log.NewSaveMessageLogUseCase(app.storages.writeTelegramStorage, app.logger)
	app.useCases.updateUserPhoneUseCase = update_user_phone.NewUpdateUserPhoneUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.updateUserPatientIdUseCase = update_user_patient_id.NewUpdateUserPatientIdUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.updateUserBirthDateUseCase = update_user_birth_date.NewUpdateUserBirthDateUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.updateUserThirdNameUseCase = update_user_full_name.NewUpdateFullNameUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.createDraftAppointmentUseCase = create_draft_appointment.NewCreateDraftAppointmentUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.updateDraftAppointmentStatusUseCase = update_appointment_status.NewUpdateAppointmentStatusUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.updateDraftAppointmentIntFieldUseCase = update_int_appointment_field.NewUpdateIntAppointmentFieldUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.updateDraftAppointmentStrFieldUseCase = update_str_appointment_field.NewUpdateStrAppointmentFieldUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.updateUserHomeAddressUseCase = update_home_address.NewUpdateUserHomeAddressUseCase(app.storages.writeUserStorage, app.logger)
	app.useCases.updateDraftAppointmentDateUseCase = update_appointment_date.NewUpdateAppointmentDate(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.cleanDraftAppointmentUseCase = clean_draft_appointment.NewCleanDraftAppointmentUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	app.useCases.fastUpdateDraftAppointmentUseCase = fast_update_draft_appointment_use_case.NewFastUpdateDraftAppointmentUseCase(app.storages.writeDraftAppointmentStorage, app.logger)
	return app
}

func (app *App) InitServices(ctx context.Context) *App {
	app.services.userService = user.NewUserService(
		app.useCases.createUserUserCase,
		app.useCases.changeLanguageUseCase,
		app.useCases.changeStatusUseCase,
		app.useCases.updateUserPhoneUseCase,
		app.useCases.updateUserHomeAddressUseCase,
		app.useCases.updateUserPatientIdUseCase,
		app.useCases.updateUserBirthDateUseCase,
		app.useCases.updateUserThirdNameUseCase,
		app.storages.readUserStorage,
		app.logger,
	)
	// todo исправить
	app.services.appointmentService = appointment.NewAppointmentService(
		app.adapters.appointmentServiceAdapter,
		app.storages.readTranslationStorage,
		app.storages.readDraftAppointmentStorage,
		app.logger,
		app.services.userService,
		app.useCases.createDraftAppointmentUseCase,
		app.useCases.updateDraftAppointmentDateUseCase,
		app.useCases.updateDraftAppointmentStatusUseCase,
		app.useCases.updateDraftAppointmentIntFieldUseCase,
		app.useCases.updateDraftAppointmentStrFieldUseCase,
		app.useCases.cleanDraftAppointmentUseCase,
		app.useCases.fastUpdateDraftAppointmentUseCase,
	)

	app.services.messageService = message.NewMessageService(
		app.useCases.saveMessageUseCase,
		app.storages.readMessageStorage,
		app.storages.readLogsStorage,
		app.logger,
		app.storages.readTranslationStorage,
	)

	app.services.notificationService = notification.NewService(
		app.gateways.WAUSMSGateway,
	)

	return app

}

func (app *App) InitMachine(ctx context.Context) *App {
	app.machine = state_machine.NewUserStateMachine(app.services.userService)
	return app
}

func (app *App) InitTelegram(ctx context.Context) *App {
	app.bot = *telegram.NewTelegramBot(*app.config, app.logger, app.services.messageService)
	return app
}

func (app *App) InitBotGateway(ctx context.Context) *App {
	app.botGateway = bot_gateway.NewBotGateway(
		app.logger, app.bot, app.services.messageService, &app.services.appointmentService,
	)
	return app
}

func (app *App) InitAdapters(ctx context.Context) *App {
	app.adapters.appointmentServiceAdapter = adapter.NewAppointmentServiceAdapter(
		&app.gateways.MisRenoGateway,
	)
	return app
}

func (app *App) InitControllers(ctx context.Context) *App {
	app.controller.telegramWebhookController = controller.NewRestController(
		*app.config,
		app.logger,
		app.bot,
		app.machine,
		app.services.userService,
		&app.services.appointmentService,
		app.services.messageService,
		app.services.notificationService,
		app.botGateway,
	)
	app.controller.telegramWebhookController.InitController()

	app.server = &http.Server{
		Addr:         app.config.HTTPServer.Address,
		Handler:      app.controller.telegramWebhookController,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return app
}
