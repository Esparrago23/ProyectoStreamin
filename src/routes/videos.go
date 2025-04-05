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

// Estructura para la solicitud de carga de video
type VideoUploadRequest struct {
    Title       string `form:"title" binding:"required"`      // El título del video
    Description string `form:"description"`                   // La descripción del video
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

    // Obtener el archivo del request
    file, err := c.FormFile("video")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
        return
    }

    // Generar un nombre de archivo único
    filename := filepath.Join("uploads", strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename)

    // Guardar el archivo
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
        return
    }

    // Crear el registro de video en la base de datos
    video := models.Video{
        Title:       req.Title,
        Description: req.Description,
        FilePath:    filename,
        UserID:      1, // Aquí deberías usar el ID del usuario autenticado
        Status:      "pending",
    }

    if err := DB.Create(&video).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video record"})
        return
    }

    // Crear el registro de procesamiento del video
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

    // Enviar a Kafka para procesamiento
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

// Inicialización de dependencias
var (
    db            *gorm.DB
    kafkaService  *services.KafkaService
)

func InitVideoRoutes(db *gorm.DB, ks *services.KafkaService) {
    DB = db
    kafkaService = ks
}
