package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	config *Config
}

type jwtCustomClaim struct {
	Id        string  `json:"id"`
	ExpiresAt int64  `json:"expires_at"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

const (
	StandardRole = "standard"
)

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{config: config}
}

func (j *JWTToken) CreateToken(userID string, isAdmin bool, ttl time.Duration) (string, error) {

	role := StandardRole
	claims := jwtCustomClaim{
		Id:        userID,
		ExpiresAt: time.Now().Add(ttl * time.Minute).Unix(),
		Role:      role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.SigningKey))

	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid authentication token")
		}
		return []byte(j.config.SigningKey), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid authentication token")
	}

	claims, ok := token.Claims.(*jwtCustomClaim)

	if !ok {
		return "", "", fmt.Errorf("invalid authentication token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", "", fmt.Errorf("token has expired")
	}

	return claims.Id, claims.Role, nil
}

// ################################################################

// function to fetch googgle public key

// func fetchGooglePublicKeys() (map[string]*rsa.PublicKey, error) {
// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching Google public keys: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %v", err)
// 	}

// 	var keys map[string]*rsa.PublicKey
// 	if err := json.Unmarshal(body, &keys); err != nil {
// 		return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
// 	}

// 	return keys, nil
// }

// func verifyGoogleIDToken(idToken string) (*jwtCustomClaim, error) {
// 	keys, err := fetchGooglePublicKeys()
// 	if err != nil {
// 		return nil, err
// 	}

// 	token, err := jwt.ParseWithClaims(idToken, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
// 		// Validate the token signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// Get the key ID from the token header
// 		kid, ok := token.Header["kid"].(string)
// 		if !ok {
// 			return nil, fmt.Errorf("missing or invalid key ID in token header")
// 		}

// 		// Retrieve the public key for the given key ID
// 		key, found := keys[kid]
// 		if !found {
// 			return nil, fmt.Errorf("public key not found for key ID: %v", kid)
// 		}

// 		return key, nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing Google ID token: %v", err)
// 	}

// 	// Extract custom claims
// 	claims, ok := token.Claims.(*jwtCustomClaim)
// 	if !ok || !token.Valid {
// 		return nil, fmt.Errorf("invalid Google ID token")
// 	}

// 	return claims, nil
// }
