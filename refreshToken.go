package main

import (
	"context"
	"errors"
	"time"

	"github.com/ajdinahmetovic/item-service/db"
	"github.com/ajdinahmetovic/item-service/proto/v1"
	"github.com/dgrijalva/jwt-go"
)

//RefreshToken func
func (s *Server) RefreshToken(ctx context.Context, request *proto.RefreshTokenReq) (*proto.RefreshTokenRes, error) {

	var accesToken, refreshToken string
	token := request.GetRefreshToken()
	newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("tajna"), nil
	})
	if err != nil {
		return &proto.RefreshTokenRes{RefreshToken: "", Token: ""}, err
	}
	if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
		userID := int(claims["CustomUserInfo"].(map[string]interface{})["ID"].(float64))
		if err != nil {
			return &proto.RefreshTokenRes{RefreshToken: "", Token: token}, err
		}
		res, err := db.FindUser(&db.User{
			ID: userID,
		})
		if err != nil {
			return &proto.RefreshTokenRes{RefreshToken: "", Token: token}, err
		}
		if len(res) < 1 {
			return nil, errors.New("User does not exist")
		}
		accesToken, err = generateToken(&userID, time.Now().Add(time.Hour*1).Unix())
		if err != nil {
			return nil, errors.New("Failed to generate token")
		}
		refreshToken, err = generateToken(&userID, time.Now().Add(time.Hour*24).Unix())
		if err != nil {
			return nil, errors.New("Failed to generate token")
		}
	} else {
		return nil, errors.New("Refresh token is not valid")
	}
	return &proto.RefreshTokenRes{RefreshToken: refreshToken, Token: accesToken}, nil
}

func generateToken(id *int, exp int64) (string, error) {
	signingKey := []byte("tajna")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": &exp,
		"CustomUserInfo": struct {
			ID int
		}{int(*id)},
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
