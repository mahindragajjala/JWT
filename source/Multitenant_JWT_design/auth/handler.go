//Login & Issue JWT
package auth

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "multitenant-jwt/util"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    TenantID string `json:"tenant_id"`
}

func LoginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Simulate user check
    if req.Username != "alice" || req.Password != "pass123" || req.TenantID != "alpha" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    token, err := util.GenerateJWT(req.Username, req.TenantID, "admin")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
