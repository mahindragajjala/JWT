package main

import (
    "github.com/gin-gonic/gin"
    "multitenant-jwt/auth"
    "multitenant-jwt/middleware"
)

func main() {
    r := gin.Default()

    // Public login route
    r.POST("/login", auth.LoginHandler)

    // Protected tenant API group
    tenant := r.Group("/api")
    tenant.Use(middleware.TenantJWTMiddleware())
    {
        tenant.GET("/dashboard", func(c *gin.Context) {
            claims, _ := c.Get("claims")
            c.JSON(200, gin.H{
                "message": "Welcome to the tenant dashboard!",
                "claims":  claims,
            })
        })
    }

    r.Run(":8080")
}
