package messagesRepo

import (
	"database/sql"
	"fmt"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

// CreateMessage inserts a new message
func CreateMessage(req models.CreateMessageRequest, createdBy int) (*models.Message, error) {
	var messageId int
	err := db.DB.QueryRow(`
		INSERT INTO messages (messageForUserId, comment, createdBy)
		OUTPUT INSERTED.messageId
		VALUES (@p1, @p2, @p3)
	`, req.MessageForUserId, req.Comment, createdBy).Scan(&messageId)

	if err != nil {
		LogService.LogError("❌ DB error creating message: ", err)
		return nil, err
	}

	// Fetch the created message
	message, fetchErr := GetMessageById(messageId)
	if fetchErr != nil {
		return nil, fetchErr
	}

	LogService.LogMessage("✅ Message created successfully")
	return message, nil
}

// GetMessageById retrieves a message by its ID
func GetMessageById(messageId int) (*models.Message, error) {
	message := &models.Message{}
	err := db.DB.QueryRow(`
		SELECT messageId, messageForUserId, comment, isActive, isRead, createdBy, createdAt
		FROM messages
		WHERE messageId = @p1
	`, messageId).Scan(
		&message.MessageId,
		&message.MessageForUserId,
		&message.Comment,
		&message.IsActive,
		&message.IsRead,
		&message.CreatedBy,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error fetching message: ", err)
		return nil, err
	}

	return message, nil
}

// UpdateMessage updates a message's content
func UpdateMessage(req models.UpdateMessageRequest, currentUserId int) error {
	result, err := db.DB.Exec(`
		UPDATE messages
		SET comment = @p1
		WHERE messageId = @p2 AND createdBy = @p3 AND isActive = 1
	`, req.Comment, req.MessageId, currentUserId)

	if err != nil {
		LogService.LogError("❌ DB error updating message: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ Error getting rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		LogService.LogMessage("⚠️ No message was updated (might not exist or not owned by user)")
		return sql.ErrNoRows
	}

	LogService.LogMessage("✅ Message updated successfully")
	return nil
}

// SoftDeleteMessage marks a message as inactive (soft delete)
func SoftDeleteMessage(messageId int, currentUserId int) error {
	result, err := db.DB.Exec(`
		UPDATE messages
		SET isActive = 0
		WHERE messageId = @p1 AND createdBy = @p2 AND isActive = 1
	`, messageId, currentUserId)

	if err != nil {
		LogService.LogError("❌ DB error deleting message: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ Error getting rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		LogService.LogMessage("⚠️ No message was deleted (might not exist or not owned by user)")
		return sql.ErrNoRows
	}

	LogService.LogMessage("✅ Message deleted successfully")
	return nil
}

// MarkMessagesAsRead marks specified messages as read
func MarkMessagesAsRead(messageIds []int, currentUserId int) error {
	if len(messageIds) == 0 {
		return nil
	}

	// Build the IN clause for message IDs
	query := `
		UPDATE messages
		SET isRead = 1
		WHERE messageForUserId = @p1 AND isActive = 1 AND messageId IN (`

	params := []interface{}{currentUserId}
	for i, id := range messageIds {
		if i > 0 {
			query += ", "
		}
		query += "@p" + fmt.Sprintf("%d", i+2) // Start from @p2
		params = append(params, id)
	}
	query += ")"

	_, err := db.DB.Exec(query, params...)
	if err != nil {
		LogService.LogError("❌ DB error marking messages as read: ", err)
		return err
	}

	LogService.LogMessage("✅ Messages marked as read")
	return nil
}

// GetChatParticipants returns all users that the current user has chatted with
// Includes last message, unread count, etc.
func GetChatParticipants(currentUserId int) ([]models.ChatParticipant, error) {
	rows, err := db.DB.Query(`
		WITH ChatUsers AS (
			-- Get all users I've sent messages to
			SELECT DISTINCT messageForUserId as otherUserId
			FROM messages
			WHERE createdBy = @p1 AND isActive = 1
			
			UNION
			
			-- Get all users who sent messages to me
			SELECT DISTINCT createdBy as otherUserId
			FROM messages
			WHERE messageForUserId = @p1 AND isActive = 1
		),
		LastMessages AS (
			SELECT 
				CASE 
					WHEN m.createdBy = @p1 THEN m.messageForUserId
					ELSE m.createdBy
				END as otherUserId,
				m.comment as lastMessage,
				m.createdAt as lastMessageTime,
				m.createdBy as lastMessageSentBy,
				ROW_NUMBER() OVER (
					PARTITION BY CASE 
						WHEN m.createdBy = @p1 THEN m.messageForUserId
						ELSE m.createdBy
					END 
					ORDER BY m.createdAt DESC
				) as rn
			FROM messages m
			WHERE (m.createdBy = @p1 OR m.messageForUserId = @p1) AND m.isActive = 1
		),
		UnreadCounts AS (
			SELECT createdBy as otherUserId, COUNT(*) as unreadCount
			FROM messages
			WHERE messageForUserId = @p1 AND isRead = 0 AND isActive = 1
			GROUP BY createdBy
		)
		SELECT 
			u.userId,
			u.userName,
			u.profilePic,
			u.roleId,
			lm.lastMessage,
			lm.lastMessageTime,
			lm.lastMessageSentBy,
			ISNULL(uc.unreadCount, 0) as unreadCount
		FROM ChatUsers cu
		INNER JOIN users u ON u.userId = cu.otherUserId
		LEFT JOIN LastMessages lm ON lm.otherUserId = cu.otherUserId AND lm.rn = 1
		LEFT JOIN UnreadCounts uc ON uc.otherUserId = cu.otherUserId
		WHERE u.isActive = 1
		ORDER BY lm.lastMessageTime DESC
	`, currentUserId)

	if err != nil {
		LogService.LogError("❌ DB error fetching chat participants: ", err)
		return nil, err
	}
	defer rows.Close()

	participants := []models.ChatParticipant{}

	for rows.Next() {
		participant := models.ChatParticipant{}
		err := rows.Scan(
			&participant.UserId,
			&participant.UserName,
			&participant.ProfilePic,
			&participant.RoleId,
			&participant.LastMessageText,
			&participant.LastMessageTime,
			&participant.LastMessageSentBy,
			&participant.UnreadCount,
		)

		if err != nil {
			LogService.LogError("❌ DB scan error fetching chat participants: ", err)
			return nil, err
		}

		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		LogService.LogError("❌ DB rows iteration error fetching chat participants: ", err)
		return nil, err
	}

	return participants, nil
}

// GetChatMessages retrieves all messages between two users and marks as read
func GetChatMessages(currentUserId int, otherUserId int) ([]models.ChatMessage, error) {
	// First, mark messages from the other user as read
	_, err := db.DB.Exec(`
		UPDATE messages
		SET isRead = 1
		WHERE messageForUserId = @p1 AND createdBy = @p2 AND isActive = 1 AND isRead = 0
	`, currentUserId, otherUserId)

	if err != nil {
		LogService.LogError("⚠️ Error marking messages as read: ", err)
		// Continue even if this fails
	}

	// Fetch all messages between the two users
	rows, err := db.DB.Query(`
		SELECT 
			m.messageId,
			m.messageForUserId,
			m.comment,
			m.isActive,
			m.isRead,
			m.createdBy,
			u.userName as createdByName,
			u.profilePic as createdByPic,
			m.createdAt
		FROM messages m
		INNER JOIN users u ON u.userId = m.createdBy
		WHERE m.isActive = 1
		AND (
			(m.createdBy = @p1 AND m.messageForUserId = @p2)
			OR
			(m.createdBy = @p2 AND m.messageForUserId = @p1)
		)
		ORDER BY m.createdAt ASC
	`, currentUserId, otherUserId)

	if err != nil {
		LogService.LogError("❌ DB error fetching chat messages: ", err)
		return nil, err
	}
	defer rows.Close()

	messages := []models.ChatMessage{}

	for rows.Next() {
		message := models.ChatMessage{}
		err := rows.Scan(
			&message.MessageId,
			&message.MessageForUserId,
			&message.Comment,
			&message.IsActive,
			&message.IsRead,
			&message.CreatedBy,
			&message.CreatedByName,
			&message.CreatedByPic,
			&message.CreatedAt,
		)

		if err != nil {
			LogService.LogError("❌ DB scan error fetching chat messages: ", err)
			return nil, err
		}

		// Set helper field
		message.IsSentByMe = message.CreatedBy == currentUserId

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		LogService.LogError("❌ DB rows iteration error fetching chat messages: ", err)
		return nil, err
	}

	return messages, nil
}

// GetUnreadMessageCount returns the total unread message count for a user
func GetUnreadMessageCount(userId int) (int, error) {
	var count int
	err := db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM messages
		WHERE messageForUserId = @p1 AND isRead = 0 AND isActive = 1
	`, userId).Scan(&count)

	if err != nil {
		LogService.LogError("❌ DB error fetching unread count: ", err)
		return 0, err
	}

	return count, nil
}

// GetAvailableChatUsers returns available users to chat with based on role
// - For trainers (roleId = 4): returns all gym members (roleId = 5)
// - For members (roleId = 5): returns all gym trainers (roleId = 4)
func GetAvailableChatUsers(currentUserId int, roleId int) ([]models.AvailableChatUser, error) {
	var users []models.AvailableChatUser

	// First, get the user's gym ID
	var gymId int
	var userRoleId int
	err := db.DB.QueryRow(`
		SELECT u.roleId, COALESCE(um.gymId, e.gymId) as gymId
		FROM users u
		LEFT JOIN (
			SELECT userId, (SELECT TOP 1 p.gymId FROM plans p WHERE p.planId = userMemberships.planId) as gymId
			FROM userMemberships
			WHERE isActive = 1 AND membershipStatus = 'A'
		) um ON u.userId = um.userId
		LEFT JOIN employees e ON u.userId = e.userId
		WHERE u.userId = @p1
	`, currentUserId).Scan(&userRoleId, &gymId)

	if err != nil {
		LogService.LogError("❌ DB error fetching user's gym: ", err)
		return nil, err
	}

	if gymId == 0 {
		LogService.LogMessage("⚠️ User has no gym association")
		return users, nil // Return empty list
	}

	// Based on role, fetch appropriate users
	var rows *sql.Rows
	if roleId == 4 { // Trainer - fetch all members
		rows, err = db.DB.Query(`
			select 
                u.userId,
				u.userName,
				u.mobile,
				u.profilePic,
				u.roleId
                from users u
                where u.userId in 
            (SELECT distinct
				u1.userId
			FROM userMemberships um
			INNER JOIN users u1 ON u1.userId = um.userId
			INNER JOIN plans p ON p.planId = um.planId
			WHERE p.gymId = @p1
			AND um.isActive = 1
			AND um.membershipStatus = 'A'
			AND um.startDate <= GETDATE()
			AND um.endDate >= GETDATE()
			AND u1.roleId = 5)
		`, gymId, currentUserId)

		if err != nil {
			LogService.LogError("❌ DB error fetching members: ", err)
			return nil, err
		}

	} else if roleId == 5 { // Member - fetch all trainers
		rows, err = db.DB.Query(`
			SELECT
				u.userId,
				u.userName,
				u.mobile,
				u.profilePic,
				u.roleId
			FROM users u
			INNER JOIN employees e ON u.userId = e.userId
			WHERE e.gymId = @p1
			AND u.roleId = 4
			AND u.isActive = 1
			AND u.userId != @p2
			ORDER BY u.userName ASC
		`, gymId, currentUserId)

		if err != nil {
			LogService.LogError("❌ DB error fetching trainers: ", err)
			return nil, err
		}
	} else {
		LogService.LogMessage(fmt.Sprintf("⚠️ Chat not available for role %d", roleId))
		return users, nil // Return empty list for other roles
	}

	if rows != nil {
		defer rows.Close()

		for rows.Next() {
			var user models.AvailableChatUser
			err := rows.Scan(
				&user.UserId,
				&user.UserName,
				&user.Mobile,
				&user.ProfilePic,
				&user.RoleId,
			)
			if err != nil {
				LogService.LogError("❌ DB scan error: ", err)
				continue
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			LogService.LogError("❌ DB rows error: ", err)
			return nil, err
		}
	}

	LogService.LogMessage(fmt.Sprintf("✅ Found %d available chat users for role %d", len(users), roleId))
	return users, nil
}
