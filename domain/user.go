package domain

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}

type UserService interface {
	SignUp(user *User) error
}
