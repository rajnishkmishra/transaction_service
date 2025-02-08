package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConn struct {
	DB *gorm.DB
}

type Backends struct {
	GinEngine *gin.Engine
	MySQLConn *MySQLConn
}

func SetupAndRun(pathHandlers func(backends Backends)) {
	router := InitializeApplicationConfig(pathHandlers)
	router.Run("0.0.0.0:8080")
}

func NewMysqlConn() *MySQLConn {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	fmt.Printf("dsn: %v\n", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to MySQL: %v", err)
	}

	fmt.Println("Connected to MySQL successfully!")
	return &MySQLConn{
		DB: db,
	}
}

func InitializeApplicationConfig(pathHandlers func(Backends)) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    strings.Split("GET, POST, PUT", ","),
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	router.Use(gin.Recovery())
	mysqlService := NewMysqlConn()

	pathHandlers(Backends{
		GinEngine: router,
		MySQLConn: mysqlService})
	return router
}
