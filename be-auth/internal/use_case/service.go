package use_case

import (
	"time"

	"github.com/pudding-hack/backend/be-auth/internal/model"
	"github.com/pudding-hack/backend/lib"
)

type userRepository interface {
	GetUserByUsername(username string) (*model.User, error)
	GetRoleByID(id int) (*model.Role, error)
	StoredAccessTokenToRedis(token string, userID string) (err error)
	StoredRefreshTokenToRedis(token string, userID string) (err error)
	GetAccessTokenFromRedis(token string) (userID string, err error)
	GetRefreshTokenFromRedis(token string) (userID string, err error)
	DeleteAccessTokenFromRedis(token string) (err error)
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

func (s *service) Login(username, password string) (res LoginResponse, err error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return res, err
	}

	if !lib.CheckPasswordHash(password, user.Password) {
		return res, lib.NewErrUnauthorized("invalid password")
	}

	role, err := s.repo.GetRoleByID(user.RoleID)
	if err != nil {
		return res, err
	}

	claimsAccesToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.LoginExirationDuration))
	accessToken, err := lib.GenerateToken(claimsAccesToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredAccessTokenToRedis(accessToken, user.ID)

	claimsRefreshToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.RefreshExpirationDuration))
	refreshToken, err := lib.GenerateToken(claimsRefreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredRefreshTokenToRedis(refreshToken, user.ID)

	return fromModuleToLoginResponse(user, role, accessToken, refreshToken), nil
}

func (s *service) GetCurrentUser(accessToken string) (res User, err error) {
	_, err = s.repo.GetAccessTokenFromRedis(accessToken)
	if err != nil {
		return res, lib.NewErrUnauthorized(lib.ErrTokenExpired)
	}

	token, err := lib.ParseToken(accessToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	user, err := s.repo.GetUserByUsername(token.Username)
	if err != nil {
		return res, err
	}

	role, err := s.repo.GetRoleByID(user.RoleID)
	if err != nil {
		return res, err
	}

	return fromModuleToUserResponse(user, role), nil
}

func (s *service) RefreshToken(refreshToken string) (res LoginResponse, err error) {
	token, err := lib.ParseToken(refreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	_, err = s.repo.GetRefreshTokenFromRedis(refreshToken)
	if err != nil {
		return res, lib.NewErrUnauthorized("refresh token not found")
	}

	user, err := s.repo.GetUserByUsername(token.Username)
	if err != nil {
		return res, err
	}

	role, err := s.repo.GetRoleByID(user.RoleID)
	if err != nil {
		return res, err
	}

	claimsAccesToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.LoginExirationDuration))
	accessToken, err := lib.GenerateToken(claimsAccesToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}
	go s.repo.StoredAccessTokenToRedis(accessToken, user.ID)

	claimsRefreshToken := lib.NewClaims(user.ID, user.Username, user.Email, s.cfg.App.Name, time.Now().Add(s.cfg.JWT.RefreshExpirationDuration))
	refreshToken, err = lib.GenerateToken(claimsRefreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return res, lib.NewErrUnauthorized(err.Error())
	}

	go s.repo.StoredRefreshTokenToRedis(refreshToken, user.ID)

	return fromModuleToLoginResponse(user, role, accessToken, refreshToken), nil
}

func (s *service) Logout(accessToken string) (err error) {
	_, err = s.repo.GetAccessTokenFromRedis(accessToken)
	if err != nil {
		return lib.NewErrUnauthorized("access token not found")
	}

	err = s.repo.DeleteAccessTokenFromRedis(accessToken)
	if err != nil {
		return err
	}

	return nil
}
