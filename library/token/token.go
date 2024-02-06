package token

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
	"wusthelper-manager-go/library/log"
)

type Token struct {
	SecretKey string
	Timeout   time.Duration
}

type SignPayload struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
}

// New 一个token对象，secret为token加密密钥，timeout为有效期，单位为天
func New(secret string, timeout time.Duration) *Token {
	return &Token{
		SecretKey: secret,
		Timeout:   time.Hour * 24 * timeout,
	}
}

func (t *Token) Sign(payload SignPayload) (token string) {
	claims := &jwt.MapClaims{
		"id":       payload.Uid,
		"username": payload.Username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Unix() + int64(t.Timeout.Seconds()),
	}
	token, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, *claims).
		SignedString([]byte(t.SecretKey))
	return
}

func (t *Token) Verify(token string) bool {
	tt, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})

	return err == nil && tt.Valid
}

func (t *Token) GetClaimVerify(token string) (*jwt.MapClaims, bool) {
	claims := new(jwt.MapClaims)
	_, err := jwt.NewParser().ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		return nil, false
	}

	return claims, true
}

func (t *Token) GetClaimWithoutVerify(token string) *jwt.MapClaims {
	claims := new(jwt.MapClaims)
	_, _, err := jwt.NewParser().ParseUnverified(token, claims)
	if err != nil {
		log.Warn("token解析（不验证）错误", zap.String("err", err.Error()))
		return nil
	}

	return claims
}
