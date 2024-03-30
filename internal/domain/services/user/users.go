package user

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserService struct {
	createUserUseCase CreateUserUseCase
	readRepo          ReadUserStorage
	logger            *slog.Logger
}

func NewUserService(
	createUserUseCase CreateUserUseCase,
	readRepo ReadUserStorage,
	logger *slog.Logger,
) UserService {
	return UserService{
		createUserUseCase: createUserUseCase,
		readRepo:          readRepo,
		logger:            logger,
	}
}

func (u UserService) RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (userId int64, err error) {

	user, err := u.readRepo.GetUserByTgID(ctx, dto.TgID)
	if err != nil {
		return 0, err
	}

	if user != (entity.User{}) {
		u.logger.Warn("user has registered")
		return 0, nil
	}

	//newUser := dto.ToDomain()

	u.logger.Info("user was not found, trying to register new user", user)

	//// Save new user
	//userId, err = u.createUserUseCase.Execute(ctx, newUser, adminEntity)
	//if err != nil {
	//	return 0, err
	//}

	return userId, nil
}
