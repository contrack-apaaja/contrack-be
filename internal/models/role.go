package models

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleRegular    UserRole = "REGULAR"
	RoleLegal      UserRole = "LEGAL"
	RoleManagement UserRole = "MANAGEMENT"
)

// IsValid checks if the role is valid
func (r UserRole) IsValid() bool {
	return r == RoleRegular || r == RoleLegal || r == RoleManagement
}

// CanUpdateContracts checks if the role can update contracts and clauses
func (r UserRole) CanUpdateContracts() bool {
	return r == RoleLegal || r == RoleManagement
}

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// GetDefaultRole returns the default role for new users
func GetDefaultRole() UserRole {
	return RoleRegular
}
