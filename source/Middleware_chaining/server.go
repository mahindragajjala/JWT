/*
    Middleware is a function that wraps the HTTP request/response lifecycle, 
    allowing you to run logic before and/or after the request is handled 
    by the final handler.
    
    Chaining means executing multiple middleware functions in order, 
    where each middleware calls the next one in the chain.


                            Incoming Request
                               ↓
                            [Middleware 1]
                               ↓
                            [Middleware 2]
                               ↓
                            [Handler]
                               ↓
                            Response


Context Propagation
Go’s context.Context is used to:
  - Carry request-scoped data across API boundaries and goroutines
  - Handle timeouts, cancellation, and deadlines
  - Inject values that can be retrieved later (e.g., user ID, correlation ID, tenant ID)
    Middleware can attach data to context.Context so that downstream handlers or 
    middleware can access it.
*/
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
)

type key string

const userKey key = "user"

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        log.Printf("Completed in %v", time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := r.Header.Get("X-User")
        if user == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        // Store in context
        ctx := context.WithValue(r.Context(), userKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value(userKey)
    fmt.Fprintf(w, "Hello, %s!", user)
}

func main() {
    final := http.HandlerFunc(finalHandler)

    // Chain middlewares
    http.Handle("/", loggingMiddleware(authMiddleware(final)))

    log.Println("Starting on :8080")
    http.ListenAndServe(":8080", nil)
}


/*
What This Does:
            loggingMiddleware: Logs the request and response time
            authMiddleware: Checks if X-User header exists, adds it to the context
            finalHandler: Retrieves the user from context and responds
*/

/*
using the gin framework
                          r := gin.Default()
                          
                          r.Use(func(c *gin.Context) {
                              // Add a value to context
                              c.Set("RequestID", "abc-123")
                              c.Next()
                          })
                          
                          r.GET("/", func(c *gin.Context) {
                              rid := c.GetString("RequestID")
                              c.JSON(200, gin.H{"message": "hello", "request_id": rid})
                          })
*/
/*
Tips for Middleware Chaining
 Feature                Description                                               
 `context.WithValue`    Store values in the request context                       
 `context.WithTimeout`  Add timeout control                                       
 `c.Next()`             In Gin, continues to next middleware                      
 `r.WithContext()`      In `net/http`, propagate modified context to next handler 

*/
/*
use cases
    ✅ Authorization: Extract token, validate it, and store user ID in context
    ✅ Request ID: Assign unique ID per request for logging/tracing
    ✅ Localization: Read Accept-Language and inject into context
    ✅ Tenant Handling: Multi-tenant apps use middleware to resolve and set tenant ID
*/
