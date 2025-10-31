package businessRoleType

type RoleType int

// Role IDs (like enum)
const (
	RoleOwner   RoleType = 1
	RoleManager RoleType = 2
	RoleStaff   RoleType = 3
	RoleTrainer RoleType = 4
	RoleMember  RoleType = 5
)

// RoleNames maps role IDs to their short names.
var RoleNames = map[RoleType]string{
	RoleOwner:   "owner",
	RoleManager: "manager",
	RoleStaff:   "staff",
	RoleTrainer: "trainer",
	RoleMember:  "member",
}

// RoleDescriptions maps role IDs to detailed descriptions.
var RoleDescriptions = map[RoleType]string{
	RoleOwner:   "Owner of the gym",
	RoleManager: "Manager of the gym with almost all of the rights (less than owner)",
	RoleStaff:   "Staff of the gym with limited rights (less than manager), usually under a manager",
	RoleTrainer: "Trainer of the gym with limited rights, client management and scheduling",
	RoleMember:  "Regular gym member with access to facilities",
}

func GetRoleTypeFromInt(roleTypeInt int) *RoleType {
	role := RoleType(roleTypeInt)
	if _, exists := RoleNames[role]; exists {
		return &role
	}
	return nil
}

func (r RoleType) Int() int {
	return int(r)
}
