package gorm

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	Username          string `gorm:"primary_key"`
	HashedPassword    string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *HttpServer) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	newUser := User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// repository.Create(user)
	result := server.store.Create(&newUser)
	if result.Error != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	rsp := newUserResponse(newUser)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	// SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user_response"`
}

func (server *HttpServer) loginUser(ctx *gin.Context) {
	// bind req body
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user := User{
		Username: req.Username,
	}
	// check user is existed and check password
	result := server.store.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	// return loginUserResponse
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
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

	rsp := loginUserResponse{
		// SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

// func (server *GrpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
// 	payload, err := server.authorizeUser(ctx)
// 	if err != nil {
// 		return nil, unauthenticationError(err)
// 	}
// 	if violations := validateUpdateUserRequest(req); violations != nil {
// 		return nil, invalidArgumentError(violations)
// 	}

// 	if payload.Username != req.Username {
// 		return nil, permissionDeniedError(err)
// 	}
// 	arg := db.UpdateUserParams{
// 		Username: req.GetUsername(),
// 		FullName: sql.NullString{
// 			String: req.GetFullName(),
// 			Valid:  req.FullName != nil,
// 		},
// 		Email: sql.NullString{
// 			String: req.GetEmail(),
// 			Valid:  req.Email != nil,
// 		},
// 	}

// 	if req.Password != nil {
// 		hashedPassword, err := util.HashedPassword(req.GetPassword())
// 		if err != nil {
// 			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
// 		}
// 		arg.HashedPassword = sql.NullString{
// 			String: hashedPassword,
// 			Valid:  true,
// 		}
// 		arg.PasswordChangedAt = sql.NullTime{
// 			Time:  time.Now(),
// 			Valid: true,
// 		}
// 	}

// 	user, err := server.store.UpdateUser(ctx, arg)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, status.Errorf(codes.NotFound, "user not founde %s", err)
// 		}
// 		return nil, status.Errorf(codes.Internal, "failed to update user %s", err)
// 	}
// 	rsp := &pb.UpdateUserResponse{
// 		User: convertUser(user),
// 	}
// 	return rsp, nil
// }

// func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if err := val.ValidateUserName(req.GetUsername()); err != nil {
// 		violations = append(violations, FieldViolation("username", err))
// 	}
// 	if req.Password != nil {
// 		if err := val.ValidatePassword(req.GetPassword()); err != nil {
// 			violations = append(violations, FieldViolation("password", err))
// 		}
// 	}
// 	if req.Email != nil {
// 		if err := val.ValidateEmail(req.GetEmail()); err != nil {
// 			violations = append(violations, FieldViolation("email", err))
// 		}
// 	}
// 	if req.FullName != nil {
// 		if err := val.ValidateFullName(req.GetFullName()); err != nil {
// 			violations = append(violations, FieldViolation("fullname", err))
// 		}
// 	}
// 	return violations
// }
