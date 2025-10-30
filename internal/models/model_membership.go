package models

import "time"

type PlanRecord struct {
	PlanId       int        `json:"planId"`
	PlanBanner   *string    `json:"planBanner"`
	PlanName     string     `json:"planName"`
	PlanDesc     string     `json:"planDesc"`
	PlanPrice    float64    `json:"planPrice"`
	PlanDuration int        `json:"planDuration"`
	IsActive     *bool      `json:"isActive"`
	CreatedBy    *int       `json:"createdBy"`
	CreatedAt    *time.Time `json:"createdAt"`
	GymId        *int       `json:"gymId"`
}

type PlanChangeLog struct {
	LogId           int        `json:"logId"`
	PlanId          *int       `json:"planId"`
	ChangedBy       *int       `json:"changedBy"`
	ChangeType      string     `json:"changeType"`
	ChangeDetails   string     `json:"changeDetails"`
	ChangedAt       *time.Time `json:"changedAt"`
	OldPlanBanner   *string    `json:"oldPlanBanner"`
	OldPlanName     string     `json:"oldPlanName"`
	OldPlanDesc     string     `json:"oldPlanDesc"`
	OldPlanPrice    float64    `json:"oldPlanPrice"`
	OldPlanDuration int        `json:"oldPlanDuration"`
	OldIsActive     bool       `json:"oldIsActive"`
	NewPlanBanner   *string    `json:"newPlanBanner"`
	NewPlanName     string     `json:"newPlanName"`
	NewPlanDesc     string     `json:"newPlanDesc"`
	NewPlanPrice    float64    `json:"newPlanPrice"`
	NewPlanDuration int        `json:"newPlanDuration"`
	NewIsActive     bool       `json:"newIsActive"`
}

// CreatePlanRequest represents the API payload to create a plan
type UpsertPlanRequest struct {
	PlanId       *int    `json:"planId"`
	PlanBanner   *string `json:"planBanner"`
	PlanName     string  `json:"planName"`
	PlanDesc     string  `json:"planDesc"`
	PlanPrice    float64 `json:"planPrice"`
	PlanDuration int     `json:"planDuration"`
	GymId        int     `json:"gymId"`
	UserAgent    string  `json:"userAgent"`
	AppVersion   string  `json:"appVersion"`
	IsActive     bool    `json:"isActive"`
}

type UserMembership struct {
	MembershipId     int    `json:"membershipId"`
	UserId           int    `json:"userId"`
	PlanId           int    `json:"planId"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	IsActive         bool   `json:"isActive"`
	MembershipStatus string `json:"membershipStatus"`
	CreatedAt        string `json:"createdAt"`
}

type RequestPlanFromUserToGymModel struct {
	PlanId int `json:"planId`
}

type CancelRequestMembershipModel struct {
	MembershipId int `json:"membershipId"`
}

type TakeActionOnMembershipRequestModel struct {
	MembershipId int    `json:"membershipId"`
	ActionTaken  string `json:"actionTaken"`
}

// MemberWithPendingFees represents a member whose membership is expired or expiring soon
type MemberWithPendingFees struct {
	UserId           int     `json:"userId"`
	UserName         string  `json:"userName"`
	Mobile           string  `json:"mobile"`
	Email            *string `json:"email"`
	ProfilePic       *string `json:"profilePic"`
	MembershipId     int     `json:"membershipId"`
	PlanId           int     `json:"planId"`
	PlanName         string  `json:"planName"`
	PlanPrice        float64 `json:"planPrice"`
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	DaysUntilExpiry  int     `json:"daysUntilExpiry"` // Negative if expired
	IsExpired        bool    `json:"isExpired"`
	MembershipStatus string  `json:"membershipStatus"`
}

// SendFeeReminderRequest represents the request to send fee reminders
type SendFeeReminderRequest struct {
	UserIds []int `json:"userIds"` // Array of user IDs to send notifications to
}
