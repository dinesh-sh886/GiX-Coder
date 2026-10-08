package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

// Client handles JWT validation and JWKS management.
type Client struct {
	jwksURL    string
	issuer     string
	audience   string
	algorithm  string
	keyCache   *sync.Map
	jwksClient *http.Client
}

// Config holds authentication configuration.
type Config struct {
	JWKSURL   string
	Issuer    string
	Audience  string
	Algorithm string
}

// NewClient creates a new authentication client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.JWKSURL == "" {
		return nil, fmt.Errorf("JWKS URL is required")
	}
	if cfg.Issuer == "" {
		return nil, fmt.Errorf("issuer is required")
	}
	if cfg.Audience == "" {
		return nil, fmt.Errorf("audience is required")
	}
	if cfg.Algorithm == "" {
		cfg.Algorithm = "RS256"
	}

	return &Client{
		jwksURL:   cfg.JWKSURL,
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
		algorithm: cfg.Algorithm,
		keyCache:  &sync.Map{},
		jwksClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// JWKS represents a JSON Web Key Set.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

// Claims represents JWT claims.
type Claims struct {
	Subject   string   `json:"sub"`
	TenantID  string   `json:"tenant_id,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
	Issuer    string   `json:"iss,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens.
func (c *Client) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := c.ExtractBearerToken(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": err.Error(),
				},
			})
			return
		}

		claims, err := c.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": err.Error(),
				},
			})
			return
		}

		// Store claims in context
		ctx.Set("claims", claims)
		ctx.Set("user_id", claims.Subject)
		ctx.Set("tenant_id", claims.TenantID)
		ctx.Set("scopes", claims.Scopes)

		ctx.Next()
	}
}

// ValidateToken validates a JWT token and returns claims.
func (c *Client) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, c.keyFunc)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Validate issuer
		if c.issuer != "" && claims.Issuer != c.issuer {
			return nil, fmt.Errorf("invalid issuer: expected %s, got %s", c.issuer, claims.Issuer)
		}

		// Validate audience
		if c.audience != "" {
			found := false
			for _, aud := range claims.Audience {
				if aud == c.audience {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("invalid audience: expected %s, not found in %v", c.audience, claims.Audience)
			}
		}

		// Validate expiry
		if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
			return nil, fmt.Errorf("token expired")
		}

		// Validate not before
		if claims.IssuedAt > 0 && time.Now().Unix() < claims.IssuedAt {
			return nil, fmt.Errorf("token not yet valid")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// keyFunc returns the key function for JWT parsing.
func (c *Client) keyFunc(token *jwt.Token) (interface{}, error) {
	// Verify signing method
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing kid in token header")
	}

	// Get key from cache
	if key, ok := c.keyCache.Load(kid); ok {
		return key, nil
	}

	// Fetch JWKS
	jwks, err := c.fetchJWKS()
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}

	// Find matching key
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			pubKey, err := c.parseRSAPublicKey(key)
			if err != nil {
				return nil, fmt.Errorf("parse RSA key: %w", err)
			}
			c.keyCache.Store(kid, pubKey)
			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("key not found: %s", kid)
}

// fetchJWKS fetches the JWKS from the configured URL.
func (c *Client) fetchJWKS() (*JWKS, error) {
	resp, err := c.jwksClient.Get(c.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}

	return &jwks, nil
}

// parseRSAPublicKey parses an RSA public key from a JWK.
func (c *Client) parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("decode n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("decode e: %w", err)
	}

	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

// extractBearerToken extracts the Bearer token from the Authorization header.
func (c *Client) ExtractBearerToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

// GetClaims extracts claims from the request context.
func GetClaims(c *gin.Context) (*Claims, bool) {
	val, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	claims, ok := val.(*Claims)
	return claims, ok
}

// GetUserID extracts the user ID from the context.
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetTenantID extracts the tenant ID from the context.
func GetTenantID(c *gin.Context) (string, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return "", false
	}
	id, ok := tenantID.(string)
	return id, ok
}

// GetScopes extracts the scopes from the context.
func GetScopes(c *gin.Context) ([]string, bool) {
	scopes, exists := c.Get("scopes")
	if !exists {
		return nil, false
	}
	s, ok := scopes.([]string)
	return s, ok
}
