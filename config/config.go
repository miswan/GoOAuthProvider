package config

const (
    // PasetoKey must be exactly 32 bytes for v2 local
    PasetoKey = "YELLOW_SUBMARINE_TEST_KEY_32_B" // 32 bytes
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
