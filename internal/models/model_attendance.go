package models

type AttendanceRecord struct {
	AttendanceId int     `json:"attendanceId"`
	UserId       int     `json:"userId"`
	PunchInTime  string  `json:"punchInTime"`
	PunchOutTime *string `json:"punchOutTime"`
	PunchInLat   string  `json:"punchInLat"`
	PunchInLong  string  `json:"punchInLong"`
	PunchOutLat  *string `json:"punchOutLat"`
	PunchOutLong *string `json:"punchOutLong"`
}

type PunchInAttendanceRequestModel struct {
	UserId       int     `json:"userId"`
	LocationLat  *string `json:"locationLat"`
	LocationLong *string `json:"locationLong"`
}

type PunchOutAttendanceRequestModel struct {
	UserId       int     `json:"userId"`
	LocationLat  *string `json:"locationLat"`
	LocationLong *string `json:"locationLong"`
}

type GetAllAttendanceRequestModel struct {
	UserId *int `json:"userId"`
}

type DeleteAttendanceRequest struct {
	AttendanceId int `json:"attendanceId"`
}
