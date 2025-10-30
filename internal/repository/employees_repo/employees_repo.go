package employeesRepo

import (
	"database/sql"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

/*
create table employees (

	employeeId          INT         primary key identity(1, 1) not null,
	userId              int         not null,
	gymId               INT         not null,
	startedWorking      DATETIME    default getdate() not null,
	FOREIGN KEY (userId) REFERENCES users(userId),
	FOREIGN KEY (gymId) REFERENCES gym(gymId)

)
*/
func InsertEmployees(employeeModel models.EmployeeModel) (*int, error) {
	var id int
	// Use OUTPUT INSERTED.employeeId to get the inserted id
	err := db.DB.QueryRow(`
        INSERT INTO employees (userId, gymId, createdBy)
        OUTPUT INSERTED.employeeId
        VALUES (@p1, @p2, @p3)`,
		employeeModel.UserId, employeeModel.GymId, employeeModel.CreatedBy).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting plan: ", err)
		return nil, err
	}
	return &id, nil
}

func GetEmployeeByUserId(userId int) (*models.EmployeeModel, error) {
	employeeModel := &models.EmployeeModel{}
	err := db.DB.QueryRow(
		`Select employeeId, userId, gymId, createdBy, startedWorking
		from employees
		where userId = @p1
		`, userId).Scan(&employeeModel.EmployeeId, &employeeModel.UserId, &employeeModel.GymId, &employeeModel.CreatedBy, &employeeModel.StartedWorking)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return employeeModel, nil
}

func GetAllEmployees(gymId int) ([]models.EmployeeModel, error) {
	rows, err := db.DB.Query(`
		Select employeeId, userId, gymId, createdBy, startedWorking
		from employees
		where gymId = @p1
	`, gymId)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var employees []models.EmployeeModel

	for rows.Next() {
		var currentEmp models.EmployeeModel
		err := rows.Scan(
			&currentEmp.EmployeeId,
			&currentEmp.UserId,
			&currentEmp.GymId,
			&currentEmp.CreatedBy,
			&currentEmp.StartedWorking,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		employees = append(employees, currentEmp)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(employees) == 0 {
		return nil, sql.ErrNoRows
	}

	return employees, nil
}

func GetEmployeeByEmployeeId(empId int) (*models.EmployeeModel, error) {
	employeeModel := &models.EmployeeModel{}
	err := db.DB.QueryRow(
		`Select employeeId, userId, gymId, createdBy, startedWorking
		from employees
		where employeeId = @p1
		`, empId).Scan(&employeeModel.EmployeeId, &employeeModel.UserId, &employeeModel.GymId, &employeeModel.CreatedBy, &employeeModel.StartedWorking)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return employeeModel, nil
}
