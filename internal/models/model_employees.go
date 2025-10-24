package models

type EmployeeModel struct {
	EmployeeId     int    `json:"employeeId"`
	UserId         int    `json:"userId"`
	GymId          int    `json:"gymId"`
	StartedWorking string `json:"startedWorking"`
}
