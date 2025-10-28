package models

type AttendanceRecord struct {
	AttendanceId    int     `json:"attendanceId"`
	UserId          int     `json:"userId"`
	PunchInTime     string  `json:"punchInTime"`
	PunchOutTime    *string `json:"punchOutTime"`
	PunchInAddress  string  `json:"punchInAddress"`
	PunchInLat      string  `json:"punchInLat"`
	PunchInLong     string  `json:"punchInLong"`
	PunchOutAddress *string `json:"punchOutAddress"`
	PunchOutLat     *string `json:"punchOutLat"`
	PunchOutLong    *string `json:"punchOutLong"`
}

type PunchInAttendanceRequestModel struct {
	UserId       int     `json:"userId"`
	Address      string  `json:"address"`
	LocationLat  *string `json:"locationLat"`
	LocationLong *string `json:"locationLong"`
}

type PunchOutAttendanceRequestModel struct {
	UserId       int     `json:"userId"`
	Address      string  `json:"address"`
	LocationLat  *string `json:"locationLat"`
	LocationLong *string `json:"locationLong"`
}

type GetAllAttendanceRequestModel struct {
	UserId *int `json:"userId"`
}

type DeleteAttendanceRequest struct {
	AttendanceId int `json:"attendanceId"`
}
