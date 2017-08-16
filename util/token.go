package util

import (
	"fmt"
	"qnmahjong/def"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
)

// Claims custom token
type Claims struct {
	PlayerID  int32 `json:"player_id"`
	Channel   int32 `json:"channel"`
	Version   int32 `json:"version"`
	LoginType int32 `json:"login_type"`
	jwt.StandardClaims
}

// CreateToken create token
func CreateToken(claims *Claims) (signedToken string, success bool) {
	claims.ExpiresAt = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCreateToken)
		return
	}

	success = true
	return
}

// ValidateToken validate token
func ValidateToken(signedToken string) (claims *Claims, success bool) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrParseToken)
		return
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		success = true
		return
	}

	log.WithFields(log.Fields{
		"token": err,
	}).Error(def.ErrValidateToken)
	return
}
