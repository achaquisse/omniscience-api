package rest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

const supabaseURL = "https://ovqjkfjuzpxhvrjuagpf.supabase.co"

var publicKey *ecdsa.PublicKey

type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Kid string `json:"kid"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func init() {
	err := fetchSupabasePublicKey()
	if err != nil {
		fmt.Printf("Warning: Failed to fetch Supabase public key: %v\n", err)
	}
}

func fetchSupabasePublicKey() error {
	jwksURL := fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", supabaseURL)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	if len(jwks.Keys) == 0 {
		return fmt.Errorf("no keys found in JWKS")
	}

	key := jwks.Keys[0]
	pubKey, err := jwkToECDSAPublicKey(key)
	if err != nil {
		return err
	}

	publicKey = pubKey
	return nil
}

func jwkToECDSAPublicKey(jwk JWK) (*ecdsa.PublicKey, error) {
	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, err
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, err
	}

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	var curve elliptic.Curve
	switch jwk.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported curve: %s", jwk.Crv)
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

func GetUserEmailFromToken(c *fiber.Ctx) (string, error) {
	if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
		if email, exists := claims["email"].(string); exists && email != "" {
			return email, nil
		}
	}
	return "", fmt.Errorf("unable to extract user email from JWT token")
}

func AuthMiddleware(c *fiber.Ctx) error {
	if c.Method() == "OPTIONS" {
		return c.Next()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ReturnUnauthorized(c, "Missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ReturnUnauthorized(c, "Invalid authorization header format. Use: Bearer <token>")
	}

	var token *jwt.Token
	var err error

	if IsTestMode() {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
	} else {
		if publicKey == nil {
			if err := fetchSupabasePublicKey(); err != nil {
				return ReturnInternalError(c, "Failed to verify token")
			}
		}

		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
	}

	if err != nil || !token.Valid {
		if err != nil {
			log.Error(err)
		}
		return ReturnUnauthorized(c, "Invalid or expired token")
	}

	c.Locals("user", token.Claims)
	return c.Next()
}

func IsTestMode() bool {
	return os.Getenv("TEST_MODE") == "true"
}
