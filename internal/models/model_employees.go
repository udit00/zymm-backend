package models

type EmployeeModel struct {
	EmployeeId     int    `json:"employeeId"`
	UserId         int    `json:"userId"`
	GymId          int    `json:"gymId"`
	CreatedBy      int    `json:"createdBy"`
	StartedWorking string `json:"startedWorking"`
}
