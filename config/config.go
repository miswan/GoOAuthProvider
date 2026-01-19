package config

// PasetoKey must be exactly 32 bytes long for v2.Local
var PasetoKey = []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // 32 bytes

const (
    AccessTokenExpiry  = 3600 // 1 hour
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
