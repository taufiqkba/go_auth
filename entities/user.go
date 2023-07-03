package entities

type User struct {
	Id        int64
	Name      string `validate:"required" label:"Full Name"`
	Email     string `validate:"required,email,is_unique=users-email"`
	Username  string `validate:"required,gte=3,is_unique=users-username"`
	Password  string `validate:"required,gte=6"`
	CPassword string `validate:"required,eqfield=Password" label:"Password Confirmation"`
}
