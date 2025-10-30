package models

type EmployeeModel struct {
	EmployeeId     int    `json:"employeeId"`
	UserId         int    `json:"userId"`
	GymId          int    `json:"gymId"`
	CreatedBy      int    `json:"createdBy"`
	StartedWorking string `json:"startedWorking"`
}

// EmployeeWithUserDetails combines employee and user information
type EmployeeWithUserDetails struct {
	EmployeeId     int     `json:"employeeId"`
	UserId         int     `json:"userId"`
	GymId          int     `json:"gymId"`
	CreatedBy      int     `json:"createdBy"`
	StartedWorking string  `json:"startedWorking"`
	UserName       string  `json:"userName"`
	Mobile         string  `json:"mobile"`
	Email          *string `json:"email"`
	Gender         string  `json:"gender"`
	ProfilePic     *string `json:"profilePic"`
	RoleId         int     `json:"roleId"`
	IsActive       bool    `json:"isActive"`
}

type DeactivateEmployeeRequest struct {
	EmployeeId int `json:"employeeId"`
}

type ActivateEmployeeRequest struct {
	EmployeeId int `json:"employeeId"`
}
