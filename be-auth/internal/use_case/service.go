package use_case

import (
	"context"
	"errors"
	"time"

	"github.com/pudding-hack/backend/be-auth/internal/model"
	"github.com/pudding-hack/backend/lib"
)

type userRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetRoleByID(ctx context.Context, id int) (*model.Role, error)
	StoredAccessTokenToRedis(ctx context.Context, token string, userID string) (err error)
	StoredRefreshTokenToRedis(ctx context.Context, token string, userID string) (err error)
	GetAccessTokenFromRedis(ctx context.Context, token string) (userID string, err error)
	GetRefreshTokenFromRedis(ctx context.Context, token string) (userID string, err error)
	DeleteAccessTokenFromRedis(ctx context.Context, token string) (err error)
}

type service struct {
	cfg  *lib.Config
	repo userRepository
}

func NewService(cfg *lib.Config, repo userRepository) *service {
	return &service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *service) Login(ctx context.Context, username, password string) (res LoginResponse, err error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return res, lib.NewErrUnauthorized("invalid username")
		}
		return res, err
	}

	if !lib.CheckPasswordHash(password, user.Password) {
		return res, lib.NewErrUnauthorized("invalid password")
	}

	role, err := s.repo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		if errors.Is(err, lib.ErrSqlErrorNotFound) {
			return res, lib.NewErrUnauthorized("role not found")
		}

		return res, err
	}

	claimsAccesToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.LoginExirationDuration))
	accessToken, err := lib.GenerateToken(claimsAccesToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredAccessTokenToRedis(ctx, accessToken, user.ID)

	claimsRefreshToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.RefreshExpirationDuration))
	refreshToken, err := lib.GenerateToken(claimsRefreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredRefreshTokenToRedis(ctx, refreshToken, user.ID)

	return fromModuleToLoginResponse(user, role, accessToken, refreshToken), nil
}

func (s *service) GetCurrentUser(ctx context.Context, accessToken string) (res User, err error) {
	_, err = s.repo.GetAccessTokenFromRedis(ctx, accessToken)
	if err != nil {
		return res, lib.NewErrUnauthorized(lib.ErrTokenExpired)
	}

	token, err := lib.ParseToken(accessToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	user, err := s.repo.GetUserByUsername(ctx, token.Username)
	if err != nil {
		return res, err
	}

	role, err := s.repo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		return res, err
	}

	return fromModuleToUserResponse(user, role), nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (res LoginResponse, err error) {
	token, err := lib.ParseToken(refreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	_, err = s.repo.GetRefreshTokenFromRedis(ctx, refreshToken)
	if err != nil {
		return res, lib.NewErrUnauthorized("refresh token not found")
	}

	user, err := s.repo.GetUserByUsername(ctx, token.Username)
	if err != nil {
		return res, err
	}

	role, err := s.repo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		return res, err
	}

	claimsAccesToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.LoginExirationDuration))
	accessToken, err := lib.GenerateToken(claimsAccesToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredAccessTokenToRedis(ctx, accessToken, user.ID)

	claimsRefreshToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.RefreshExpirationDuration))
	refreshToken, err = lib.GenerateToken(claimsRefreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	go s.repo.StoredRefreshTokenToRedis(ctx, refreshToken, user.ID)

	return fromModuleToLoginResponse(user, role, accessToken, refreshToken), nil
}

func (s *service) Logout(ctx context.Context, accessToken string) (err error) {
	_, err = s.repo.GetAccessTokenFromRedis(ctx, accessToken)
	if err != nil {
		return lib.NewErrUnauthorized("access token not found")
	}

	err = s.repo.DeleteAccessTokenFromRedis(ctx, accessToken)
	if err != nil {
		return err
	}

	return nil
}
