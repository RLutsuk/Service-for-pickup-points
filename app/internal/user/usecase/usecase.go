package usecase

import (
	"errors"
	"log/slog"
	"regexp"
	"time"

	userRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/repository"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "2b42e820074d4141beaf4b3018c5360a71ee0b0f05cc0153646dd73ec5f9a3c9"
	tokenTll   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   string `json:"id"`
	UserRole string `json:"role"`
}

type UseCaseI interface {
	CreateUser(user *models.User) error
	AuthUser(user *models.User) (string, error)
	TestUser(role string) (string, error)
}

type useCase struct {
	userRepository userRep.RepositoryI
	logger         *slog.Logger
}

func New(userRepository userRep.RepositoryI, logger *slog.Logger) UseCaseI {
	return &useCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (uc *useCase) CreateUser(user *models.User) error {

	if err := uc.userRepository.ChekUserByEmail(user.Email); err != nil {
		uc.logger.Error("error with the email user search", slog.String("error", err.Error()))
		return err
	}

	if ok := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(user.Email); !ok {
		uc.logger.Error("invalid email")
		return models.ErrBadEmail
	}

	if user.Role != "employee" && user.Role != "moderator" {
		uc.logger.Error("invalid role")
		return models.ErrBadData
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("error with hash password", slog.String("error", err.Error()))
		return err
	}
	user.Password = string(hash)

	user, err = uc.userRepository.CreateUser(user)
	if err != nil {
		uc.logger.Error("error with creation user", slog.String("error", err.Error()))
		return err
	}

	user.Password = ""
	return err
}

func (uc *useCase) AuthUser(user *models.User) (string, error) {

	password := user.Password
	err := uc.userRepository.GetUserByEmail(user)
	if err != nil {
		uc.logger.Error("error with the email user search", slog.String("error", err.Error()))
		return "", models.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		uc.logger.Error("invalid password", slog.String("error", err.Error()))
		return "", models.ErrUserNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Role,
	})
	return token.SignedString([]byte(signingKey))
}

func (uc *useCase) TestUser(role string) (string, error) {

	if role != "employee" && role != "moderator" {
		uc.logger.Error("invalid role")
		return "", models.ErrBadData
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		"testuserID",
		role,
	})
	return token.SignedString([]byte(signingKey))
}

func Parsetoken(accessToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return "", "", models.ErrBadAuthorizated
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", "", models.ErrBadAuthorizated
	}

	return claims.UserId, claims.UserRole, nil
}
