package models

type LoginUserDataModel struct {
	UserId      int
	DisplayName string
	Password    string
	RoleId      int
}

type LoginRequestModel struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
	AppVersion    string `json:"app_version"`
	UserAgent     string `json:"user_agent"`
	LocationLat   string `json:"location_lat"`
	LocationLong  string `json:"location_long"`
	IpAddress     string `json:"ip_address"`
}

type LoginResponseModel struct {
	DisplayName  string `json:"displayName"`
	AuthCheckSum string `json:"authCheckSum"`
}
