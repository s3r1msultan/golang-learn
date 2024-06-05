package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"server/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	u := &User{
		ObjectId: primitive.NewObjectID(),
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
	}

	user, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &CreateUserRes{
		ID:       user.ObjectId.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

type MyJWTClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	context, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.Repository.GetUserByEmail(context, req.Email)
	if err != nil {
		return nil, err
	}
	err = util.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return &LoginUserRes{
		accessToken: tokenString,
		ID:          user.ObjectId.Hex(),
		Username:    user.Username,
		Email:       user.Email,
	}, nil
}
