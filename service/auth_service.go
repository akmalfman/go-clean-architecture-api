package service

import (
	"errors"
	"time"

	"github.com/akmalfsalman/go-clean-architecture-api/models"
	"github.com/akmalfsalman/go-clean-architecture-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("kunci_rahasia_super_aman_banget_lho")

type JwtClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(user models.User) (models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(user models.User) (models.User, error) {
	if user.Email == "" || user.Password == "" {
		return models.User{}, errors.New("email dan password tidak boleh kosong")
	}
	if len(user.Password) < 6 {
		return models.User{}, errors.New("password minimal 6 karakter")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	newID, err := s.userRepo.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}

	user.ID = newID
	user.Password = ""
	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	foundUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	token, err := s.generateJWT(foundUser.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) generateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
