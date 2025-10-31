package attendanceRepo

import (
	"database/sql"
	"errors"
	"fmt"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

func InsertSubmitAttendance(punchInRequestModel models.PunchInAttendanceRequestModel) (*int, error) {
	var id int
	if punchInRequestModel.LocationLat == nil || punchInRequestModel.LocationLong == nil || *punchInRequestModel.LocationLat == "" || *punchInRequestModel.LocationLong == "" {
		return nil, errors.New("Insertion Failed, incomplete data.")
	}
	err := db.DB.QueryRow(`
		INSERT INTO attendance (userId, punchInTime, punchInAddress, punchInLat, punchInLong)
		OUTPUT INSERTED.attendanceId
		VALUES (@p1, getDate(), @p2, @p3, @p4)`,
		punchInRequestModel.UserId, punchInRequestModel.Address, *punchInRequestModel.LocationLat, *punchInRequestModel.LocationLong,
	).Scan(&id)
	if err != nil {
		LogService.LogError(" DB error inserting userMembershipsChangesLogs: ", err)
		return nil, err
	}
	return &id, nil
}

func UpdatePunchOutTime(punchOutRequestModel models.PunchOutAttendanceRequestModel, attendanceId int) error {
	if punchOutRequestModel.LocationLat == nil || punchOutRequestModel.LocationLong == nil || *punchOutRequestModel.LocationLat == "" || *punchOutRequestModel.LocationLong == "" {
		return errors.New("Updation Failed, incomplete data.")
	}
	_, err := db.DB.Exec(`UPDATE attendance SET punchOutTime = getDate(), punchOutAddress = @p1, punchOutLat = @p2, punchOutLong = @p3 WHERE attendanceId = @p4`, &punchOutRequestModel.Address, *punchOutRequestModel.LocationLat, *punchOutRequestModel.LocationLong, attendanceId)
	if err != nil {
		LogService.LogError(" DB error: ", err)
		return err
	}
	return nil
}

func LastAttendanceRecordOfTheUser(userId int) (*models.AttendanceRecord, error) {
	attendanceModel := &models.AttendanceRecord{}
	err := db.DB.QueryRow(
		`SELECT [attendanceId],[userId],[punchInTime],[punchOutTime],[punchInAddress],[punchInLat],[punchInLong],[punchOutAddress],[punchOutLat],[punchOutLong]
		FROM attendance
		WHERE userId = @p1
		ORDER BY punchInTime DESC
		`, userId).Scan(&attendanceModel.AttendanceId, &attendanceModel.UserId, &attendanceModel.PunchInTime, &attendanceModel.PunchOutTime, &attendanceModel.PunchInAddress, &attendanceModel.PunchInLat, &attendanceModel.PunchInLong, &attendanceModel.PunchOutAddress, &attendanceModel.PunchOutLat, &attendanceModel.PunchOutLong)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError(" DB error: ", err)
		return nil, err
	}
	return attendanceModel, nil

}

func GetAllAttendanceByUserId(userId int) ([]models.AttendanceRecord, error) {
	rows, err := db.DB.Query(`
		SELECT [attendanceId],[userId],[punchInTime],[punchOutTime],[punchInAddress],[punchInLat],[punchInLong],[punchOutAddress],[punchOutLat],[punchOutLong]
		FROM attendance
		WHERE userId = @p1
		ORDER BY punchInTime DESC`, userId)
	if err != nil {
		LogService.LogError(" DB query error:", err)
		return nil, err
	}
	defer rows.Close()

	var records []models.AttendanceRecord

	for rows.Next() {
		var record models.AttendanceRecord
		if err := rows.Scan(
			&record.AttendanceId,
			&record.UserId,
			&record.PunchInTime,
			&record.PunchOutTime,
			&record.PunchInAddress,
			&record.PunchInLat,
			&record.PunchInLong,
			&record.PunchOutAddress,
			&record.PunchOutLat,
			&record.PunchOutLong,
		); err != nil {
			LogService.LogError(" Row scan error:", err)
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError(" Rows iteration error:", err)
		return nil, err
	}

	return records, nil
}

func GetAttendanceRecord(attendanceId int) (*models.AttendanceRecord, error) {
	attendanceModel := &models.AttendanceRecord{}
	err := db.DB.QueryRow(
		`SELECT [attendanceId],[userId],[punchInTime],[punchOutTime],[punchInAddress],[punchInLat],[punchInLong],[punchOutAddress],[punchOutLat],[punchOutLong]
		FROM attendance
		WHERE attendanceId = @p1
		ORDER BY punchInTime DESC
		`, attendanceId).Scan(&attendanceModel.AttendanceId, &attendanceModel.UserId, &attendanceModel.PunchInTime, &attendanceModel.PunchOutTime, &attendanceModel.PunchInAddress, &attendanceModel.PunchInLat, &attendanceModel.PunchInLong, &attendanceModel.PunchOutAddress, &attendanceModel.PunchOutLat, &attendanceModel.PunchOutLong)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError(" DB error: ", err)
		return nil, err
	}
	return attendanceModel, nil

}

func DeleteAttendanceRecord(attendanceId int) error {
	exec, err := db.DB.Exec(`delete attendance WHERE attendanceId = @p1`, attendanceId)
	rowsAffected, rowsAffectedErr := exec.RowsAffected()
	if rowsAffectedErr != nil {
		return rowsAffectedErr
	}
	LogService.LogMessage("Rows Affected: " + fmt.Sprint(rowsAffected))
	if err != nil {
		LogService.LogError(" DB error: ", err)
		return err
	}
	return nil

}
