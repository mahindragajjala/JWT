/*
[CLIENT]                           [SERVER]
   |                                   |
   |---- Login with credentials -----> |
   |                                   |
   | <- Access Token + Refresh Token --|
   |                                   |
   |---- API request with Access ----> |  ✅ valid
   |                                   |
   | <- Protected Resource Response ---|
   |
   |---- API request with Access ----> |  ❌ expired
   | <- 401 Unauthorized --------------|
   |
   |---- Request New Token using ----> |
   |     Refresh Token                 |
   |                                   |
   | <- New Access Token (and maybe) --|
   |    new Refresh Token              |
   |
   |---- Retry API request with ------>|
   |     New Access Token              |
   |                                   |
   | <- Protected Resource Response ---|

*/
func RefreshTokenHandler(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 1. Parse and validate the refresh token
    token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
        return []byte("refresh-secret"), nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

    claims := token.Claims.(jwt.MapClaims)
    userID := claims["sub"].(string)

    // 2. (Optional) Check in DB if this refresh token is still valid

    // 3. Create a new access token
    newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(15 * time.Minute).Unix(),
    })

    accessTokenString, _ := newAccessToken.SignedString([]byte("access-secret"))

    c.JSON(http.StatusOK, gin.H{
        "access_token": accessTokenString,
    })
}
