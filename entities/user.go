package entities

type User struct {
	Id        int64
	Name      string `validate:"required" label:"Full Name"`
	Email     string `validate:"required,email"`
	Username  string `validate:"required,gte=3"`
	Password  string `validate:"required,gte=6"`
	CPassword string `validate:"required,eqfield=Password" label:"Password Confirmation"`
}
