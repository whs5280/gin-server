package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gin-server/app/models"
	"github.com/dgrijalva/jwt-go"
	"io"
	"math/big"
	"net/http"
	"time"
)

type AppleService struct {
	Config     models.AppleConfig
	HttpClient *http.Client
	PublicKeys map[string]*rsa.PublicKey
	KeysExpiry time.Time
}

// VerifyAppleToken 验证Apple身份令牌
func (s *AppleService) VerifyAppleToken(identityToken string) (*models.AppleUser, error) {
	// 1. 获取Apple的公钥
	if err := s.fetchApplePublicKeys(); err != nil {
		//if err := s.fetchApplePublicKeysByRedis(); err != nil {
		return nil, fmt.Errorf("failed to fetch Apple public keys: %v", err)
	}

	// 2. 解析JWT令牌(不验证签名)
	token, _, err := new(jwt.Parser).ParseUnverified(identityToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// 3. 获取kid并查找对应的公钥
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid not found in token header")
	}

	publicKey, exists := s.PublicKeys[kid]
	if !exists {
		return nil, fmt.Errorf("public key for kid %s not found", kid)
	}

	// 4. 完整验证JWT令牌
	claims := jwt.MapClaims{}
	token, err = jwt.ParseWithClaims(identityToken, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	fmt.Println(claims)
	// 5. 验证claims
	if err := s.validateClaims(claims); err != nil {
		return nil, err
	}

	// 6. 提取用户信息
	user := &models.AppleUser{
		ID:            claims["sub"].(string),
		Email:         s.getStringClaim(claims, "email"),
		EmailVerified: s.getStringClaim(claims, "email_verified"),
	}

	return user, nil
}

// fetchApplePublicKeys 获取Apple的公钥
func (s *AppleService) fetchApplePublicKeys() error {
	// 如果公钥未过期，直接使用缓存
	if time.Now().Before(s.KeysExpiry) && len(s.PublicKeys) > 0 {
		return nil
	}

	resp, err := s.HttpClient.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jwkSet models.AppleJWKSet
	if err := json.Unmarshal(body, &jwkSet); err != nil {
		return err
	}

	// 解析并存储公钥
	for _, key := range jwkSet.Keys {
		if key.Use == "sig" && key.Alg == "RS256" {
			publicKey, err := jwkToRSAPublicKey(key)
			if err != nil {
				fmt.Printf("Failed to convert JWK to RSA key for kid %s: %v\n", key.Kid, err)
				continue
			}

			s.PublicKeys[key.Kid] = publicKey
			fmt.Printf("Successfully loaded key for kid: %s\n", key.Kid)
		}
	}

	fmt.Printf("Total keys loaded: %d\n", len(s.PublicKeys))
	for kid := range s.PublicKeys {
		fmt.Printf(" - Kid: %s\n", kid)
	}

	if len(s.PublicKeys) == 0 {
		return errors.New("no valid public keys found")
	}

	// 设置公钥过期时间(Apple建议每24小时刷新一次)
	s.KeysExpiry = time.Now().Add(23 * time.Hour)
	return nil
}

// validateClaims 验证JWT claims
func (s *AppleService) validateClaims(claims jwt.MapClaims) error {
	// 调试输出
	fmt.Printf("Validating claims - aud: %v, config clientID: %v\n", claims["aud"], s.Config.ClientID)

	// 验证issuer
	iss, ok := claims["iss"].(string)
	if !ok || iss != "https://appleid.apple.com" {
		return errors.New("invalid issuer")
	}

	// 验证audience
	aud, ok := claims["aud"].(string)
	if !ok {
		return errors.New("aud claim missing")
	}

	if aud != s.Config.ClientID {
		return fmt.Errorf("invalid audience, got %s, want %s", aud, s.Config.ClientID)
	}

	// 验证expiry
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("exp claim missing")
	}

	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf("token expired at %v", time.Unix(int64(exp), 0))
	}

	return nil
}

func jwkToRSAPublicKey(key models.AppleJWK) (*rsa.PublicKey, error) {
	// 解码Base64URL编码的模数(N)和指数(E)
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %v", err)
	}

	// 将字节转换为大整数
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes).Int64()

	return &rsa.PublicKey{
		N: n,
		E: int(e),
	}, nil
}

// getStringClaim 安全获取字符串类型的claim
func (s *AppleService) getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
