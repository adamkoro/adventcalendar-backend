package api

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginResponse struct {
	Username    string `json:"username"`
	TokenString string `json:"token"`
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type getUserRequest struct {
	Username string `json:"username" binding:"required"`
}
type getUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Status string `json:"status"`
}
