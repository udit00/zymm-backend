package models

import (
	"time"
)

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
	IsActive   bool
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

type SelfDataResponse struct {
	UserId                  int             `json:"userId"`
	UserName                string          `json:"userName"`
	Mobile                  string          `json:"mobile"`
	Email                   *string         `json:"email"`
	Gender                  string          `json:"gender"`
	RoleId                  int             `json:"roleId"`
	ProfilePic              *string         `json:"profilePic"`
	MembershipId            *int            `json:"membershipId"`
	PlanId                  *int            `json:"planId"`
	GymId                   *int            `json:"gymId"`
	ActiveMembershipDetails *UserMembership `json:"activeMembershipDetails"`
	VisitedToday            bool            `json:"visitedToday"`
	PlanDetails             *PlanRecord     `json:"planDetails"`
}

type ChangePasswordRequestModel struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponseModel struct {
	AuthToken string `json:"authToken"`
}

type DeleteProfileRequestModel struct {
	Password string `json:"password"`
}
