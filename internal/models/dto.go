package models

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SetReferrerRequest struct {
	ReferrerID int `json:"referrer_id"`
}

type CompleteTaskRequest struct {
	TaskID int `json:"task_id"`
}
