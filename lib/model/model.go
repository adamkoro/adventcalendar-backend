package model

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
