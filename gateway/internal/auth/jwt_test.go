package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// testPrivateKey generates a test RSA private key (2048 bits)
func testPrivateKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}

// testJWKS returns a test JWKS with the public key
func testJWKS(pubKey *rsa.PublicKey) string {
	n := pubKey.N.Bytes()
	e := big.NewInt(int64(pubKey.E)).Bytes()

	nEnc := base64.RawURLEncoding.EncodeToString(n)
	eEnc := base64.RawURLEncoding.EncodeToString(e)

	return `{"keys":[{"kty":"RSA","kid":"test-key-1","use":"sig","n":"` + nEnc + `","e":"` + eEnc + `","alg":"RS256"}]}`
}

func TestValidateToken_ValidJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		Scopes:    []string{"read", "write"},
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "https://test-issuer.com",
		Audience:  []string{"test-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	validatedClaims, err := client.ValidateToken(tokenString)
	require.NoError(t, err)
	require.Equal(t, "test-user", validatedClaims.Subject)
	require.Equal(t, "test-tenant", validatedClaims.TenantID)
	require.Equal(t, []string{"read", "write"}, validatedClaims.Scopes)
}

func TestValidateToken_ExpiredJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		IssuedAt:  time.Now().Unix() - 7200,
		ExpiresAt: time.Now().Unix() - 3600,
		Issuer:    "https://test-issuer.com",
		Audience:  []string{"test-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	_, err = client.ValidateToken(tokenString)
	require.Error(t, err)
	require.Contains(t, err.Error(), "expired")
}

func TestValidateToken_InvalidIssuer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://expected-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "https://wrong-issuer.com",
		Audience:  []string{"test-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	_, err = client.ValidateToken(tokenString)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid issuer")
}

func TestValidateToken_InvalidAudience(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "expected-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "https://test-issuer.com",
		Audience:  []string{"wrong-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	_, err = client.ValidateToken(tokenString)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid audience")
}

func TestValidateToken_MalformedJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	_, err = client.ValidateToken("invalid.token.string")
	require.Error(t, err)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()
	differentKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "https://test-issuer.com",
		Audience:  []string{"test-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(differentKey)
	require.NoError(t, err)

	_, err = client.ValidateToken(tokenString)
	require.Error(t, err)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	claims := Claims{
		Subject:   "test-user",
		TenantID:  "test-tenant",
		Scopes:    []string{"read", "write"},
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600,
		Issuer:    "https://test-issuer.com",
		Audience:  []string{"test-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-1"
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	called := false
	handler := func(c *gin.Context) {
		called = true
		userID, _ := GetUserID(c)
		require.Equal(t, "test-user", userID)
		tenantID, _ := GetTenantID(c)
		require.Equal(t, "test-tenant", tenantID)
		scopes, _ := GetScopes(c)
		require.Equal(t, []string{"read", "write"}, scopes)
	}

	middleware := client.AuthMiddleware()
	middleware(c)
	handler(c)

	require.True(t, called)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := client.AuthMiddleware()
	middleware(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "AUTH_FAILED")
	require.Contains(t, w.Body.String(), "missing authorization header")
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey := testPrivateKey()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJWKS(&privateKey.PublicKey)))
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		JWKSURL:   ts.URL,
		Issuer:    "https://test-issuer.com",
		Audience:  "test-audience",
		Algorithm: "RS256",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat token")

	middleware := client.AuthMiddleware()
	middleware(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.Contains(t, w.Body.String(), "AUTH_FAILED")
	require.Contains(t, w.Body.String(), "invalid authorization header format")
}