package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pudding-hack/backend/lib/auth"
)

type UserContextKey string
type JWTContextKey string

const (
	UserContextKeyKey UserContextKey = "USER_CONTEXT_KEY"
	JWTContextKeyKey  JWTContextKey  = "JWT_CONTEXT_KEY"
)

func ValidateCurrentUser(cfg APIConfig, r *http.Request) (*auth.User, error) {
	// Get the token from the request header

	secretKey := r.Header.Get("X-Secret-Key")
	if secretKey != "" {
		if secretKey != cfg.SecretKey {
			return nil, fmt.Errorf("Unauthorized")
		}

		return &auth.User{
			ID:   0,
			Name: "System",
		}, nil
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, fmt.Errorf("Unauthorized")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/current-user", cfg.AuthURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response auth.GetCurrentUserResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("error: %+v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("response: %+v", response)
		return nil, NewErrUnauthorized(response.Message)
	}

	return &response.Data, nil
}

func SetUserContext(r *http.Request, user *auth.User) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, UserContextKeyKey, user)
	ctx = context.WithValue(ctx, JWTContextKeyKey, r.Header.Get("Authorization"))
	return r.WithContext(ctx)
}

func GetUserContext(ctx context.Context) *auth.User {
	user := ctx.Value(UserContextKeyKey)
	if user == nil {
		return nil
	}

	return user.(*auth.User)
}

func GetJWTContext(ctx context.Context) string {
	jwt := ctx.Value(JWTContextKeyKey)
	if jwt == nil {
		return ""
	}

	return jwt.(string)
}
