package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestResolveConfigPathFromServerDirectory 负责验证从 backend/cmd/server 启动时也能相对定位到配置文件。
func TestResolveConfigPathFromServerDirectory(t *testing.T) {
	// tempDir 提供测试期间可安全写入的临时目录。
	tempDir := t.TempDir()

	// projectRoot 模拟最小化的 backend 项目根目录。
	projectRoot := filepath.Join(tempDir, "backend")
	// configDir 模拟服务默认配置目录。
	configDir := filepath.Join(projectRoot, "configs")
	// serverDir 模拟 dlv 从 backend/cmd/server 启动时的工作目录。
	serverDir := filepath.Join(projectRoot, "cmd", "server")

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	if err := os.MkdirAll(serverDir, 0o755); err != nil {
		t.Fatalf("创建服务目录失败: %v", err)
	}

	// configPath 指向本次测试使用的默认配置文件。
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("app:\n  name: test\n"), 0o644); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	// previousDir 保存测试前的工作目录，便于结束后恢复现场。
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前工作目录失败: %v", err)
	}

	defer func() {
		if chdirErr := os.Chdir(previousDir); chdirErr != nil {
			t.Fatalf("恢复工作目录失败: %v", chdirErr)
		}
	}()

	if err := os.Chdir(serverDir); err != nil {
		t.Fatalf("切换到服务目录失败: %v", err)
	}

	// resolvedPath 保存默认配置定位逻辑实际找到的配置文件路径。
	resolvedPath, err := resolveConfigPath("")
	if err != nil {
		t.Fatalf("期望可以从 server 目录定位到配置文件，但返回错误: %v", err)
	}

	// expectedPath 保存对配置文件路径做符号链接归一化后的结果。
	expectedPath, err := filepath.EvalSymlinks(configPath)
	if err != nil {
		t.Fatalf("解析期望配置路径失败: %v", err)
	}

	// actualPath 保存对实际定位结果做符号链接归一化后的结果。
	actualPath, err := filepath.EvalSymlinks(resolvedPath)
	if err != nil {
		t.Fatalf("解析实际配置路径失败: %v", err)
	}

	if actualPath != expectedPath {
		t.Fatalf("期望定位到 %s，实际为 %s", expectedPath, actualPath)
	}
}

// TestAuthConfigValidate 负责验证鉴权配置的关键安全约束会被正确校验。
func TestAuthConfigValidate(t *testing.T) {
	// testCases 保存当前测试要覆盖的鉴权配置校验场景。
	testCases := []struct {
		// Name 表示当前子测试的名称。
		Name string
		// Config 表示待校验的鉴权配置。
		Config AuthConfig
		// WantErr 表示当前场景是否期望返回校验错误。
		WantErr bool
	}{
		{
			Name: "合法配置可以通过校验",
			Config: AuthConfig{
				Issuer:          "operation-admin-backend",
				JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
				AccessTokenTTL:  2 * time.Hour,
				RefreshTokenTTL: 168 * time.Hour,
			},
			WantErr: false,
		},
		{
			Name: "密钥过短会被拒绝",
			Config: AuthConfig{
				Issuer:          "operation-admin-backend",
				JWTSecret:       "short-secret",
				AccessTokenTTL:  2 * time.Hour,
				RefreshTokenTTL: 168 * time.Hour,
			},
			WantErr: true,
		},
		{
			Name: "AccessTokenTTL 必须大于零",
			Config: AuthConfig{
				Issuer:          "operation-admin-backend",
				JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
				AccessTokenTTL:  0,
				RefreshTokenTTL: 168 * time.Hour,
			},
			WantErr: true,
		},
		{
			Name: "RefreshTokenTTL 必须大于零",
			Config: AuthConfig{
				Issuer:          "operation-admin-backend",
				JWTSecret:       "abcdefghijklmnopqrstuvwxyz123456",
				AccessTokenTTL:  2 * time.Hour,
				RefreshTokenTTL: 0,
			},
			WantErr: true,
		},
	}

	for _, testCase := range testCases {
		// testCase 保存当前循环内执行的测试用例副本，避免闭包引用同一地址。
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			// err 保存本次鉴权配置校验返回的结果。
			err := testCase.Config.Validate()

			if testCase.WantErr && err == nil {
				t.Fatalf("期望返回校验错误，但实际为 nil")
			}

			if !testCase.WantErr && err != nil {
				t.Fatalf("期望校验通过，但返回错误: %v", err)
			}
		})
	}
}
