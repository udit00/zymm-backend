package models

type LoginRequestModel struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	AppVersion string `json:"app_version"`
	Platform   string `json:"platform"`
}

type LoginResponseModel struct {
	DisplayName  string `json:"displayName"`
	AuthCheckSum string `json:"authCheckSum"`
	LastLoggedIn string `json:"lastLoggedIn"`
}
