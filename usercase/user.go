package usercase

import (
	"context"
	"time"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/util"
	"github.com/pkg/errors"
)

type usersTableUseCase struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	sessionUseCase       domain.SessionUseCase
	usersrepo            domain.UsersRepository
	tokenMaker           token.Maker
}

func NewusersTableUserCase(config util.Config, repo domain.UsersRepository, sessionUseCase domain.SessionUseCase, tokenMaker token.Maker) domain.UsersTableUseCase {
	return &usersTableUseCase{
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
		sessionUseCase:       sessionUseCase,
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

func (u *usersTableUseCase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		return nil, domain.ErrorInternalServerError
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
		if err == domain.ErrorInternalServerError || err == domain.ErrorUniqueViolation {
			return nil, err
		}
	}

	rsp := newUserResponse(*user)
	return &rsp, nil
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

func (u *usersTableUseCase) LoginUser(ctx context.Context, req domain.LoginUserRequest) (*domain.LoginUserResponse, error) {

	// check user is existed and check password
	user, err := u.usersrepo.GetByUsername(context.Background(), req.Username)
	if err != nil {
		return nil, err
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, errors.Wrapf(domain.ErrorPermissionNowAllowed, "Mismatched Password")
	}
	// return loginUserResponse
	accessToken, accessPayload, err := u.CreateToken(user.Username, u.accessTokenDuration)
	if err != nil {
		return nil, errors.Wrapf(domain.ErrorInternalServerError, "Create AccessToken Failed")
	}

	refreshToken, refreshPayload, err := u.CreateToken(user.Username, u.refreshTokenDuration)
	if err != nil {
		return nil, errors.Wrapf(domain.ErrorInternalServerError, "Create RefreshToken Failed")
	}

	session := domain.Session{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		// UserAgent:    ctx.Request.UserAgent(),
		// ClientIp:     ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	}

	session, err = u.sessionUseCase.CreateSession(context.Background(), session)
	if err != nil {
		return nil, errors.Wrapf(domain.ErrorInternalServerError, "Create Session Failed")
	}

	return &domain.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(*user),
	}, nil
}
