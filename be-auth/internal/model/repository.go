package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

type repository struct {
	cfg   *lib.Config
	db    conn.Connection
	redis *redis.Pool
}

func New(cfg *lib.Config, db conn.Connection, redis *redis.Pool) *repository {
	return &repository{
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.Get(ctx, &user, "SELECT * FROM users WHERE username = $1 AND deleted_at is NULL", username)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (r *repository) GetRoleByID(ctx context.Context, id int) (*Role, error) {
	var role Role
	err := r.db.Get(ctx, &role, "SELECT * FROM roles WHERE id = $1 AND deleted_at is NULL", id)
	if err != nil {
		return &Role{}, err
	}

	return &role, nil
}

func (r *repository) StoredAccessTokenToRedis(ctx context.Context, token string, userID string) (err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token:%s", token)
	expired := time.Now().Add(r.cfg.JWT.LoginExirationDuration).Unix()

	_, err = conn.Do("SET", key, userID, "EX", expired)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) StoredRefreshTokenToRedis(ctx context.Context, token string, userID string) (err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("refresh_token:%s", token)
	expired := time.Now().Add(r.cfg.JWT.RefreshExpirationDuration).Unix()

	_, err = conn.Do("SET", key, userID, "EX", expired)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAccessTokenFromRedis(ctx context.Context, token string) (userID string, err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token:%s", token)

	userID, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return
}

func (r *repository) GetRefreshTokenFromRedis(ctx context.Context, token string) (userID string, err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("refresh_token:%s", token)

	userID, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return
}

func (r *repository) DeleteAccessTokenFromRedis(ctx context.Context, token string) (err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token:%s", token)

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
