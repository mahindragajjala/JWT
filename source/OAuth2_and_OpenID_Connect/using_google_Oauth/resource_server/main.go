// resource_server/main.go

package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "strings"
    "github.com/coreos/go-oidc/v3/oidc"
)

var (
    verifier *oidc.IDTokenVerifier
)

func init() {
    ctx := context.Background()
    provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
    if err != nil {
        log.Fatalf("Failed to get provider: %v", err)
    }
    verifier = provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")})
}

func main() {
    http.HandleFunc("/api/data", handleData)
    log.Println("Resource server running on :9000")
    log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleData(w http.ResponseWriter, r *http.Request) {
    auth := r.Header.Get("Authorization")
    if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
        http.Error(w, "Bearer token required", http.StatusUnauthorized)
        return
    }

    token := strings.TrimPrefix(auth, "Bearer ")
    idToken, err := verifier.Verify(context.Background(), token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    var claims map[string]interface{}
    if err := idToken.Claims(&claims); err != nil {
        http.Error(w, "failed to parse claims", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Protected data! Hello %s", claims["email"])
}
