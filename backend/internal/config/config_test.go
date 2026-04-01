package config

import (
	"os"
	"path/filepath"
	"testing"
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
