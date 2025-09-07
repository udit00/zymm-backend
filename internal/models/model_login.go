package models

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	AppVersion string `json:"app_version"`
	Platform   string `json:"platform"`
}
