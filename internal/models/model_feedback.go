package models

import "time"

// type FeedbackRecord struct {
// 	FeedbackId int       `json:"feedbackId"`
// 	Rating     *int      `json:"rating"`
// 	Comments   *string   `json:"comments"`
// 	GymId      int       `json:"gymId"`
// 	CreatedBy  int       `json:"createdBy"`
// 	CreatedAt  time.Time `json:"createdAt"`
// }

type CreateFeedbackRequest struct {
	Rating   *int    `json:"rating"`
	Comments *string `json:"comments"`
	GymId    int     `json:"gymId"`
}

type UpdateFeedbackRequest struct {
	FeedbackId int     `json:"feedbackId"`
	Rating     *int    `json:"rating"`
	Comments   *string `json:"comments"`
}

type FeedbackRecord struct {
	FeedbackId          int       `json:"feedbackId"`
	Rating              *int      `json:"rating"`
	Comments            *string   `json:"comments"`
	GymId               int       `json:"gymId"`
	CreatedBy           int       `json:"createdBy"`
	CreatedByName       *string   `json:"createdByName"`
	CreatedByProfilePic *string   `json:"profilePic"`
	CreatedAt           time.Time `json:"createdAt"`
}
