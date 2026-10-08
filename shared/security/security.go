package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// GenerateRandomBytes generates cryptographically secure random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// GenerateRandomString generates a random string of the given length.
func GenerateRandomString(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateAPIKey generates a secure API key.
func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword checks a password against a bcrypt hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// DeriveKey derives a key from a password using scrypt.
func DeriveKey(password, salt string) ([]byte, error) {
	return scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
}

// GenerateSalt generates a random salt.
func GenerateSalt() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashWithSalt hashes a value with a salt using HMAC-SHA256.
func HashWithSalt(value, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC verifies an HMAC.
func VerifyHMAC(value, salt, expectedHash string) bool {
	computed := HashWithSalt(value, salt)
	return subtle.ConstantTimeCompare([]byte(computed), []byte(expectedHash)) == 1
}

// JWTClaims represents JWT claims.
type JWTClaims struct {
	Subject   string   `json:"sub"`
	TenantID  string   `json:"tenant_id,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
	Issuer    string   `json:"iss,omitempty"`
	Audience  string   `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token.
func GenerateJWT(claims *JWTClaims, signingKey []byte) (string, error) {
	if claims == nil {
		return "", fmt.Errorf("claims cannot be nil")
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   claims.Subject,
		IssuedAt:  jwt.NewNumericDate(time.Unix(claims.IssuedAt, 0)),
		ExpiresAt: jwt.NewNumericDate(time.Unix(claims.ExpiresAt, 0)),
		Issuer:    claims.Issuer,
		Audience:  jwt.ClaimStrings{claims.Audience},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signingKey)
}

// ParseJWT parses and validates a JWT token.
func ParseJWT(tokenString string, verificationKey []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verificationKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateJWT validates a JWT token without parsing claims.
func ValidateJWT(tokenString string, verificationKey []byte) error {
	_, err := ParseJWT(tokenString, verificationKey)
	return err
}

// ExtractBearerToken extracts Bearer token from Authorization header.
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return parts[1], nil
}

// SanitizeForLog removes sensitive data from strings for logging.
func SanitizeForLog(input string) string {
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"key",
		"credential",
		"authorization",
		"api_key",
		"apikey",
		"access_token",
		"refresh_token",
		"client_secret",
	}

	result := input
	for _, pattern := range sensitivePatterns {
		// Case-insensitive replace
		result = replaceCaseInsensitive(result, pattern, "[REDACTED]")
	}
	return result
}

func replaceCaseInsensitive(s, old, replacement string) string {
	if old == "" {
		return s
	}
	var result strings.Builder
	lowerS := strings.ToLower(s)
	lowerOld := strings.ToLower(old)
	index := 0
	for {
		i := strings.Index(lowerS[index:], lowerOld)
		if i == -1 {
			result.WriteString(s[index:])
			break
		}
		result.WriteString(s[index : index+i])
		result.WriteString(replacement)
		index += i + len(old)
	}
	return result.String()
}

// ConstantTimeCompare compares two strings in constant time.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// GenerateID generates a unique ID.
func GenerateID() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateExecutionID generates an execution ID (UUID v4).
func GenerateExecutionID() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	// Set version 4 (random)
	b[6] = (b[6] & 0x0f) | 0x40
	// Set variant
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// GenerateCorrelationID generates a correlation ID.
func GenerateCorrelationID() (string, error) {
	return GenerateExecutionID()
}

// GenerateTraceID generates a W3C trace ID (32 hex chars).
func GenerateTraceID() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateSpanID generates a W3C span ID (16 hex chars).
func GenerateSpanID() (string, error) {
	b, err := GenerateRandomBytes(8)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// MaskSecret masks a secret for display.
func MaskSecret(secret string) string {
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}

// ValidateSecretFormat validates a secret reference format.
func ValidateSecretFormat(secret string) bool {
	return strings.HasPrefix(secret, "op://") || strings.HasPrefix(secret, "vault:")
}

// ParseSecretRef parses a secret reference.
func ParseSecretRef(secret string) (scheme, path string, ok bool) {
	if strings.HasPrefix(secret, "op://") {
		return "op", strings.TrimPrefix(secret, "op://"), true
	}
	if strings.HasPrefix(secret, "vault:") {
		return "vault", strings.TrimPrefix(secret, "vault:"), true
	}
	return "", "", false
}
