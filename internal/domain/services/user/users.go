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
	updateUserFullNameUseCase  UpdateUserFullNameUseCase
}

func NewUserService(
	createUserUseCase CreateUserUseCase,
	changeLanguageUseCase ChangeLanguageUseCase,
	changeStateUseCase ChangeStateUseCase,
	updateUserPhoneUseCase UpdateUserPhoneUseCase,
	updateUserPatientIdUseCase UpdateUserPatientIdUseCase,
	updateUserBirthDateUseCase UpdateUserBirthDateUseCase,
	updateUserFullNameUseCase UpdateUserFullNameUseCase,
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
		updateUserFullNameUseCase:  updateUserFullNameUseCase,
	}
}

func (u UserService) GetUser(ctx context.Context, tgId int64) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.GetUser"

	user, err = u.readRepo.GetUserByTgID(ctx, tgId)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.RegisterNewUser"

	user, err = u.GetUser(ctx, dto.TgID)
	if err != nil {
		return entity.User{}, err
	}
	//todo add pointer
	if !reflect.ValueOf(user.GetTgId()).IsZero() {
		u.logger.Warn("user has registered")
		return user, nil
	}

	newUser := dto.ToDomain(
		[]entity.UserOpt{
			entity.WithUsrUsername(&dto.UserName),
			entity.WithUsrLastName(&dto.LastName),
		},
	)

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
	// todo add pointer to user
	user.SetLanguageCode(languageCode)

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

func (u UserService) ChangeState(ctx context.Context, tgId int64, state string) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.ChangeState"
	user, err = u.readRepo.GetUserByTgID(ctx, tgId)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}

	user.SetState(state)
	err = u.changeStateUseCase.Execute(ctx, user)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, result bool, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdatePhone"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}

	if !u.validatePhoneMessage(phone) {
		return user, false, nil
	}

	err = u.updateUserPhoneUseCase.Execute(ctx, user, phone)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}
	return user, false, nil
}

func (u UserService) UpdateFullName(ctx context.Context, dto tg.TgUserDTO, fullName string) (user entity.User, result bool, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdateThirdName"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}

	if !u.validateNameMessage(fullName) {
		return user, false, nil
	}

	err = u.updateUserFullNameUseCase.Execute(ctx, user, fullName)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}
	return user, true, nil
}

func (u UserService) UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, result bool, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.UpdateBirthDate"
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}

	if !u.validateBirthDateMessage(birthDate) {
		return user, false, nil
	}

	err = u.updateUserBirthDateUseCase.Execute(ctx, user, birthDate)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return entity.User{}, false, err
	}
	return user, true, nil
}
