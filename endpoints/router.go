package endpoints

import (
    "database/sql"
    "github.com/gin-gonic/gin"
)

var dbConn *sql.DB

// SetupRouter initializes the Gin router and registers all API endpoints.
func SetupRouter(db *sql.DB) *gin.Engine {
    dbConn = db
    router := gin.Default()

    // Register endpoints
    router.GET("/api/corredor", getDrivers)
    router.GET("/api/corredor/detalle/:id", getDriverDetail)
    router.GET("/api/carrera", getCarreras)
    router.GET("/api/carrera/detalle/:id", getCarreraDetail)
    router.GET("/api/resumen", getResumen)

    return router
}
