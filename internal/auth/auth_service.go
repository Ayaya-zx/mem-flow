package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

var sercretKey = []byte("super-secret-key")

type AuthData struct {
	Name     string
	Password string
}

type AuthService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (as *AuthService) RegUser(authData *AuthData) (string, error) {
	if authData.Name == "" || authData.Password == "" {
		return "", common.InvalidAuthData(
			"name and password cannot be empty strings")
	}

	u := &entity.User{
		Name:       authData.Name,
		PasswdHash: hashPassword(authData.Password),
	}

	err := as.userRepo.AddUser(u)
	if err != nil {
		return "", err
	}
	return createToken(u)
}

func (as *AuthService) AuthUser(authData *AuthData) (string, error) {
	u, err := as.userRepo.GetUser(authData.Name)
	if err != nil {
		return "", err
	}
	if u.PasswdHash != hashPassword(authData.Password) {
		return "", common.InvalidAuthData(
			"incorect name or passowrd",
		)
	}
	return createToken(u)
}

func (as *AuthService) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sercretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", common.InvalidToken("invalid token")
	}

	raw, ok := claims["sub"]
	if !ok {
		return "", common.InvalidToken("invalid token")
	}

	name, ok := raw.(string)
	if !ok {
		return "", common.InvalidToken("invalid token")
	}

	return name, nil
}

func hashPassword(passwd string) string {
	h := sha256.New()
	h.Write([]byte(passwd))
	return string(h.Sum(nil))
}

func createToken(u *entity.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": u.Name,
			"iss": "mem-flow",
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Unix(),
		},
	)
	return token.SignedString(sercretKey)
}
