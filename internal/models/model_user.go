package models

import "time"

type UserRecord struct {
	UserId     int
	UserName   string
	Mobile     string
	Email      *string
	UserPass   string
	DisplayPic *string
	Gender     string
	ProfilePic *string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	RoleId     int
}

type LoginLogsRecord struct {
	UserId       int
	LoginDate    time.Time
	AppVersion   string
	UserAgent    string
	LocationLat  *string
	LocationLong *string
	IpAddress    *string
}
