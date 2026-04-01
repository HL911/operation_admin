# 前端开发规范

> 当前仓库前端位于 `frontend/`，采用 `Next.js App Router + TypeScript + Tailwind + shadcn/ui` 作为统一基座。

---

## 概览

当前前端仍处于后台骨架阶段，已经明确的真实约定如下：

- 页面默认以服务端组件为先，只有交互需要时才落客户端组件
- 用户可见文案、注释、模板和说明统一使用中文
- 后端访问统一走 `frontend/app/api/admin/*`，不要在页面里直接请求后端地址
- 通用组件沉淀在 `frontend/components/ui` 与 `frontend/components/admin`
- `v0` 只负责视觉稿和初版代码，入库前必须做结构收敛与注释补齐

## 规范索引

| 规范 | 说明 | 状态 |
|------|------|------|
| [目录结构](./directory-structure.md) | 前端目录与模块边界 | 已落地 |
| [组件规范](./component-guidelines.md) | UI 组件、后台骨架组件和 props 约定 | 已落地 |
| [Hook 规范](./hook-guidelines.md) | 当前阶段 Hook 使用边界 | 已落地 |
| [状态管理](./state-management.md) | 服务端数据、URL 状态和本地状态策略 | 已落地 |
| [质量规范](./quality-guidelines.md) | lint、typecheck、文档与注释要求 | 已落地 |
| [类型安全](./type-safety.md) | TypeScript、接口契约与禁止模式 | 已落地 |

## 编码前检查

开始任何前端任务前，至少确认以下几点：

1. 目标页面是放在 `app/` 路由层，还是先沉淀成可复用组件
2. 该功能是否应该直接接后端代理，还是先用 mock contract 占位
3. 是否需要新增 `features/*` 配置、类型或代理工具，而不是把数据写死在组件里
4. 产物是否来自 `v0`；如果是，是否已经拆解成共享组件和业务组件
5. 所有函数、变量、配置说明是否都补齐了中文注释

## 相关文档

- `frontend/README.md`：前端本地启动、目录和命令说明
- `docs/frontend-v0-workflow.md`：`v0` 提示词与入库工作流
