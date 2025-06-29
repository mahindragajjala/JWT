// auth_client/main.go

package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var (
    clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
    clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
    redirectURL  = "http://localhost:8080/callback"
    ctx          = context.Background()
)

func main() {
    provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
    if err != nil {
        log.Fatalf("failed to get provider: %v", err)
    }

    state := "example-state" // In production, use secure random state
    oauth2Config := oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        Endpoint:     provider.Endpoint(),
        RedirectURL:  redirectURL,
        Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
    }

    verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        link := oauth2Config.AuthCodeURL(state)
        http.Redirect(w, r, link, http.StatusFound)
    })

    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("state") != state {
            http.Error(w, "state mismatch", http.StatusBadRequest)
            return
        }
        oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
        if err != nil {
            http.Error(w, "token exchange failed", http.StatusInternalServerError)
            return
        }

        rawIDToken, ok := oauth2Token.Extra("id_token").(string)
        if !ok {
            http.Error(w, "no id_token", http.StatusInternalServerError)
            return
        }

        idToken, err := verifier.Verify(ctx, rawIDToken)
        if err != nil {
            http.Error(w, "invalid id token", http.StatusUnauthorized)
            return
        }

        var claims struct {
            Email         string `json:"email"`
            EmailVerified bool   `json:"email_verified"`
            Name          string `json:"name"`
            Picture       string `json:"picture"`
        }
        if err := idToken.Claims(&claims); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Hello %s!\nID Token claims: %+v\nAccess Token: %s", claims.Name, claims, oauth2Token.AccessToken)
    })

    log.Println("Client app running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
