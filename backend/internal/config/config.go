package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// Config 描述整个应用的运行配置。
type Config struct {
	// App 描述应用名称等基础元信息，用于日志、监控和对外标识。
	App AppConfig `mapstructure:"app"`
	// Server 描述 HTTP 服务监听地址、运行模式与超时参数。
	Server ServerConfig `mapstructure:"server"`
	// Database 描述 PostgreSQL 连接地址、凭证与连接池限制。
	Database DatabaseConfig `mapstructure:"database"`
	// Auth 描述后台开发者鉴权使用的 JWT 与 Refresh Token 配置。
	Auth AuthConfig `mapstructure:"auth"`
	// Log 描述日志级别、输出路径与文件滚动策略。
	Log LogConfig `mapstructure:"log"`
}

// AppConfig 描述应用级元信息。
type AppConfig struct {
	// Name 指定应用名称，默认值为 `operation-admin-backend`，用于健康检查和启动日志。
	Name string `mapstructure:"name"`
}

// ServerConfig 描述 HTTP 服务监听和超时配置。
type ServerConfig struct {
	// Host 指定 HTTP 服务监听地址，默认值为 `0.0.0.0`，表示监听所有网卡。
	Host string `mapstructure:"host"`
	// Port 指定 HTTP 服务监听端口，默认值为 `8080`，取值范围必须在 1 到 65535 之间。
	Port int `mapstructure:"port"`
	// Mode 指定 Gin 运行模式，支持 `debug`、`release`、`test`，默认值为 `debug`。
	Mode string `mapstructure:"mode"`
	// ReadTimeout 指定单次请求读取超时，默认值为 `10s`，单位为 Go duration 字符串。
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
	// WriteTimeout 指定单次响应写出超时，默认值为 `15s`，单位为 Go duration 字符串。
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	// ShutdownTimeout 指定优雅停机等待时间，默认值为 `10s`，单位为 Go duration 字符串。
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig 描述 PostgreSQL 连接池配置。
type DatabaseConfig struct {
	// Host 指定 PostgreSQL 主机地址，默认值为 `127.0.0.1`。
	Host string `mapstructure:"host"`
	// Port 指定 PostgreSQL 端口，默认值为 `15432`，需要与根目录 Docker Compose 保持一致。
	Port int `mapstructure:"port"`
	// User 指定数据库用户名，默认值为 `postgres`。
	User string `mapstructure:"user"`
	// Password 指定数据库密码，默认值为 `postgres`，生产环境应改为安全凭证。
	Password string `mapstructure:"password"`
	// DBName 指定默认连接的数据库名称，默认值为 `pincermarket`。
	DBName string `mapstructure:"dbname"`
	// SSLMode 指定 PostgreSQL SSL 模式，开发环境默认值为 `disable`。
	SSLMode string `mapstructure:"sslmode"`
	// MaxIdleConns 指定连接池最大空闲连接数，默认值为 `5`。
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// MaxOpenConns 指定连接池最大打开连接数，默认值为 `20`，必须不小于空闲连接数。
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// ConnMaxLifetime 指定单条连接在池中的最大存活时长，默认值为 `30m`。
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// AuthConfig 描述后台开发者鉴权配置。
type AuthConfig struct {
	// Issuer 指定 JWT 签发方标识，默认值为 `operation-admin-backend`。
	Issuer string `mapstructure:"issuer"`
	// JWTSecret 指定 HS256 签名密钥，开发环境也必须至少为 32 个字符。
	JWTSecret string `mapstructure:"jwt_secret"`
	// AccessTokenTTL 指定 Access Token 有效期，默认值为 `2h`，单位为 Go duration 字符串。
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
	// RefreshTokenTTL 指定 Refresh Token 有效期，默认值为 `168h`，单位为 Go duration 字符串。
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// LogConfig 描述日志级别与日志落盘策略。
type LogConfig struct {
	// Level 指定日志最小输出级别，支持 `debug`、`info`、`warn`、`error`，默认值为 `info`。
	Level string `mapstructure:"level"`
	// Filename 指定滚动日志文件路径，默认值为 `logs/app.log`，相对路径会基于服务根目录解析。
	Filename string `mapstructure:"filename"`
	// MaxSizeMB 指定单个日志文件滚动前的最大体积，默认值为 `20`，单位为 MB。
	MaxSizeMB int `mapstructure:"max_size_mb"`
	// MaxBackups 指定保留的历史日志文件数量，默认值为 `5`。
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAgeDays 指定历史日志文件保留天数，默认值为 `14`，单位为天。
	MaxAgeDays int `mapstructure:"max_age_days"`
	// Compress 指定历史日志是否压缩，默认值为 `false`。
	Compress bool `mapstructure:"compress"`
}

// Load 负责加载配置文件，并允许环境变量覆盖默认值。
func Load(configPath string) (*Config, error) {
	// resolvedPath 用于记录最终选中的配置文件路径。
	resolvedPath, err := resolveConfigPath(configPath)
	if err != nil {
		return nil, err
	}

	// loader 是本次配置读取所使用的独立 Viper 实例。
	loader := viper.New()
	loader.SetConfigFile(resolvedPath)
	loader.SetConfigType("yaml")
	loader.SetEnvPrefix("OPERATION_ADMIN")
	loader.AutomaticEnv()

	// replacer 负责把嵌套配置键转换成环境变量可识别的下划线格式。
	replacer := strings.NewReplacer(".", "_")
	loader.SetEnvKeyReplacer(replacer)

	setDefaults(loader)

	if err := loader.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// loadedConfig 保存解析后的结构化配置。
	var loadedConfig Config

	if err := loader.Unmarshal(&loadedConfig, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		)
	}); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// runtimeRoot 用于把相对路径日志文件映射到服务根目录。
	runtimeRoot := resolveRuntimeRoot(resolvedPath)
	loadedConfig.Log.Filename = resolveRelativePath(runtimeRoot, loadedConfig.Log.Filename)

	if err := loadedConfig.Validate(); err != nil {
		return nil, err
	}

	return &loadedConfig, nil
}

// Address 负责组合服务监听地址。
func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DSN 负责组合 PostgreSQL 所需的数据源连接串。
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}

// Validate 负责校验整体配置是否满足服务启动要求。
func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}

	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	if err := c.Auth.Validate(); err != nil {
		return err
	}

	if err := c.Log.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate 负责校验应用元信息是否满足最小要求。
func (c AppConfig) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("app.name 不能为空")
	}

	return nil
}

// Validate 负责校验 HTTP 服务配置的端口、模式和超时是否合法。
func (c ServerConfig) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return errors.New("server.host 不能为空")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("server.port 超出合法范围: %d", c.Port)
	}

	if !isSupportedGinMode(c.Mode) {
		return fmt.Errorf("server.mode 不受支持: %s", c.Mode)
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("server.read_timeout 必须大于 0，当前为 %s", c.ReadTimeout)
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("server.write_timeout 必须大于 0，当前为 %s", c.WriteTimeout)
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("server.shutdown_timeout 必须大于 0，当前为 %s", c.ShutdownTimeout)
	}

	return nil
}

// Validate 负责校验数据库连接参数与连接池限制是否合法。
func (c DatabaseConfig) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return errors.New("database.host 不能为空")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("database.port 超出合法范围: %d", c.Port)
	}

	if strings.TrimSpace(c.User) == "" {
		return errors.New("database.user 不能为空")
	}

	if strings.TrimSpace(c.DBName) == "" {
		return errors.New("database.dbname 不能为空")
	}

	if strings.TrimSpace(c.SSLMode) == "" {
		return errors.New("database.sslmode 不能为空")
	}

	if c.MaxIdleConns < 0 {
		return fmt.Errorf("database.max_idle_conns 不能小于 0，当前为 %d", c.MaxIdleConns)
	}

	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("database.max_open_conns 必须大于 0，当前为 %d", c.MaxOpenConns)
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf(
			"database.max_idle_conns 不能大于 database.max_open_conns，当前为 %d > %d",
			c.MaxIdleConns,
			c.MaxOpenConns,
		)
	}

	if c.ConnMaxLifetime <= 0 {
		return fmt.Errorf("database.conn_max_lifetime 必须大于 0，当前为 %s", c.ConnMaxLifetime)
	}

	return nil
}

// Validate 负责校验鉴权配置是否满足 JWT 与 Refresh Token 的安全要求。
func (c AuthConfig) Validate() error {
	if strings.TrimSpace(c.Issuer) == "" {
		return errors.New("auth.issuer 不能为空")
	}

	if len(strings.TrimSpace(c.JWTSecret)) < 32 {
		return errors.New("auth.jwt_secret 长度不能少于 32 个字符")
	}

	if c.AccessTokenTTL <= 0 {
		return fmt.Errorf("auth.access_token_ttl 必须大于 0，当前为 %s", c.AccessTokenTTL)
	}

	if c.RefreshTokenTTL <= 0 {
		return fmt.Errorf("auth.refresh_token_ttl 必须大于 0，当前为 %s", c.RefreshTokenTTL)
	}

	return nil
}

// Validate 负责校验日志输出级别和滚动策略是否合法。
func (c LogConfig) Validate() error {
	if _, err := parseLogLevel(c.Level); err != nil {
		return err
	}

	if c.Filename == "" {
		return errors.New("log.filename 不能为空")
	}

	if c.MaxSizeMB <= 0 {
		return fmt.Errorf("log.max_size_mb 必须大于 0，当前为 %d", c.MaxSizeMB)
	}

	if c.MaxBackups < 0 {
		return fmt.Errorf("log.max_backups 不能小于 0，当前为 %d", c.MaxBackups)
	}

	if c.MaxAgeDays < 0 {
		return fmt.Errorf("log.max_age_days 不能小于 0，当前为 %d", c.MaxAgeDays)
	}

	return nil
}

// setDefaults 负责为缺省配置补齐默认值。
func setDefaults(loader *viper.Viper) {
	loader.SetDefault("app.name", "operation-admin-backend")

	loader.SetDefault("server.host", "0.0.0.0")
	loader.SetDefault("server.port", 8080)
	loader.SetDefault("server.mode", "debug")
	loader.SetDefault("server.read_timeout", "10s")
	loader.SetDefault("server.write_timeout", "15s")
	loader.SetDefault("server.shutdown_timeout", "10s")

	loader.SetDefault("database.host", "127.0.0.1")
	loader.SetDefault("database.port", 15432)
	loader.SetDefault("database.user", "postgres")
	loader.SetDefault("database.password", "postgres")
	loader.SetDefault("database.dbname", "pincermarket")
	loader.SetDefault("database.sslmode", "disable")
	loader.SetDefault("database.max_idle_conns", 5)
	loader.SetDefault("database.max_open_conns", 20)
	loader.SetDefault("database.conn_max_lifetime", "30m")

	loader.SetDefault("auth.issuer", "operation-admin-backend")
	loader.SetDefault("auth.jwt_secret", "dev-change-me-please-use-a-strong-secret")
	loader.SetDefault("auth.access_token_ttl", "2h")
	loader.SetDefault("auth.refresh_token_ttl", "168h")

	loader.SetDefault("log.level", "info")
	loader.SetDefault("log.filename", "logs/app.log")
	loader.SetDefault("log.max_size_mb", 20)
	loader.SetDefault("log.max_backups", 5)
	loader.SetDefault("log.max_age_days", 14)
	loader.SetDefault("log.compress", false)
}

// resolveConfigPath 负责自动探测默认配置文件位置。
func resolveConfigPath(configPath string) (string, error) {
	if configPath != "" {
		return configPath, nil
	}

	// currentDir 表示当前进程的工作目录，调试模式下可能位于 backend/cmd/server。
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	// searchRoots 描述从当前目录逐级向上回溯时需要检查的目录集合。
	searchRoots := buildSearchRoots(currentDir)

	for _, searchRoot := range searchRoots {
		// resolvedPath 表示在当前候选根目录下找到的配置文件路径。
		resolvedPath, found := findConfigInRoot(searchRoot)
		if found {
			return resolvedPath, nil
		}
	}

	return "", errors.New("未找到配置文件，请通过 -config 指定路径")
}

// resolveRelativePath 负责把配置中的相对路径映射到服务运行根目录。
func resolveRelativePath(baseDir, targetPath string) string {
	if targetPath == "" || filepath.IsAbs(targetPath) {
		return targetPath
	}

	return filepath.Join(baseDir, targetPath)
}

// resolveRuntimeRoot 负责根据配置文件位置推导服务根目录。
func resolveRuntimeRoot(configPath string) string {
	// configDir 代表配置文件所在目录。
	configDir := filepath.Dir(configPath)
	if filepath.Base(configDir) == "configs" {
		return filepath.Dir(configDir)
	}

	return configDir
}

// buildSearchRoots 负责从当前工作目录构建逐级向上的配置搜索路径集合。
func buildSearchRoots(startDir string) []string {
	// roots 按“当前目录到仓库根方向”的顺序保存待检查目录。
	roots := make([]string, 0, 8)
	// visited 用于避免符号链接或异常路径导致重复检查同一目录。
	visited := make(map[string]struct{})
	// currentDir 保存当前这一轮需要检查的目录。
	currentDir := startDir

	for {
		if _, exists := visited[currentDir]; exists {
			break
		}

		visited[currentDir] = struct{}{}
		roots = append(roots, currentDir)

		// parentDir 表示 currentDir 的父目录，用于继续向上回溯。
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}

		currentDir = parentDir
	}

	return roots
}

// findConfigInRoot 负责在单个候选根目录下查找支持的默认配置文件位置。
func findConfigInRoot(rootDir string) (string, bool) {
	// candidatePaths 描述当前仓库支持的两种默认配置文件位置。
	candidatePaths := []string{
		filepath.Join(rootDir, "configs", "config.yaml"),
		filepath.Join(rootDir, "backend", "configs", "config.yaml"),
	}

	for _, candidatePath := range candidatePaths {
		if _, err := os.Stat(candidatePath); err == nil {
			return candidatePath, true
		}
	}

	return "", false
}

// isSupportedGinMode 负责判断给定模式是否属于 Gin 支持的运行模式。
func isSupportedGinMode(mode string) bool {
	// normalizedMode 用于兼容大小写和额外空白字符差异。
	normalizedMode := strings.TrimSpace(strings.ToLower(mode))

	switch normalizedMode {
	case "debug", "release", "test":
		return true
	default:
		return false
	}
}

// parseLogLevel 负责在配置校验阶段验证日志级别是否可识别。
func parseLogLevel(level string) (string, error) {
	// normalizedLevel 用于统一处理日志级别大小写和空白差异。
	normalizedLevel := strings.TrimSpace(strings.ToLower(level))

	switch normalizedLevel {
	case "", "debug", "info", "warn", "warning", "error":
		return normalizedLevel, nil
	default:
		return "", fmt.Errorf("log.level 不受支持: %s", level)
	}
}
