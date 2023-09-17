package model

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
}
type UserResponse struct {
	Id       int    `json:"id"`
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

type Session struct {
	Username string
	Token    string
	SourceIP string
	LoginAt  string
}

type MQMessage struct {
	EmailTo string `json:"emailto"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
