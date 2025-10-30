package businessNotificationType

import "testing"

func TestNotificationMetadata(t *testing.T) {
    for typ, meta := range notificationsData {
        if meta.Name == "" {
            t.Errorf("notification type %v has empty name", typ)
        }
        if meta.Description == "" {
            t.Errorf("notification type %v has empty description", typ)
        }
    }

    if _, ok := notificationsData[NotificationZymm]; !ok {
        t.Fatalf("expected NotificationZymm metadata to exist")
    }
}

