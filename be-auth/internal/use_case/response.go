package use_case

import "github.com/pudding-hack/backend/be-auth/internal/model"

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	Signature string `json:"signature"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func fromModuleToUserResponse(user *model.User, role *model.Role) User {
	return User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: Role{
			ID:   role.ID,
			Name: role.Name,
		},
	}
}

func fromModuleToLoginResponse(user *model.User, role *model.Role, accessToken, refreshToken string) LoginResponse {
	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role: Role{
				ID:   role.ID,
				Name: role.Name,
			},
		},
	}
}
