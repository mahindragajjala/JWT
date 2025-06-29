package util

import (
    "errors"
    "io/ioutil"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var tenantKeyMap = map[string]string{
    "alpha": "keys/alpha_private.pem",
}

var tenantPubMap = map[string]string{
    "alpha": "keys/alpha_public.pem",
}

func GenerateJWT(username, tenantID, role string) (string, error) {
    privateKeyData, err := ioutil.ReadFile(tenantKeyMap[tenantID])
    if err != nil {
        return "", err
    }

    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
    if err != nil {
        return "", err
    }

    claims := jwt.MapClaims{
        "sub":       username,
        "tenant_id": tenantID,
        "role":      role,
        "exp":       time.Now().Add(15 * time.Minute).Unix(),
        "iss":       "auth." + tenantID + ".com",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    token.Header["kid"] = tenantID + "-key"

    return token.SignedString(privateKey)
}

func ValidateJWT(tokenString string, tenantID string) (jwt.MapClaims, error) {
    pubKeyPath := tenantPubMap[tenantID]
    pubKeyData, err := ioutil.ReadFile(pubKeyPath)
    if err != nil {
        return nil, err
    }

    pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
    if err != nil {
        return nil, err
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return pubKey, nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
