package models

import "time"

type RegistrationLogsRecord struct {
	UserId      int       `json:"userId"`
	UserName    string    `json:"userName"`
	UserPass    string    `json:"userPass"`
	Gender      string    `json:"gender"`
	Mobile      string    `json:"mobile"`
	Email       *string   `json:"email"`
	ProfilePic  *string   `json:"profilePic"`
	RoleId      int       `json:"roleId"`
	AppVersion  string    `json:"appVersion"`
	AppPlatform string    `json:"appPlatform"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RegistrationRequestModel represents the expected payload for user registration
type RegistrationApiRequestModel struct {
	DisplayPic  *string `json:"displayPic"`
	DisplayName string  `json:"displayName"`
	Mobile      string  `json:"mobile"`
	Email       *string `json:"email"`
	Password    string  `json:"password"`
	Gender      string  `json:"gender"`
	AppVersion  string  `json:"appVersion"`
	AppPlatform string  `json:"appPlatform"`
}

type RegistrationOwnerApiRequestModel struct {
	DisplayPic         *string `json:"displayPic"`
	DisplayName        string  `json:"displayName"`
	Mobile             string  `json:"mobile"`
	OwnerPersonalEmail *string `json:"ownerPersonalEmail"`
	Password           string  `json:"password"`
	Gender             string  `json:"gender"`
	AppVersion         string  `json:"appVersion"`
	AppPlatform        string  `json:"appPlatform"`

	GymName       string `json:"gymName"`
	State         string `json:"state"`
	City          string `json:"city"`
	GymAddress    string `json:"gymAddress"`
	ContactNo     string `json:"gymOfficialContactNo"`
	OfficialEmail string `json:"gymOfficialEmail"`
	LocationLat   string `json:"gymOfficialLocationLat"`
	LocationLong  string `json:"gymOfficialLocationLong"`
}

// RegistrationResponseModel represents the response after successful registration
type RegistrationApiResponseModel struct {
	UserId     int    `json:"userId"`
	UserName   string `json:"userName"`
	Registered bool   `json:"registered"`
	Message    string `json:"message"`
}
