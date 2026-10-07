package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// 轻量 HS256 JWT（避免引入额外依赖），authMiddleware 解析同一套签名。

type jwtClaims struct {
	Sub      string `json:"sub"`
	Role     string `json:"role"`
	TeamCode string `json:"team"`
	Name     string `json:"name"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

// AuthIdentity 从令牌中解析出的身份。
type AuthIdentity struct {
	Username string
	Role     string
	TeamCode string
	Name     string
}

func b64(input []byte) string { return base64.RawURLEncoding.EncodeToString(input) }

func sign(input string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return b64(mac.Sum(nil))
}

// IssueJWT 签发 HS256 令牌。
func IssueJWT(id AuthIdentity, secret string, ttl time.Duration) string {
	now := time.Now()
	claims := jwtClaims{
		Sub: id.Username, Role: id.Role, TeamCode: id.TeamCode, Name: id.Name,
		Iat: now.Unix(), Exp: now.Add(ttl).Unix(),
	}
	body, _ := json.Marshal(claims)
	header := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := b64(body)
	signingInput := header + "." + payload
	return signingInput + "." + sign(signingInput, secret)
}

// ParseJWT 校验并解析令牌。
func ParseJWT(token, secret string) (*AuthIdentity, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("malformed token")
	}
	if sign(parts[0]+"."+parts[1], secret) != parts[2] {
		return nil, errors.New("bad signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("bad payload")
	}
	var claims jwtClaims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return nil, errors.New("bad claims")
	}
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}
	return &AuthIdentity{Username: claims.Sub, Role: claims.Role, TeamCode: claims.TeamCode, Name: claims.Name}, nil
}

// HashPassword 与种子用户口令一致使用 SHA256 摘要（本地演示，不引 bcrypt）。
func HashPassword(password, salt string) string {
	mac := hmac.New(sha256.New, []byte(salt))
	mac.Write([]byte(password))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
