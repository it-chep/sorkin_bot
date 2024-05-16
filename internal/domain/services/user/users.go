package user

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserService struct {
	createUserUseCase          CreateUserUseCase
	readRepo                   ReadUserStorage
	logger                     *slog.Logger
	changeLanguageUseCase      ChangeLanguageUseCase
	changeStateUseCase         ChangeStateUseCase
	updateUserPhoneUseCase     UpdateUserPhoneUseCase
	updateUserPatientIdUseCase UpdateUserPatientIdUseCase
	updateUserBirthDateUseCase UpdateUserBirthDateUseCase
	updateUserThirdNameUseCase UpdateUserThirdNameUseCase
}

func NewUserService(
	createUserUseCase CreateUserUseCase,
	changeLanguageUseCase ChangeLanguageUseCase,
	changeStateUseCase ChangeStateUseCase,
	updateUserPhoneUseCase UpdateUserPhoneUseCase,
	updateUserPatientIdUseCase UpdateUserPatientIdUseCase,
	updateUserBirthDateUseCase UpdateUserBirthDateUseCase,
	updateUserThirdNameUseCase UpdateUserThirdNameUseCase,
	readRepo ReadUserStorage,
	logger *slog.Logger,
) UserService {
	return UserService{
		createUserUseCase:          createUserUseCase,
		readRepo:                   readRepo,
		logger:                     logger,
		changeStateUseCase:         changeStateUseCase,
		changeLanguageUseCase:      changeLanguageUseCase,
		updateUserPhoneUseCase:     updateUserPhoneUseCase,
		updateUserPatientIdUseCase: updateUserPatientIdUseCase,
		updateUserBirthDateUseCase: updateUserBirthDateUseCase,
		updateUserThirdNameUseCase: updateUserThirdNameUseCase,
	}
}

func (u UserService) GetUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.GetUser"

	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.RegisterNewUser"

	user, err = u.GetUser(ctx, dto)
	if err != nil {
		return entity.User{}, err
	}
	//todo add pointer
	if !reflect.ValueOf(user.GetTgId()).IsZero() {
		u.logger.Warn("user has registered")
		return user, nil
	}

	newUser := dto.ToDomain()

	u.logger.Info("user was not found, trying to register new user", user, op)

	//// Save new user
	_, err = u.createUserUseCase.Execute(ctx, newUser)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u UserService) ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.ChangeLanguage"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("%s %s", err, op))
		return entity.User{}, err
	}

	err = u.changeLanguageUseCase.Execute(ctx, user, languageCode)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) UpdatePatientId(ctx context.Context, user entity.User, patientId *int) (err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdatePatientId"
	err = u.updateUserPatientIdUseCase.Execute(ctx, user, *patientId)
	if err != nil {
		u.logger.Error(fmt.Sprintf("%s %s", err, op))
		return err
	}
	return err
}

func (u UserService) ChangeState(ctx context.Context, dto tg.TgUserDTO, state string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.ChangeState"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}

	err = u.changeStateUseCase.Execute(ctx, user, state)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdatePhone"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}

	err = u.updateUserPhoneUseCase.Execute(ctx, user, phone)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) UpdateThirdName(ctx context.Context, dto tg.TgUserDTO, thirdName string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdateThirdName"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}

	err = u.updateUserThirdNameUseCase.Execute(ctx, user, thirdName)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdateBirthDate"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}

	err = u.updateUserBirthDateUseCase.Execute(ctx, user, birthDate)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}
