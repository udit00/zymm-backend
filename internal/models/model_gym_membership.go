package models

type ActiveMembershipUser struct {
	UserId           int     `json:"userId"`
	MembershipId     int     `json:"membershipId"`
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	PlanId           int     `json:"planId"`
	MembershipStatus string  `json:"membershipStatus"`
	UserName         string  `json:"userName"`
	Email            *string `json:"email,omitempty"`
	Mobile           string  `json:"mobile"`
	PlanName         string  `json:"planName"`
	PlanDuration     int     `json:"planDuration"`
	PlanPrice        float64 `json:"planPrice"`
}

// GymMemberWithDetails represents a gym member with their active plan and feedback rating
type GymMemberWithDetails struct {
	UserId           int     `json:"userId"`
	UserName         string  `json:"userName"`
	Mobile           string  `json:"mobile"`
	Email            *string `json:"email"`
	Gender           string  `json:"gender"`
	ProfilePic       *string `json:"profilePic"`
	MembershipId     int     `json:"membershipId"`
	PlanId           int     `json:"planId"`
	PlanName         string  `json:"planName"`
	PlanPrice        float64 `json:"planPrice"`
	PlanDuration     int     `json:"planDuration"`
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	MembershipStatus string  `json:"membershipStatus"`
	FeedbackId       *int    `json:"feedbackId"`
	Rating           *int    `json:"rating"`
	Comments         *string `json:"comments"`
	FeedbackDate     *string `json:"feedbackDate"`
}



