package usercase

import (
	"context"
	"time"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/util"
)

type usersTableUseCase struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	usersrepo            domain.UsersRepository
	tokenMaker           token.Maker
}

func NewusersTableUserCase(config util.Config, repo domain.UsersRepository, tokenMaker token.Maker) domain.UsersTableUseCase {
	return &usersTableUseCase{
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
		usersrepo:            repo,
		tokenMaker:           tokenMaker,
	}
}

func newUserResponse(user domain.User) domain.UserResponse {
	return domain.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (u *usersTableUseCase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.UserResponse, error) {

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		return domain.UserResponse{}, err
	}

	newUser := domain.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// repository.Create(user)
	user, err := u.usersrepo.Create(ctx, newUser)
	if err != nil {
		return domain.UserResponse{}, err
	}

	rsp := newUserResponse(user)
	return rsp, nil
}

// func (u *usersTableUseCase) PrintLog() string {
// 	return "HI"
// }

func (u *usersTableUseCase) CreateToken(username string, duration time.Duration) (string, *token.Payload, error) {
	token, accessPayload, err := u.tokenMaker.CreateToken(username, duration)
	if err != nil {
		return "", nil, err
	}
	return token, accessPayload, nil
}

func (u *usersTableUseCase) LoginUser(ctx context.Context, req domain.LoginUserRequest) (domain.LoginUserResponse, error) {

	// check user is existed and check password
	user, err := u.usersrepo.GetByUsername(context.Background(), req.Username)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}
	// return loginUserResponse
	accessToken, accessPayload, err := u.CreateToken(user.Username, u.accessTokenDuration)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	refreshToken, refreshPayload, err := u.CreateToken(user.Username, u.refreshTokenDuration)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	// session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
	// 	ID:           refreshPayload.ID,
	// 	Username:     refreshPayload.Username,
	// 	RefreshToken: refreshToken,
	// 	UserAgent:    ctx.Request.UserAgent(),
	// 	ClientIp:     ctx.ClientIP(),
	// 	IsBlocked:    false,
	// 	ExpiresAt:    refreshPayload.ExpiredAt,
	// })

	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errResponse(err))
	// }

	rsp := domain.LoginUserResponse{
		// SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	return rsp, nil
}
