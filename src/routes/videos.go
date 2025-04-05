package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "path/filepath"
    "strconv"
    "time"
   
	"streaming-service/src/services"
    "streaming-service/src/models"
    "gorm.io/gorm"
)

type VideoUploadRequest struct {
    Title       string `form:"title" binding:"required"`
    Description string `form:"description"`
}

func SetupVideoRoutes(router *gin.Engine) {
    videos := router.Group("/api/videos")
    {
        videos.POST("/upload", uploadVideo)
        videos.GET("/:id", getVideo)
        videos.GET("/", listVideos)
    }
}

func uploadVideo(c *gin.Context) {
    var req VideoUploadRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get the file from the request
    file, err := c.FormFile("video")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
        return
    }

    // Generate unique filename
    filename := filepath.Join("uploads", strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename)

    // Save the file
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
        return
    }

    // Create video record in database
    video := models.Video{
        Title:       req.Title,
        Description: req.Description,
        FilePath:    filename,
        UserID:      1, // Replace with actual user ID from authentication
        Status:      "pending",
    }

    if err := DB.Create(&video).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video record"})
        return
    }

    // Create video processing record
    videoProcessing := models.VideoProcessing{
        VideoID:    video.ID,
        KafkaTopic: "video-processing",
        Partition:  0,
        Status:     "queued",
    }

    if err := DB.Create(&videoProcessing).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create processing record"})
        return
    }

    // Send to Kafka for processing
    if err := kafkaService.SendVideoForProcessing(video.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue video for processing"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Video uploaded successfully",
        "video":   video,
    })
}

func getVideo(c *gin.Context) {
    id := c.Param("id")
    var video models.Video

    if err := DB.First(&video, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }

    c.JSON(http.StatusOK, video)
}

func listVideos(c *gin.Context) {
    var videos []models.Video
    
    if err := DB.Find(&videos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
        return
    }

    c.JSON(http.StatusOK, videos)
}

// Add these variables at package level
var (
    db *gorm.DB
    kafkaService *services.KafkaService
)

// Add this function to initialize the dependencies
func InitVideoRoutes(db *gorm.DB, ks *services.KafkaService) {
    DB = db
    kafkaService = ks
}