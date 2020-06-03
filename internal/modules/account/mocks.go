package account

import (
	"github.com/gidyon/antibug/internal/pkg/auth"
)

// AuthAPI is auth API
type AuthAPI interface {
	auth.Interface
}
