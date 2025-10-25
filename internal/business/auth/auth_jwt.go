package bussinessAuth

import (
	"errors"
	"fmt"
	"os"
	"time"
	LogService "zymm/internal/service/log_service"

	"github.com/golang-jwt/jwt/v5"
)

var MY_SECRET_KEY string = ""
var MySigningMethod = jwt.SigningMethodHS256

type MyCustomClaims struct {
	UserId int
	RoleId int
	Expiry int64
}

func getSecretKey() string {
	if MY_SECRET_KEY == "" {
		MY_SECRET_KEY = os.Getenv("ZYMM_BACKEND_AUTH_SECRET_KEY")
	}

	if MY_SECRET_KEY == "" {
		LogService.LogMessage("TOKEN_GET_SECRET_KEY is empty! Check environment variable.")
	} else {
		LogService.LogMessage("TOKEN_GET_SECRET_KEY loaded successfully.")
	}

	return MY_SECRET_KEY
}

var authTokenUserIdKey = "userId"
var authTokenRoleIdKey = "roleId"
var authTokenExpKey = "exp"

func GenerateJWTToken(userId int, roleId int) (*string, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return nil, fmt.Errorf("secret key is not set")
	}
	expirationDays := 30
	expirationTime := time.Now().Add(time.Hour * 24 * time.Duration(expirationDays)).Unix()
	claims := jwt.MapClaims{
		authTokenUserIdKey: userId,
		authTokenRoleIdKey: roleId,
		authTokenExpKey:    expirationTime,
	}

	token := jwt.NewWithClaims(MySigningMethod, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return nil, fmt.Errorf("missing JWT secret key")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return token, nil
}

func tokenClaimsValidation(claims jwt.MapClaims) error {
	if claims[authTokenUserIdKey] == nil ||
		claims[authTokenRoleIdKey] == nil ||
		claims[authTokenExpKey] == nil ||
		claims[authTokenUserIdKey] == "" ||
		claims[authTokenRoleIdKey] == "" ||
		claims[authTokenExpKey] == "" {
		return fmt.Errorf("invalid token claims")
	}
	return nil
}

func getPrintableClaims(claims *MyCustomClaims) string {
	return fmt.Sprintf("UserId: %d, RoleId: %d, Expiry: %s", claims.UserId, claims.RoleId, time.Unix(claims.Expiry, 0))
}

func GetDataFromJWTToken(tokenString string) (*MyCustomClaims, error) {
	claims, err := parseJWTToken(tokenString)
	if err != nil {
		return nil, err
	}
	claimValidationErr := tokenClaimsValidation(claims)
	if claimValidationErr != nil {
		return nil, claimValidationErr
	}
	return &MyCustomClaims{
		UserId: int(claims[authTokenUserIdKey].(float64)),
		RoleId: int(claims[authTokenRoleIdKey].(float64)),
		Expiry: int64(claims[authTokenExpKey].(float64)),
	}, nil
}

func parseJWTToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return nil, fmt.Errorf("secret key is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return claims, nil
}
