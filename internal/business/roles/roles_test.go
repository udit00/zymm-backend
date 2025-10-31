package businessRoles

import (
    "testing"
    roleTypes "zymm/internal/business/roles/roles_type"
)

func TestMembershipActionPermissions(t *testing.T) {
    cases := []struct {
        role     roleTypes.RoleType
        allowed  bool
    }{
        {roleTypes.RoleOwner, true},
        {roleTypes.RoleManager, true},
        {roleTypes.RoleStaff, false},
        {roleTypes.RoleTrainer, false},
        {roleTypes.RoleMember, false},
    }

    for _, tc := range cases {
        if IsAllowedToTakeActionOnMemberships(tc.role) != tc.allowed {
            t.Errorf("role %v expected allowed=%t", tc.role, tc.allowed)
        }
        if IsNotAllowedToTakeActionOnMemberships(tc.role) == tc.allowed {
            t.Errorf("role %v expected inverse permission", tc.role)
        }
    }
}

func TestPlanManagementPermissions(t *testing.T) {
    cases := []struct {
        role    roleTypes.RoleType
        allowed bool
    }{
        {roleTypes.RoleOwner, true},
        {roleTypes.RoleManager, true},
        {roleTypes.RoleStaff, false},
        {roleTypes.RoleTrainer, false},
        {roleTypes.RoleMember, false},
    }

    for _, tc := range cases {
        if IsAllowedToManagePlans(tc.role) != tc.allowed {
            t.Errorf("role %v expected allowed=%t", tc.role, tc.allowed)
        }
        if IsNotAllowedToManagePlans(tc.role) == tc.allowed {
            t.Errorf("role %v expected inverse permission", tc.role)
        }
    }
}

func TestAttendancePermissions(t *testing.T) {
    allowedRoles := []roleTypes.RoleType{roleTypes.RoleOwner, roleTypes.RoleManager, roleTypes.RoleStaff, roleTypes.RoleTrainer}
    for _, role := range allowedRoles {
        if !IsAllowedToViewAttendance(role) {
            t.Errorf("expected role %v to be allowed", role)
        }
        if IsNotAllowedToViewAttendance(role) {
            t.Errorf("expected role %v to be allowed (inverse)", role)
        }
    }

    if IsAllowedToViewAttendance(roleTypes.RoleMember) {
        t.Errorf("members should not be allowed to view attendance")
    }
    if !IsNotAllowedToViewAttendance(roleTypes.RoleMember) {
        t.Errorf("members should be denied to view attendance")
    }
}

func TestEmployeeManagementPermissions(t *testing.T) {
    cases := []struct {
        role    roleTypes.RoleType
        allowed bool
    }{
        {roleTypes.RoleOwner, true},
        {roleTypes.RoleManager, true},
        {roleTypes.RoleStaff, false},
        {roleTypes.RoleTrainer, false},
        {roleTypes.RoleMember, false},
    }

    for _, tc := range cases {
        if IsAllowedToManageEmployee(tc.role) != tc.allowed {
            t.Errorf("role %v expected allowed=%t", tc.role, tc.allowed)
        }
        if IsNotAllowedToManageEmployees(tc.role) == tc.allowed {
            t.Errorf("role %v expected inverse permission", tc.role)
        }
    }
}

