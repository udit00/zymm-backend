package feedbackRepo

import (
	"database/sql"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

// InsertFeedback inserts a feedback row and returns the newly created id
func InsertFeedback(f models.FeedbackRecord) (*int, error) {
	var id int
	err := db.DB.QueryRow(`
        INSERT INTO feedback (rating, comments, gymId, createdBy)
        OUTPUT INSERTED.feedbackId
        VALUES (@p1, @p2, @p3, @p4)`,
		f.Rating, f.Comments, f.GymId, f.CreatedBy,
	).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting feedback: ", err)
		return nil, err
	}
	return &id, nil
}

// GetFeedbackById retrieves a single feedback by feedbackId
func GetFeedbackById(feedbackId int) (*models.FeedbackRecord, error) {
	feedback := &models.FeedbackRecord{}
	err := db.DB.QueryRow(`
		SELECT feedbackId, rating, comments, gymId, createdBy, createdAt
		FROM feedback
		WHERE feedbackId = @p1`,
		feedbackId).Scan(
		&feedback.FeedbackId, &feedback.Rating, &feedback.Comments, &feedback.GymId, &feedback.CreatedBy, &feedback.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return feedback, nil
}

// GetAllFeedbacksByGymId retrieves all feedbacks for a specific gym
func GetAllFeedbacksByGymId(gymId int) ([]models.FeedbackRecord, error) {
	query := `SELECT feedbackId, rating, comments, gymId, createdBy, createdAt
		FROM feedback
		WHERE gymId = @p1
		ORDER BY createdAt DESC`
	
	rows, err := db.DB.Query(query, gymId)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.FeedbackRecord

	for rows.Next() {
		var feedback models.FeedbackRecord
		err := rows.Scan(
			&feedback.FeedbackId,
			&feedback.Rating,
			&feedback.Comments,
			&feedback.GymId,
			&feedback.CreatedBy,
			&feedback.CreatedAt,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(feedbacks) == 0 {
		return []models.FeedbackRecord{}, nil
	}

	return feedbacks, nil
}

// UpdateFeedback updates an existing feedback
func UpdateFeedback(feedback models.FeedbackRecord) error {
	_, err := db.DB.Exec(`
		UPDATE feedback SET 
		rating = @p1, 
		comments = @p2
		WHERE feedbackId = @p3`,
		feedback.Rating,
		feedback.Comments,
		feedback.FeedbackId)
	if err != nil {
		LogService.LogError("❌ DB error updating feedback: ", err)
		return err
	}
	return nil
}

// DeleteFeedback deletes a feedback by feedbackId
func DeleteFeedback(feedbackId int) error {
	result, err := db.DB.Exec(`
		DELETE FROM feedback
		WHERE feedbackId = @p1`,
		feedbackId)
	if err != nil {
		LogService.LogError("❌ DB error deleting feedback: ", err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ DB error getting rows affected: ", err)
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}



