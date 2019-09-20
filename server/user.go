package main

import (
	"context"
	"time"

	"github.com/ajdinahmetovic/item-service/db"
	"github.com/ajdinahmetovic/item-service/proto/v1"
	"github.com/dgrijalva/jwt-go"
)

//Login func
func (s *Server) Login(ctx context.Context, request *proto.LoginUserReq) (*proto.LoginUserRes, error) {
	userCredidentials := request.GetUserCredidentials()
	id, err := db.Login(&db.UserCredentials{
		Username: userCredidentials.Username,
		Password: userCredidentials.Password,
	})
	if err != nil {
		return &proto.LoginUserRes{Message: "Login failed wrong username or password"}, err
	}
	signingKey := []byte("tajna")

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"CustomUserInfo": struct {
			ID int
		}{*id},
	})
	refreshTokenStr, err := refreshToken.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"CustomUserInfo": struct {
			ID int
		}{*id},
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return &proto.LoginUserRes{Message: "Log in failed"}, err
	}
	return &proto.LoginUserRes{Message: "Log in successfull", Token: tokenString, RefreshToken: refreshTokenStr, UserID: int32(*id)}, nil
}

//CreateUser func
func (s *Server) CreateUser(ctx context.Context, request *proto.CreateUserReq) (*proto.CreateUserRes, error) {
	user := request.GetUser()
	id, err := db.AddUser(&db.User{
		ID:       int(user.ID),
		Username: user.Username,
		FullName: user.FullName,
		Password: user.Password,
	})
	if err != nil {
		return &proto.CreateUserRes{Message: "Failed to cerate user"}, err
	}

	signingKey := []byte("tajna")
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"CustomUserInfo": struct {
			ID int
		}{*id},
	})
	refreshTokenStr, err := refreshToken.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"CustomUserInfo": struct {
			ID int
		}{int(*id)},
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return &proto.CreateUserRes{Message: "Failed to sign token"}, err
	}
	return &proto.CreateUserRes{Message: "Successfully created user", Token: tokenString, RefreshToken: refreshTokenStr, UserID: int32(*id)}, nil
}

//GetUser func
func (s *Server) GetUser(ctx context.Context, request *proto.GetUserReq) (*proto.GetUserRes, error) {
	user := request.GetUser()
	users, err := db.FindUser(&db.User{
		ID:       int(user.ID),
		Username: user.Username,
		FullName: user.FullName,
	})
	if err != nil {
		return &proto.GetUserRes{Users: nil, Message: "Failed to query database"}, err
	}
	res := make([]*proto.User, 0)
	for _, i := range users {
		items := make([]*proto.Item, 0)
		for _, j := range i.Items {
			items = append(items, &proto.Item{
				ID:          int32(j.ID),
				Title:       j.Title,
				Description: j.Description,
				UserID:      int32(j.UserID),
			})
		}
		res = append(res, &proto.User{
			ID:       int32(i.ID),
			Username: i.Username,
			FullName: i.FullName,
			Items:    items,
		})
	}
	return &proto.GetUserRes{Users: res, Message: "Users found"}, err
}

//UpdateUser func
func (s *Server) UpdateUser(ctx context.Context, request *proto.UpdateUserReq) (*proto.UpdateUserRes, error) {
	user := request.GetUser()
	err := db.UpdateUser(&db.User{
		ID:       int(user.ID),
		Username: user.Username,
		FullName: user.FullName,
	})
	if err != nil {
		return &proto.UpdateUserRes{Message: "Failed to update user"}, err
	}
	return &proto.UpdateUserRes{Message: "User updated"}, err
}

//DeleteUser func
func (s *Server) DeleteUser(ctx context.Context, request *proto.DeleteUserReq) (*proto.DeleteUserRes, error) {
	id := request.GetID()
	err := db.DeleteUser(int(id))
	if err != nil {
		return &proto.DeleteUserRes{Message: "Failed to delete user"}, err
	}
	return &proto.DeleteUserRes{Message: "User deleted"}, err
}
