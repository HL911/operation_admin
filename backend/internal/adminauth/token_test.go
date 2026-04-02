package adminauth

import (
	"testing"
	"time"

	"operation_admin/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TestTokenManagerGenerateAndParseAccessToken 负责验证 Access Token 可以正常签发和解析。
func TestTokenManagerGenerateAndParseAccessToken(t *testing.T) {
	// tokenManager 提供当前测试场景使用的 Token 签发和校验能力。
	tokenManager := NewTokenManager(config.AuthConfig{
		Issuer:          "operation-admin-backend",
		JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 168 * time.Hour,
	})

	// generatedToken 保存本次测试场景签发得到的 Access Token 结果。
	generatedToken, err := tokenManager.GenerateAccessToken(AdminUser{
		AdminUserID: "ADM-1001",
		LoginName:   "admin.root",
		RoleCode:    RoleCodeAdmin,
	})
	if err != nil {
		t.Fatalf("期望 access token 签发成功，但返回错误: %v", err)
	}

	// claims 保存 Access Token 解析成功后得到的业务 Claims。
	claims, err := tokenManager.ParseAccessToken(generatedToken.Token)
	if err != nil {
		t.Fatalf("期望 access token 解析成功，但返回错误: %v", err)
	}

	if claims.Subject != "ADM-1001" {
		t.Fatalf("期望 sub 为 ADM-1001，实际为 %s", claims.Subject)
	}

	if claims.LoginName != "admin.root" {
		t.Fatalf("期望 loginName 为 admin.root，实际为 %s", claims.LoginName)
	}

	if claims.RoleCode != RoleCodeAdmin {
		t.Fatalf("期望 roleCode 为 admin，实际为 %s", claims.RoleCode)
	}

	if claims.TokenType != AccessTokenType {
		t.Fatalf("期望 tokenType 为 access，实际为 %s", claims.TokenType)
	}
}

// TestTokenManagerParseAccessTokenFailure 负责验证常见异常 Access Token 会被拒绝。
func TestTokenManagerParseAccessTokenFailure(t *testing.T) {
	// validManager 提供签发合法 Access Token 所需的 Token 管理器。
	validManager := NewTokenManager(config.AuthConfig{
		Issuer:          "operation-admin-backend",
		JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 168 * time.Hour,
	})

	// wrongIssuerManager 提供错误 issuer 校验场景使用的 Token 管理器。
	wrongIssuerManager := NewTokenManager(config.AuthConfig{
		Issuer:          "another-backend",
		JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 168 * time.Hour,
	})

	// expiredManager 提供过期 Access Token 校验场景使用的 Token 管理器。
	expiredManager := NewTokenManager(config.AuthConfig{
		Issuer:          "operation-admin-backend",
		JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
		AccessTokenTTL:  -1 * time.Minute,
		RefreshTokenTTL: 168 * time.Hour,
	})

	// validToken 保存用于错误 issuer 场景的合法 Access Token。
	validToken, err := validManager.GenerateAccessToken(AdminUser{
		AdminUserID: "ADM-1001",
		LoginName:   "admin.root",
		RoleCode:    RoleCodeAdmin,
	})
	if err != nil {
		t.Fatalf("生成合法 access token 失败: %v", err)
	}

	// expiredToken 保存一个已经过期的 Access Token。
	expiredToken, err := expiredManager.GenerateAccessToken(AdminUser{
		AdminUserID: "ADM-1002",
		LoginName:   "admin.expired",
		RoleCode:    RoleCodeAdmin,
	})
	if err != nil {
		t.Fatalf("生成过期 access token 失败: %v", err)
	}

	// wrongTypeClaims 保存 tokenType 不合法的 Access Token Claims。
	wrongTypeClaims := AccessClaims{
		LoginName: "admin.root",
		RoleCode:  RoleCodeAdmin,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "ADM-1003",
			ID:        "jti-test",
			Issuer:    "operation-admin-backend",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	// wrongTypeToken 保存使用错误 tokenType 签发得到的 JWT 字符串。
	wrongTypeToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongTypeClaims).SignedString([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	if err != nil {
		t.Fatalf("生成错误 tokenType 的 JWT 失败: %v", err)
	}

	// testCases 保存本次测试需要覆盖的异常 Token 解析场景。
	testCases := []struct {
		// Name 表示当前子测试名称。
		Name string
		// Manager 表示当前场景用于解析 Access Token 的 Token 管理器。
		Manager *TokenManager
		// Token 表示当前场景待解析的原始 Access Token 字符串。
		Token string
	}{
		{
			Name:    "错误 issuer 会被拒绝",
			Manager: wrongIssuerManager,
			Token:   validToken.Token,
		},
		{
			Name:    "过期 token 会被拒绝",
			Manager: validManager,
			Token:   expiredToken.Token,
		},
		{
			Name:    "错误 tokenType 会被拒绝",
			Manager: validManager,
			Token:   wrongTypeToken,
		},
	}

	for _, testCase := range testCases {
		// testCase 保存当前循环执行的测试用例副本，避免闭包复用同一地址。
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			// err 保存当前异常 Token 解析返回的错误结果。
			_, err := testCase.Manager.ParseAccessToken(testCase.Token)
			if err == nil {
				t.Fatalf("期望 access token 解析失败，但实际成功")
			}
		})
	}
}
