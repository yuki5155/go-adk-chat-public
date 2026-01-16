package user

// Role represents the user's role in the system
type Role string

const (
	// RoleUser represents a regular user
	RoleUser Role = "user"
	// RoleSubscriber represents a user with chatbot access
	RoleSubscriber Role = "subscriber"
	// RoleAdmin represents an admin user with role management privileges
	RoleAdmin Role = "admin"
	// RoleRoot represents a root/super admin user
	RoleRoot Role = "root"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	return r == RoleUser || r == RoleSubscriber || r == RoleAdmin || r == RoleRoot
}

// IsRoot checks if the role is root
func (r Role) IsRoot() bool {
	return r == RoleRoot
}

// IsAdmin checks if the role is admin or root
func (r Role) IsAdmin() bool {
	return r == RoleAdmin || r == RoleRoot
}

// IsSubscriber checks if the role has subscriber access (subscriber, admin, or root)
func (r Role) IsSubscriber() bool {
	return r == RoleSubscriber || r == RoleAdmin || r == RoleRoot
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}
