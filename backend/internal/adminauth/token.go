package adminauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"operation_admin/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// accessTokenJTIByteLength 定义 Access Token JTI 随机字节长度。
	accessTokenJTIByteLength = 16
	// refreshTokenByteLength 定义原始 Refresh Token 随机字节长度。
	refreshTokenByteLength = 32
	// refreshTokenJTIByteLength 定义 Refresh Token JTI 随机字节长度。
	refreshTokenJTIByteLength = 16
)

// TokenManager 描述后台鉴权所需的 Token 签发与校验能力。
type TokenManager struct {
	// Issuer 保存当前服务用于写入和校验 JWT 的签发方标识。
	Issuer string
	// JWTSecret 保存当前服务用于 HS256 签名和验签的密钥。
	JWTSecret []byte
	// AccessTokenTTL 保存 Access Token 的统一有效期。
	AccessTokenTTL time.Duration
	// RefreshTokenTTL 保存 Refresh Token 的统一有效期。
	RefreshTokenTTL time.Duration
}

// AccessClaims 描述后台 Access Token 中承载的业务 Claims。
type AccessClaims struct {
	// LoginName 表示当前 Access Token 对应的后台账号登录名。
	LoginName string `json:"loginName"`
	// RoleCode 表示当前 Access Token 对应的后台账号角色编码。
	RoleCode string `json:"roleCode"`
	// TokenType 表示当前 JWT 的业务类型，固定为 access。
	TokenType string `json:"tokenType"`
	// RegisteredClaims 保存 sub、jti、iss、iat、exp 等标准 JWT 字段。
	jwt.RegisteredClaims
}

// GeneratedAccessToken 描述一次 Access Token 签发结果。
type GeneratedAccessToken struct {
	// Token 保存最终返回给客户端的 JWT 字符串。
	Token string
	// TokenJTI 保存当前 Access Token 对应的唯一标识。
	TokenJTI string
	// ExpiresAt 保存当前 Access Token 的过期时间。
	ExpiresAt time.Time
	// ExpiresInSeconds 保存当前 Access Token 的有效期秒数。
	ExpiresInSeconds int64
}

// GeneratedRefreshToken 描述一次 Refresh Token 签发结果。
type GeneratedRefreshToken struct {
	// Token 保存最终返回给客户端的原始 Refresh Token。
	Token string
	// TokenHash 保存落库使用的 SHA-256 哈希值。
	TokenHash string
	// TokenJTI 保存当前 Refresh Token 记录的唯一标识。
	TokenJTI string
	// ExpiresAt 保存当前 Refresh Token 的过期时间。
	ExpiresAt time.Time
	// ExpiresInSeconds 保存当前 Refresh Token 的有效期秒数。
	ExpiresInSeconds int64
}

// NewTokenManager 负责根据运行配置创建后台鉴权 Token 管理器。
func NewTokenManager(cfg config.AuthConfig) *TokenManager {
	return &TokenManager{
		Issuer:          cfg.Issuer,
		JWTSecret:       []byte(cfg.JWTSecret),
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

// GenerateAccessToken 负责为指定后台账号签发新的 Access Token。
func (m *TokenManager) GenerateAccessToken(adminUser AdminUser) (*GeneratedAccessToken, error) {
	// tokenJTI 保存当前 Access Token 对应的唯一标识。
	tokenJTI, err := generateRandomString(accessTokenJTIByteLength)
	if err != nil {
		return nil, err
	}

	// issuedAt 保存当前 Access Token 的签发时间。
	issuedAt := time.Now()
	// expiresAt 保存当前 Access Token 的过期时间。
	expiresAt := issuedAt.Add(m.AccessTokenTTL)

	// claims 保存当前后台账号需要写入 JWT 的业务 Claims。
	claims := AccessClaims{
		LoginName: adminUser.LoginName,
		RoleCode:  adminUser.RoleCode,
		TokenType: AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   adminUser.AdminUserID,
			ID:        tokenJTI,
			Issuer:    m.Issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// token 保存构造完成但尚未签名的 JWT 对象。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// tokenString 保存最终返回给调用方的 JWT 字符串。
	tokenString, err := token.SignedString(m.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("签发 access token 失败: %w", err)
	}

	return &GeneratedAccessToken{
		Token:            tokenString,
		TokenJTI:         tokenJTI,
		ExpiresAt:        expiresAt,
		ExpiresInSeconds: int64(m.AccessTokenTTL / time.Second),
	}, nil
}

// ParseAccessToken 负责解析并校验 Access Token 的签名、签发方和业务类型。
func (m *TokenManager) ParseAccessToken(rawToken string) (*AccessClaims, error) {
	// claims 保存 JWT 解析成功后得到的业务 Claims。
	claims := &AccessClaims{}

	// parsedToken 保存 JWT 解析器返回的原始 token 对象。
	parsedToken, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrAccessTokenInvalid
			}

			return m.JWTSecret, nil
		},
		jwt.WithIssuer(m.Issuer),
	)
	if err != nil {
		return nil, ErrAccessTokenInvalid
	}

	if !parsedToken.Valid {
		return nil, ErrAccessTokenInvalid
	}

	if claims.TokenType != AccessTokenType {
		return nil, ErrAccessTokenInvalid
	}

	if claims.Subject == "" || claims.LoginName == "" || claims.RoleCode == "" || claims.ID == "" {
		return nil, ErrAccessTokenInvalid
	}

	return claims, nil
}

// GenerateRefreshToken 负责为当前登录会话生成原始 Refresh Token 及其落库元数据。
func (m *TokenManager) GenerateRefreshToken() (*GeneratedRefreshToken, error) {
	// rawToken 保存最终返回给客户端的原始 Refresh Token。
	rawToken, err := generateRandomString(refreshTokenByteLength)
	if err != nil {
		return nil, err
	}

	// tokenJTI 保存当前 Refresh Token 记录的唯一标识。
	tokenJTI, err := generateRandomString(refreshTokenJTIByteLength)
	if err != nil {
		return nil, err
	}

	// expiresAt 保存当前 Refresh Token 的过期时间。
	expiresAt := time.Now().Add(m.RefreshTokenTTL)

	return &GeneratedRefreshToken{
		Token:            rawToken,
		TokenHash:        HashToken(rawToken),
		TokenJTI:         tokenJTI,
		ExpiresAt:        expiresAt,
		ExpiresInSeconds: int64(m.RefreshTokenTTL / time.Second),
	}, nil
}

// HashToken 负责把原始 Refresh Token 计算为落库使用的 SHA-256 十六进制哈希。
func HashToken(rawToken string) string {
	// hashBytes 保存原始 Refresh Token 对应的 SHA-256 摘要。
	hashBytes := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hashBytes[:])
}

// generateRandomString 负责生成 URL 安全的随机字符串。
func generateRandomString(byteLength int) (string, error) {
	// randomBytes 保存当前随机字符串底层使用的随机字节数组。
	randomBytes := make([]byte, byteLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("生成随机字符串失败: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
