package domain

// UserRegistration represents the user registration data
type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLogin represents the user login data
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
