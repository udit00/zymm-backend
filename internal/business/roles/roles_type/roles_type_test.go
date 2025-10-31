package businessRoleType

import "testing"

func TestGetRoleTypeFromInt(t *testing.T) {
    for role, name := range RoleNames {
        got := GetRoleTypeFromInt(role.Int())
        if got == nil {
            t.Fatalf("expected role %v to be resolved", role)
        }
        if RoleNames[*got] != name {
            t.Fatalf("expected role name %s, got %s", name, RoleNames[*got])
        }
    }

    if got := GetRoleTypeFromInt(999); got != nil {
        t.Fatalf("expected nil for unknown role, got %v", *got)
    }
}

func TestRoleTypeInt(t *testing.T) {
    if RoleOwner.Int() != int(RoleOwner) {
        t.Fatalf("expected Int() to return underlying int for RoleOwner")
    }
}

