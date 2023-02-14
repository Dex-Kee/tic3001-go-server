package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"tic3001-go-server/common/constant"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/entity"
	"tic3001-go-server/utils"
	"time"
)

type authService struct{}

var (
	AuthService   = new(authService)
	jwtSigningKey = []byte("jwtSigningKey")
	digestKey     = "digestSHA256"
)

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	Auth dto.AuthDto
}

func (s authService) Login(form dto.LoginForm) (string, error) {
	// check existence of user
	user := s.findUserByName(form.Username)
	if user == nil {
		return "", fmt.Errorf("error: username [%s] is not found", form.Username)
	}

	// check credential of user
	if utils.SHA256Digest(form.Password+digestKey) != user.HashedPassword {
		return "", fmt.Errorf("error: incorrect password")
	}

	// create jwt claim
	var claims = JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    constant.TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constant.TokenValidityDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Auth: s.createAuthDto(user),
	}

	token, err := s.createToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s authService) findUserByName(username string) *entity.User {
	// hardcode for now
	if username == "admin" {
		return &entity.User{
			Id:             uuid.NewString(),
			Name:           "admin",
			HashedPassword: utils.SHA256Digest("admin123456" + digestKey),
			Role:           "admin",
		}
	}
	if username == "test" {
		return &entity.User{
			Id:             uuid.NewString(),
			Name:           "test",
			HashedPassword: utils.SHA256Digest("test123456" + digestKey),
			Role:           "customer",
		}
	}
	return nil
}

func (s authService) createAuthDto(user *entity.User) dto.AuthDto {
	return dto.AuthDto{
		Id:            user.Id,
		UserName:      user.Name,
		Role:          user.Role,
		LastLoginTime: time.Now(),
	}
}

func (s authService) createToken(claims JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSigningKey)
}

func (s authService) ParserToken(tokenString string) (*JwtCustomClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is not found")
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*JwtCustomClaims); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, errors.New("token is invalid")
}

func (s authService) FindAccessibleResourceByRole(role string) map[string]bool {
	// query accessible resource by role
	// hard code for now
	resourceMap := make(map[string]bool)
	resourceMap["/api/notes/list"] = true
	resourceMap["/api/notes/create"] = true
	resourceMap["/api/notes/update"] = true
	resourceMap["/api/notes/delete"] = true
	resourceMap["/api/user/delete"] = false
	return resourceMap
}
