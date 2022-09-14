package model

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
type RegisterUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

type UpdatePasswordRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
