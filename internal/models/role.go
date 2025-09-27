package models

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleLegal      UserRole = "legal"
	RoleManagement UserRole = "management"
)

// IsValid checks if the role is valid
func (r UserRole) IsValid() bool {
	return r == RoleUser || r == RoleLegal || r == RoleManagement
}

// CanAccessLegalReview checks if the role can access legal review functions
func (r UserRole) CanAccessLegalReview() bool {
	return r == RoleLegal || r == RoleManagement
}

// CanAccessManagementApproval checks if the role can access management approval functions
func (r UserRole) CanAccessManagementApproval() bool {
	return r == RoleManagement
}

// CanAccessBasicFeatures checks if the role can access contracts, clauses, dashboard
func (r UserRole) CanAccessBasicFeatures() bool {
	return r == RoleUser || r == RoleLegal || r == RoleManagement
}

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// GetDefaultRole returns the default role for new users
func GetDefaultRole() UserRole {
	return RoleUser
}
