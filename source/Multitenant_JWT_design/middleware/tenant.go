//JWT & Tenant Validation
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "multitenant-jwt/util"
)

func TenantJWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        // Extract tenant ID from query or header for routing
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing tenant ID"})
            return
        }

        claims, err := util.ValidateJWT(tokenStr, tenantID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        if claims["tenant_id"] != tenantID {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant mismatch"})
            return
        }

        c.Set("claims", claims)
        c.Next()
    }
}
