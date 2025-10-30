package notification

import (
    "testing"
    notificationType "zymm/internal/business/notification/notification_type"
)

func TestGetNotificationTitleAndDesc(t *testing.T) {
    cases := []struct {
        notifType notificationType.NotificationType
        data      string
        wantTitle string
        wantDesc  string
    }{
        {notificationType.NotificationZymm, "Update", "Attention Gym Owners/ Managers", "Update"},
        {notificationType.NotificationGeneralInfo, "Info", "Attention!!!", "Info"},
        {notificationType.NotificationGymBroadcast, "Gym", "Attention Members!!!", "Gym"},
        {notificationType.NotificationMembershipReq, "Alice", "Membership Request!", "High Priority: You have a new membership request from Alice, take action now."},
        {notificationType.NotificationMembershipRes, "GymX", "Membership Response!", "GymX management has taken action on your membership request."},
        {notificationType.NotificationFeedbackReceived, "", "Feedback!", "A new feedback was submitted by a user."},
        {notificationType.NotificationPendingFees, "Gold Plan", "Fee Reminder!", "Your Gold Plan membership needs renewal. Please pay to continue."},
    }

    for _, tc := range cases {
        title, desc := GetNotificationTitleAndDesc(tc.notifType, tc.data)
        if title != tc.wantTitle || desc != tc.wantDesc {
            t.Errorf("unexpected values for type %v: got (%s, %s) want (%s, %s)", tc.notifType, title, desc, tc.wantTitle, tc.wantDesc)
        }
    }

    title, desc := GetNotificationTitleAndDesc(999, "")
    if title != "" || desc != "" {
        t.Errorf("expected empty values for unknown notification type, got (%s,%s)", title, desc)
    }
}

func TestGetTitleAndDescForPendingFeesLength(t *testing.T) {
    title, desc := GetTitleAndDescForPendingFees("Platinum")
    if title == "" || desc == "" {
        t.Fatalf("expected non-empty title and description")
    }
    if len(desc) > 100 {
        t.Fatalf("expected description to respect DB length constraint, got len=%d", len(desc))
    }
}

