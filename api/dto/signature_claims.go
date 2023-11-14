package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type SignatureClaims struct {
	PayloadChecksum string `json:"cs"`

	jwt.RegisteredClaims
}
