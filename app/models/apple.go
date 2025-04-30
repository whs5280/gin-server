package models

import (
	"crypto/rsa"
	"time"
)

type AppleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type AppleCode struct {
	Code string `json:"code"`
}

type AppleToken struct {
	IdentityToken string `json:"identity_token" binding:"required"`
	Email         string `json:"email,omitempty"`
	FullName      string `json:"full_name,omitempty"`
}

// AppleConfig Apple登录配置
type AppleConfig struct {
	ClientID    string `json:"client_id"` // 例如：com.example.app
	TeamID      string `json:"team_id"`   // Apple开发者账号的Team ID
	KeyID       string `json:"key_id"`    // 私钥的Key ID
	RedirectURI string `json:"redirect_uri"`
}

// AppleUser Apple返回的用户信息
type AppleUser struct {
	ID            string `json:"id"`             // 用户唯一标识(sub)
	Email         string `json:"email"`          // 用户邮箱(可能为空)
	EmailVerified string `json:"email_verified"` // 邮箱是否验证
	FullName      string `json:"full_name"`      // 用户全名
}

// AppleJWK Apple的公钥结构
type AppleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// AppleJWKSet Apple公钥集合
type AppleJWKSet struct {
	Keys []AppleJWK `json:"keys"`
}

// ApplePublicKeyCache Apple公钥缓存
type ApplePublicKeyCache struct {
	Keys   map[string]*rsa.PublicKey
	Expiry time.Time
}
