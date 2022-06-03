package signature

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func AllowedSignMethods() []string {
	return []string{
		jwt.SigningMethodHS256.Alg(),
		jwt.SigningMethodHS384.Alg(),
		jwt.SigningMethodHS512.Alg(),
	}
}

type SignerKey func(signer string) (key []byte, err error)

// ParseToken extracts fields from the given JWT and validates them
func ParseToken(jwtToken string, signerKey SignerKey) (a JWTPayload, err error) {
	parser := jwt.Parser{
		ValidMethods:         AllowedSignMethods(),
		SkipClaimsValidation: false,
	}

	// extract claims but do not validate token for now
	var claims JWTClaims
	t, err := parser.ParseWithClaims(jwtToken, &claims, func(t *jwt.Token) (any, error) {
		// issuer is project login, check it
		if !signerRex.MatchString(claims.Issuer) {
			return nil, fmt.Errorf("invalid signer/issuer format")
		}
		return signerKey(claims.Issuer)
	})
	if err != nil {
		err = fmt.Errorf("parsing token: %w", err)
		return
	}

	if !t.Valid {
		err = fmt.Errorf("suddenly jwt is invalid")
		return
	}

	// extract payload
	payload := JWTPayload{
		Method:      claims.Method,
		URL:         claims.URL,
		Nonce:       claims.ID,
		Signer:      claims.Issuer,
		Subject:     claims.Subject,
		Recipients:  claims.Audience,
		BodyHash:    claims.BodyHash,
		BodyHashAlg: claims.BodyHashAlg,
	}

	if err = newValidator().Struct(payload); err != nil {
		err = fmt.Errorf("invalid jwt payload: %w", err)
		return
	}

	a = payload
	return
}
