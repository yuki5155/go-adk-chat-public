package user

// Role represents the user's role in the system
type Role string

const (
	// RoleUser represents a regular user
	RoleUser Role = "user"
	// RoleRoot represents a root/admin user
	RoleRoot Role = "root"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	return r == RoleUser || r == RoleRoot
}

// IsRoot checks if the role is root
func (r Role) IsRoot() bool {
	return r == RoleRoot
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}
