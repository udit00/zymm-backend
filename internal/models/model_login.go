package models

type LoginRequestModel struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
	AppVersion    string `json:"app_version"`
	Platform      string `json:"platform"`
}

type LoginResponseModel struct {
	DisplayName  string `json:"displayName"`
	AuthCheckSum string `json:"authCheckSum"`
	LastLoggedIn string `json:"lastLoggedIn"`
}
