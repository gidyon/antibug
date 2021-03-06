package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var (
	signingKey                      = []byte(os.Getenv("JWT_TOKEN"))
	signingMethod jwt.SigningMethod = jwt.SigningMethodHS256
	defaultAPI                      = &authAPI{string(signingKey)}
)

// Payload contains jwt payload
type Payload struct {
	ID           string
	FirstName    string
	LastName     string
	PhoneNumber  string
	EmailAddress string
	Group        string
	Label        string
}

// Claims contains JWT claims information
type Claims struct {
	*Payload
	jwt.StandardClaims
}

// AuthenticateRequest authenticates a request whether it contains valid jwt in metadata
func AuthenticateRequest(ctx context.Context) error {
	return defaultAPI.AuthenticateRequest(ctx)
}

// AuthenticateActor authenticates actor
func AuthenticateActor(ctx context.Context, actorID string) (*Payload, error) {
	return defaultAPI.AuthorizeActor(ctx, actorID)
}

// AuthorizeGroup authorizes an actor group against allowed groups
func AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error) {
	return defaultAPI.AuthorizeGroup(ctx, allowedGroups...)
}

// AuthorizeStrict authenticates and authorizes an actor and group against allowed groups
func AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error) {
	return defaultAPI.AuthorizeStrict(ctx, actorID, allowedGroups...)
}

// GenToken generates jwt
func GenToken(ctx context.Context, payload *Payload, expires int64) (string, error) {
	return defaultAPI.GenToken(ctx, payload, expires)
}

// AddMD adds metadata to token
func AddMD(ctx context.Context, actorID, group string) context.Context {
	return defaultAPI.AddMD(ctx, actorID, group)
}

// ParseToken parses a jwt token and return claims
func ParseToken(tokenString string) (claims *Claims, err error) {
	return defaultAPI.ParseToken(tokenString)
}

// ParseFromCtx jwt token from context
func ParseFromCtx(ctx context.Context) (*Claims, error) {
	return defaultAPI.ParseFromCtx(ctx)
}
