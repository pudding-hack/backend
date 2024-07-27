package model

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/pudding-hack/backend/lib"
)

type repository struct {
	cfg   *lib.Config
	db    *sqlx.DB
	redis *redis.Pool
}

func New(cfg *lib.Config, db *sqlx.DB, redis *redis.Pool) *repository {
	return &repository{
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
}

func (r *repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1 AND deleted_at is NULL", username)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (r *repository) GetRoleByID(id int) (*Role, error) {
	var role Role
	err := r.db.Get(&role, "SELECT * FROM roles WHERE id = $1 AND deleted_at is NULL", id)
	if err != nil {
		return &Role{}, err
	}

	return &role, nil
}

func (r *repository) StoredAccessTokenToRedis(token string, userID int64) (err error) {
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

func (r *repository) StoredRefreshTokenToRedis(token string, userID int64) (err error) {
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

func (r *repository) GetAccessTokenFromRedis(token string) (userID int64, err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token:%s", token)

	userID, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return
}

func (r *repository) GetRefreshTokenFromRedis(token string) (userID int64, err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("refresh_token:%s", token)

	userID, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return
}

func (r *repository) DeleteAccessTokenFromRedis(token string) (err error) {
	conn := r.redis.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token:%s", token)

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
