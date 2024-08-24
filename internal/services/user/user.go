package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func Register(login string, password string) error {
	isLoginExists, err := store.HasUserByLogin(login)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed check register user login", args)
		return err
	}

	if isLoginExists {
		args := map[string]interface{}{"login": login}
		logger.Info("failed register user: login already exists", args)
		return ErrorUserExists
	}

	hashPass, err := encodePassword(password)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed encode user password", args)
		return err
	}

	err = store.CreateUser(login, hashPass)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed register user", args)
		return err
	}

	return nil
}

func Login(login string, password string) (string, error) {
	u, err := store.GetUserByLogin(login)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed get user by login to auth", args)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed user login: compare password is fail", args)
		return "", err
	}

	tokenExpired := time.Now().Add(tokenLifetime)
	token, err := generateToken(u.ID, tokenExpired)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed generate token", args)
		return "", err
	}

	return token, nil
}

func CheckAuthorize(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	appConfig, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed get app config", args)
		return false
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(appConfig.GetJWTSecret()), nil
		})
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed parse token", args)
		return false
	}

	if !token.Valid || claims.UserID == 0 {
		return false
	}

	return true
}

func GetAuthUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, err
	}

	appConfig, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed get app config", args)
		return 0, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(appConfig.GetJWTSecret()), nil
		})
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed parse token", args)
		return 0, err
	}

	if !token.Valid || claims.UserID == 0 {
		return 0, err
	}

	return claims.UserID, nil
}

func encodePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed encode user password", args)
		return "", err
	}

	return string(hash), nil
}

func generateToken(userID int, expired time.Time) (string, error) {
	cnf, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed get app config", args)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expired),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(cnf.GetJWTSecret()))
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed generate token", args)
		return "", err
	}

	return tokenString, nil
}
