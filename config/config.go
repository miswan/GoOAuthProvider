package config

const (
    JWTSecret = "your-secret-key-here-must-be-32b" // 32 bytes for Paseto v2
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
