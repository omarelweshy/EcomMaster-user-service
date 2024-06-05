package form

type UpdateUserForm struct {
	FirstName string `form:"firstName" binding:"required"`
	LastName  string `form:"lastName" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
}
