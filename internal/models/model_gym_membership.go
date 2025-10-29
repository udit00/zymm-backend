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


