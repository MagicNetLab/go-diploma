package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
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
		return err
	}

	if isLoginExists {
		return ErrorUserExists
	}

	hashPass, err := encodePassword(password)
	if err != nil {
		return err
	}

	err = store.CreateUser(login, hashPass)
	if err != nil {
		return err
	}

	return nil
}

func Login(login string, password string) (string, error) {
	u, err := store.GetUserByLogin(login)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := generateToken(u.ID)
	if err != nil {
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
		return "", err
	}

	return string(hash), nil
}

func generateToken(userID int) (string, error) {
	cnf, err := config.GetAppConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(cnf.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
