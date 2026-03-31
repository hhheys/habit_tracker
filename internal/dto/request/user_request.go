package request

// UserRole is the role of a user.
type UserRole string

const (
	// RoleUser is the role for a normal user.
	RoleUser UserRole = "user"
	// RoleAdmin is the role for an admin user.
	RoleAdmin UserRole = "admin"
)

// UserRegisterRequest is the request body for user registration.
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLoginRequest is the request body for user login.
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
