package businessRoles

import businessRoles "zymm/internal/business/roles/roles_type"

func IsAllowedToTakeActionOnMemberships(roleType businessRoles.RoleType) bool {
	return roleType == businessRoles.RoleOwner || roleType == businessRoles.RoleManager
}

func IsNotAllowedToTakeActionOnMemberships(roleType businessRoles.RoleType) bool {
	return !IsAllowedToTakeActionOnMemberships(roleType)
}

func IsAllowedToManagePlans(roleType businessRoles.RoleType) bool {
	return roleType == businessRoles.RoleOwner || roleType == businessRoles.RoleManager
}

func IsNotAllowedToManagePlans(roleType businessRoles.RoleType) bool {
	return !IsAllowedToManagePlans(roleType)
}

func IsAllowedToViewAttendance(roleType businessRoles.RoleType) bool {
	return roleType == businessRoles.RoleOwner || roleType == businessRoles.RoleManager || roleType == businessRoles.RoleStaff || roleType == businessRoles.RoleTrainer
}

func IsNotAllowedToViewAttendance(roleType businessRoles.RoleType) bool {
	return !IsAllowedToViewAttendance(roleType)
}
