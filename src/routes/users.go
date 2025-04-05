package routes

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "streaming-service/src/models"
    "gorm.io/gorm"
	
)

var DB *gorm.DB

// Add this function to initialize the database connection
func InitUserRoutes(db *gorm.DB) {
    DB = db
}

// Update the User struct reference
func registerUser(c *gin.Context) {
type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    user := models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
    }

    if err := DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func loginUser(c *gin.Context) {
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Here you would typically generate a JWT token
    // For now, we'll just return a success message
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}


func SetupUserRoutes(router *gin.Engine) {
    users := router.Group("/api/users")
    {
        users.POST("/register", registerUser)
        users.POST("/login", loginUser)
    }
}