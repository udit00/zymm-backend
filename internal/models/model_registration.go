package models

import "time"

type RegistrationLogsRecord struct {
	UserId     int       `json:"userId"`
	UserName   string    `json:"userName"`
	UserPass   string    `json:"userPass"`
	Gender     string    `json:"gender"`
	Mobile     string    `json:"mobile"`
	Email      *string   `json:"email"`
	ProfilePic *string   `json:"profilePic"`
	RoleId     int       `json:"roleId"`
	AppVersion string    `json:"appVersion"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
}

// RegistrationRequestModel represents the expected payload for user registration
type RegistrationApiRequestModel struct {
	DisplayPic   *string `json:"displayPic"`
	DisplayName  string  `json:"displayName"`
	Mobile       string  `json:"mobile"`
	Email        *string `json:"email"`
	Password     string  `json:"password"`
	Gender       string  `json:"gender"`
	AppVersion   string  `json:"appVersion"`
	UserAgent    string  `json:"userAgent"`
	LocationLat  string  `json:"locationLat"`
	LocationLong string  `json:"locationLong"`
	IpAddress    string  `json:"ipAddress"`
}

type RegistrationOwnerApiRequestModel struct {
	DisplayPic         *string `json:"displayPic"`
	DisplayName        string  `json:"displayName"`
	Mobile             string  `json:"mobile"`
	OwnerPersonalEmail *string `json:"ownerPersonalEmail"`
	Password           string  `json:"password"`
	Gender             string  `json:"gender"`
	AppVersion         string  `json:"appVersion"`
	UserAgent          string  `json:"userAgent"`
	IpAddress          string  `json:"ipAddress"`

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
	DisplayName  string `json:"displayName"`
	AuthCheckSum string `json:"authCheckSum"`
}
