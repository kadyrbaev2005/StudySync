package services

import (
    "errors"
    "time"

    "github.com/kadyrbayev2005/studysync/internal/utils"
    "golang.org/x/crypto/bcrypt"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func InitJWT() {
    secret := utils.GetEnv("JWT_SECRET", "your-secure-secret-key")
    jwtKey = []byte(secret)
}

func GenerateJWT(userID uint, role string) (string, error) {
    expiration := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiration),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}

func HashPassword(pw string) string {
    b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(b)
}

func CheckPasswordHash(pw, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
    return err == nil
}