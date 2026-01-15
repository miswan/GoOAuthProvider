package config

const (
    // Must be 32 bytes for Paseto v2
    JWTSecret = "your-secret-key-here-12345678901"
    AccessTokenExpiry = 3600 // 1 hour
    RefreshTokenExpiry = 7200 // 2 hours
)

type OAuth2Config struct {
    AuthorizeEndpoint string
    TokenEndpoint     string
    UserInfoEndpoint  string
}

var DefaultConfig = OAuth2Config{
    AuthorizeEndpoint: "/authorize",
    TokenEndpoint:     "/token",
    UserInfoEndpoint:  "/userinfo",
}
