package rolestype

// Role IDs (like enum)
const (
	RoleOwner   = 1
	RoleManager = 2
	RoleStaff   = 3
	RoleTrainer = 4
	RoleMember  = 5
)

// RoleNames maps role IDs to their short names.
var RoleNames = map[int]string{
	RoleOwner:   "owner",
	RoleManager: "manager",
	RoleStaff:   "staff",
	RoleTrainer: "trainer",
	RoleMember:  "member",
}

// RoleDescriptions maps role IDs to detailed descriptions.
var RoleDescriptions = map[int]string{
	RoleOwner:   "Owner of the gym",
	RoleManager: "Manager of the gym with almost all of the rights (less than owner)",
	RoleStaff:   "Staff of the gym with limited rights (less than manager), usually under a manager",
	RoleTrainer: "Trainer of the gym with limited rights, client management and scheduling",
	RoleMember:  "Regular gym member with access to facilities",
}
