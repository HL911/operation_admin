# Go 后端基础框架

## 技术栈

- `gin`：HTTP 路由与中间件
- `gorm`：数据库访问层
- `sonic`：JSON 编解码
- `viper`：配置加载与环境变量覆盖
- `zap`：结构化日志
- `lumberjack.v2`：日志文件滚动切割

## 目录结构

```text
backend/
├── cmd/server                 # 应用启动入口
├── configs/config.yaml        # 默认配置文件
├── internal/app               # 应用装配与生命周期管理
├── internal/config            # 配置加载与配置契约
├── internal/database          # Gorm / PostgreSQL 初始化
├── internal/http
│   ├── handler                # HTTP 处理器
│   ├── response               # Sonic JSON 响应封装
│   ├── router                 # 路由注册
├── internal/logger            # Zap + Lumberjack 日志初始化
├── logs                       # 本地日志目录
├── go.mod
└── Makefile
```

## 默认能力

- 提供 `GET /healthz` 健康检查接口
- 预留 `GET /admin/v1/healthz` 作为后台管理接口分组示例
- 启动时自动读取配置并初始化日志、数据库与 HTTP 服务
- 支持通过环境变量覆盖配置文件中的字段
- 日志默认同时输出到控制台与 `backend/logs/app.log`

## 本地启动

### 1. 启动 PostgreSQL

项目根目录已经提供本地 PostgreSQL：

```bash
docker compose up -d
```

默认连接信息：

- Host: `127.0.0.1`
- Port: `15432`
- Database: `pincermarket`
- User: `postgres`
- Password: `postgres`

### 2. 安装依赖

```bash
cd backend
make tidy
```

### 3. 启动服务

```bash
cd backend
make run
```

服务默认监听：

- `http://127.0.0.1:8080/healthz`
- `http://127.0.0.1:8080/admin/v1/healthz`

## 配置说明

默认配置文件路径：

- 在 `backend/` 目录执行时：`configs/config.yaml`
- 在仓库根目录执行时：`backend/configs/config.yaml`

支持使用环境变量覆盖，前缀统一为 `OPERATION_ADMIN`。

例如：

```bash
export OPERATION_ADMIN_SERVER_PORT=9090
export OPERATION_ADMIN_DATABASE_HOST=127.0.0.1
```

## 构建与测试

```bash
cd backend
make test
make build
```
