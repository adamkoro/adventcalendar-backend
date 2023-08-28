package api

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	Username    string `json:"username"`
	TokenString string `json:"token"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
}
type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Status string `json:"status"`
}
