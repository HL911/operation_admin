# 运营后台前端骨架

## 项目定位

`frontend/` 用于承接运营后台的前端实现，当前阶段提供：

- `Next.js App Router` 的基础工程
- 面向后台管理场景的通用页面骨架
- 统一的 `/api/admin/*` 代理入口
- `Tailwind + shadcn/ui` 的主题与组件落点
- 面向 `v0` 的提示词与落库工作流约定

## 启动方式

### 1. 准备环境变量

复制 `.env.example` 并按需调整：

```bash
cp .env.example .env.local
```

### 2. 安装依赖

```bash
pnpm install
```

### 3. 启动开发服务器

```bash
pnpm dev
```

默认访问地址：

- 前端：`http://127.0.0.1:3000`
- 后端代理基址：`http://127.0.0.1:8080`

## 常用命令

- `pnpm dev`：启动开发服务器
- `pnpm lint`：执行 ESLint 检查
- `pnpm typecheck`：执行 TypeScript 类型检查
- `pnpm check`：串行执行 `lint` 与 `typecheck`
- `pnpm build`：执行生产构建

## 目录说明

```text
frontend/
├── app/                    # App Router 路由与 Route Handler
├── components/
│   ├── admin/              # 后台壳层与业务骨架组件
│   └── ui/                 # 可复用基础 UI 组件
├── features/
│   └── admin/              # 后台模块配置、代理逻辑与后续业务扩展
├── lib/                    # 通用工具函数
├── styles/                 # 全局主题变量与动效样式
└── public/                 # 静态资源
```

## 开发约定

- 用户可见文案、说明文档与模板统一使用中文。
- 每个函数都必须补充职责注释，变量需要有业务语义说明。
- 页面不要直接请求后端地址，统一走 `/api/admin/*`。
- `v0` 生成的代码必须经过组件拆分、类型补齐、中文注释补齐后才能入库。

## v0 使用建议

`v0` 适合负责：

- 页面视觉方向探索
- 复杂后台列表或详情布局初稿
- 表单编排与信息层级初稿

仓库内最终落地时，请遵循以下顺序：

1. 把 PM 文档拆成页面目标、字段、状态与操作流
2. 用中文 prompt 在 `v0` 中生成视觉稿或初版代码
3. 将产物拆成 `components/ui`、`components/admin` 与 `features/*`
4. 接入真实接口或 mock contract
5. 补齐中文注释、验收态与错误态
