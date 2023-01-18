package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer string
	Audience string
	Secret string
	TokenExpiry time.Duration
	RefreshExpiry time.Duration
	CookieDomain string
	CookiePath string
	CookieName string
}

type jwtUser struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

type TokenPairs struct {
	Token string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims) //token.Claims.(jwt.MapClaims) : token.Claims을 jwt.MapClaims로 Casting
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID) // subject
	claims["aud"] = j.Audience //Audience
	claims["iss"] = j.Issuer // issuer
	claims["iat"] = time.Now().UTC().Unix() //issued at
	claims["typ"] = "JWT"// type

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()// expiry

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create a signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token: signedAccessToken,
		RefreshToken:  signedRefreshToken,
	}

	// Return TokenPairs
	return tokenPairs, nil
}

// GetRefreshCookie returns a cookie containing the refresh token. Note that the cookie is http only, secure,
func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path : j.CookiePath,
		Value: refreshToken,
		Expires: time.Now().Add(j.RefreshExpiry),
		MaxAge: int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}

// GetExpiredRefreshCookie is a convenience method to return a cookie suitable for forcing a user's browser to delete existing cookie.
func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	//Delete the cookie from user browser
	return &http.Cookie{
		Name: j.CookieName,
		Path : j.CookiePath,
		Value: "",
		Expires: time.Unix(0,0),
		MaxAge: -1,
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error){
	w.Header().Add("Vary", "Authorization")

	// get auth header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2{
		// Authorization header has to be formed as 'Bearer afdsafasfddsx~~'
		return "", nil, errors.New("Invalid auth header")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		// Authorization header has to be formed as 'Bearer afdsafasfddsx~~'
		return "", nil, errors.New("Invalid auth header")
	}

	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// if signing method is not expected one
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by"){
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}