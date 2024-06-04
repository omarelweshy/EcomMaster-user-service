package form

type RegistrationForm struct {
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
	Username  string `form:"username" binding:"required"`
	Password  string `form:"password" binding:"required,min=8"`
	Email     string `form:"email" binding:"required,email"`
}
