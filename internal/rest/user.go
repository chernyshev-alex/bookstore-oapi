package rest

type UserHandler struct {
	//booksSpi service.BooksService
}

var _ ServerInterface = (NewBooksHandler)(nil)

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// func (h *UserHandler) RegisterUser(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		c.Abort()
// 		return
// 	}
// 	if err := user.HashPassword(user.Password); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		c.Abort()
// 		return
// 	}
// 	record := database.Instance.Create(&user)
// 	if record.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
// 		c.Abort()
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
// }
