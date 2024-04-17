package user

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserService struct {
	createUserUseCase     CreateUserUseCase
	readRepo              ReadUserStorage
	logger                *slog.Logger
	changeLanguageUseCase ChangeLanguageUseCase
}

func NewUserService(
	createUserUseCase CreateUserUseCase,
	changeLanguageUseCase ChangeLanguageUseCase,
	readRepo ReadUserStorage,
	logger *slog.Logger,
) UserService {
	return UserService{
		createUserUseCase:     createUserUseCase,
		readRepo:              readRepo,
		logger:                logger,
		changeLanguageUseCase: changeLanguageUseCase,
	}
}

func (u UserService) RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error) {

	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("%s", err))
		return entity.User{}, err
	}

	if user != (entity.User{}) {
		u.logger.Warn("user has registered")
		return user, nil
	}

	newUser := dto.ToDomain()

	u.logger.Info("user was not found, trying to register new user", user)

	//// Save new user
	_, err = u.createUserUseCase.Execute(ctx, newUser)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u UserService) ChangeLanguage(ctx context.Context, dto tg.TgUserDTO, languageCode string) (user entity.User, err error) {
	user, err = u.readRepo.GetUserByTgID(ctx, dto.TgID)

	if err != nil {
		u.logger.Error(fmt.Sprintf("%s", err))
		return entity.User{}, err
	}

	err = u.changeLanguageUseCase.Execute(ctx, user, languageCode)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
