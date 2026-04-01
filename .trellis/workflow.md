# 开发工作流

> 基于 [Effective Harnesses for Long-Running Agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents)

---

## 目录

1. [快速开始（先做这些）](#快速开始先做这些)
2. [工作流概览](#工作流概览)
3. [会话启动流程](#会话启动流程)
4. [开发流程](#开发流程)
5. [会话结束](#会话结束)
6. [文件说明](#文件说明)
7. [最佳实践](#最佳实践)

---

## 快速开始（先做这些）

### 第 0 步：初始化开发者身份（仅首次需要）

> **多开发者支持**：每位开发者或 Agent 都需要先初始化自己的身份

```bash
# 检查是否已初始化
python3 ./.trellis/scripts/get_developer.py

# 如果尚未初始化，运行：
python3 ./.trellis/scripts/init_developer.py <your-name>
# 示例：python3 ./.trellis/scripts/init_developer.py cursor-agent
```

这会创建：
- `.trellis/.developer` - 你的身份文件（已加入 gitignore，不提交）
- `.trellis/workspace/<your-name>/` - 你的个人工作区目录

**命名建议**：
- 人类开发者：使用你的名字，例如 `john-doe`
- Cursor AI：`cursor-agent` 或 `cursor-<task>`
- Claude Code：`claude-agent` 或 `claude-<task>`
- iFlow cli：`iflow-agent` 或 `iflow-<task>`

### 第 1 步：理解当前上下文

```bash
# 一条命令获取完整上下文
python3 ./.trellis/scripts/get_context.py

# 或手动分别查看：
python3 ./.trellis/scripts/get_developer.py      # 你的身份
python3 ./.trellis/scripts/task.py list          # 活动任务
git status && git log --oneline -10              # Git 状态
```

### 第 2 步：阅读项目规范 [必做]

**关键**：在编写任何代码前先阅读规范：

```bash
# 阅读前端规范索引（如适用）
cat .trellis/spec/frontend/index.md

# 阅读后端规范索引（如适用）
cat .trellis/spec/backend/index.md
```

**为什么两个都要看？**
- 了解完整的项目架构
- 掌握整个代码库的编码标准
- 看清前后端如何交互
- 了解整体代码质量要求

### 第 3 步：开始编码前，阅读具体规范（必做）

根据你的任务，阅读更细的规范文档：

**前端任务**：
```bash
cat .trellis/spec/frontend/hook-guidelines.md      # Hook 规范
cat .trellis/spec/frontend/component-guidelines.md # 组件规范
cat .trellis/spec/frontend/type-safety.md          # 类型规范
```

**后端任务**：
```bash
cat .trellis/spec/backend/database-guidelines.md   # 数据库规范
cat .trellis/spec/backend/type-safety.md           # 类型规范
cat .trellis/spec/backend/logging-guidelines.md    # 日志规范
```

---

## 工作流概览

### 核心原则

1. **先读后写** - 开始前先理解上下文
2. **遵循标准** - [!] **编码前必须阅读 `.trellis/spec/` 规范**
3. **增量开发** - 一次只完成一个任务
4. **及时记录** - 完成后立即更新追踪文件
5. **控制文档体积** - [!] **单个日志文档最多 2000 行**
6. **统一中文输出** - 所有输出文档、模板与用户可见说明默认使用中文
7. **强制补充注释** - 所有函数、变量、结构体字段、配置项都要有清晰注释

### 文件系统

```
.trellis/
|-- .developer           # 开发者身份（gitignored）
|-- scripts/
|   |-- __init__.py          # Python 包初始化
|   |-- common/              # 共享工具（Python）
|   |   |-- __init__.py
|   |   |-- paths.py         # 路径工具
|   |   |-- developer.py     # 开发者管理
|   |   +-- git_context.py   # Git 上下文实现
|   |-- multi_agent/         # 多 Agent 流水线脚本
|   |   |-- __init__.py
|   |   |-- start.py         # 启动 worktree agent
|   |   |-- status.py        # 监控 agent 状态
|   |   |-- create_pr.py     # 创建 PR
|   |   +-- cleanup.py       # 清理 worktree
|   |-- init_developer.py    # 初始化开发者身份
|   |-- get_developer.py     # 获取当前开发者名称
|   |-- task.py              # 管理任务
|   |-- get_context.py       # 获取会话上下文
|   +-- add_session.py       # 一键记录会话
|-- workspace/           # 开发者工作区
|   |-- index.md         # 工作区索引 + 会话模板
|   +-- {developer}/     # 每位开发者的目录
|       |-- index.md     # 个人索引（带 @@@auto 标记）
|       +-- journal-N.md # 日志文件（顺序编号）
|-- tasks/               # 任务跟踪
|   +-- {MM}-{DD}-{name}/
|       +-- task.json
|-- spec/                # [!] 编码前必须阅读
|   |-- frontend/        # 前端规范（如适用）
|   |   |-- index.md               # 从这里开始 - 规范索引
|   |   +-- *.md                   # 主题文档
|   |-- backend/         # 后端规范（如适用）
|   |   |-- index.md               # 从这里开始 - 规范索引
|   |   +-- *.md                   # 主题文档
|   +-- guides/          # 思考指南
|       |-- index.md                      # 指南索引
|       |-- cross-layer-thinking-guide.md # 跨层思考清单
|       +-- *.md                          # 其他指南
+-- workflow.md             # 本文档
```

---

## 会话启动流程

### 第 1 步：获取会话上下文

使用统一的上下文脚本：

```bash
# 一条命令获取全部上下文
python3 ./.trellis/scripts/get_context.py

# 或获取 JSON 格式
python3 ./.trellis/scripts/get_context.py --json
```

### 第 2 步：阅读开发规范 [!] 必做

**[!] 关键：编写任何代码前，必须先阅读规范**

根据你将要进行的开发类型，阅读对应规范：

**前端开发**（如适用）：
```bash
# 先读索引，再按任务阅读具体文档
cat .trellis/spec/frontend/index.md
```

**后端开发**（如适用）：
```bash
# 先读索引，再按任务阅读具体文档
cat .trellis/spec/backend/index.md
```

**跨层功能**：
```bash
# 适用于跨多个层级的功能
cat .trellis/spec/guides/cross-layer-thinking-guide.md
```

### 第 3 步：选择要开发的任务

使用任务管理脚本：

```bash
# 列出活动任务
python3 ./.trellis/scripts/task.py list

# 创建新任务（会创建包含 task.json 的目录）
python3 ./.trellis/scripts/task.py create "<title>" --slug <name>
```

---

## 开发流程

### 任务开发流

```
1. 创建或选择任务
   --> python3 ./.trellis/scripts/task.py create "<title>" --slug <name> 或 list

2. 按规范编写代码
   --> 阅读与你任务相关的 `.trellis/spec/` 文档
   --> 若涉及跨层：阅读 `.trellis/spec/guides/`

3. 自测
   --> 运行项目的 lint/test 命令（参见规范文档）
   --> 手动测试功能

4. 提交代码
   --> git add <files>
   --> git commit -m "type(scope): description"
       格式：feat/fix/docs/refactor/test/chore

5. 记录会话（一条命令）
   --> python3 ./.trellis/scripts/add_session.py --title "会话标题" --commit "提交哈希"
```

### 代码质量检查清单

**提交前必须通过**：
- [OK] Lint 检查通过
- [OK] 类型检查通过（如适用）
- [OK] 手动功能测试通过
- [OK] 所有新增或更新的文档均使用中文
- [OK] 所有函数都已添加职责注释
- [OK] 所有变量都有对应注释，或由紧邻的总注释统一说明
- [OK] 所有结构体及其字段都有对应注释
- [OK] 所有配置文件和配置结构体字段都有用途、默认值或约束说明

**项目特定检查**：
- 前端参见 `.trellis/spec/frontend/quality-guidelines.md`
- 后端参见 `.trellis/spec/backend/quality-guidelines.md`

---

## 会话结束

### 一键记录会话

代码提交后，使用：

```bash
python3 ./.trellis/scripts/add_session.py \
  --title "会话标题" \
  --commit "abc1234" \
  --summary "简要说明"
```

这会自动：
1. 检测当前日志文件
2. 若超过 2000 行则创建新文件
3. 追加会话内容
4. 更新 `index.md`（会话数、历史表）

### 结束前检查清单

使用 `/trellis:finish-work` 命令依次检查：
1. [OK] 所有代码已提交，且提交信息符合约定
2. [OK] 已通过 `add_session.py` 记录会话
3. [OK] 无 lint/test 错误
4. [OK] 工作目录干净，或已注明 WIP
5. [OK] 如有需要，已更新规范文档

---

## 文件说明

### 1. workspace/ - 开发者工作区

**用途**：记录每位 AI Agent 会话的工作内容

**结构**（多开发者支持）：
```
workspace/
|-- index.md              # 主索引
+-- {developer}/          # 每位开发者的目录
    |-- index.md          # 个人索引（带 @@@auto 标记）
    +-- journal-N.md      # 日志文件（顺序编号：1、2、3...）
```

**何时更新**：
- [OK] 每次会话结束时
- [OK] 完成重要任务时
- [OK] 修复重要 bug 时

### 2. spec/ - 开发规范

**用途**：沉淀一致开发方式的标准文档

**结构**（多文档形式）：
```
spec/
|-- frontend/           # 前端文档（如适用）
|   |-- index.md        # 从这里开始
|   +-- *.md            # 主题文档
|-- backend/            # 后端文档（如适用）
|   |-- index.md        # 从这里开始
|   +-- *.md            # 主题文档
+-- guides/             # 思考指南
    |-- index.md        # 从这里开始
    +-- *.md            # 各类指南文档
```

**何时更新**：
- [OK] 发现新的模式时
- [OK] 修复 bug 暴露出规范缺失时
- [OK] 确立新的团队约定时

### 3. 任务跟踪

每个任务是一个目录，内部包含 `task.json`：

```
tasks/
|-- 01-21-my-task/
|   +-- task.json
+-- archive/
    +-- 2026-01/
        +-- 01-15-old-task/
            +-- task.json
```

**常用命令**：
```bash
python3 ./.trellis/scripts/task.py create "<title>" [--slug <name>]   # 创建任务目录
python3 ./.trellis/scripts/task.py archive <name>  # 归档到 archive/{year-month}/
python3 ./.trellis/scripts/task.py list            # 列出活动任务
python3 ./.trellis/scripts/task.py list-archive    # 列出已归档任务
```

---

## 最佳实践

### [OK] 应该做的事

1. **会话开始前**：
   - 运行 `python3 ./.trellis/scripts/get_context.py` 获取完整上下文
   - [!] **必须阅读** 相关 `.trellis/spec/` 文档

2. **开发过程中**：
   - [!] **遵循** `.trellis/spec/` 规范
   - 对于跨层功能，使用 `/trellis:check-cross-layer`
   - 一次只开发一个任务
   - 频繁运行 lint 和测试
   - 输出文档统一使用中文
   - 为函数、变量、结构体字段、配置项补齐注释，不要留空白语义

3. **开发完成后**：
   - 使用 `/trellis:finish-work` 做收尾检查
   - 修完 bug 后，使用 `/trellis:break-loop` 做深入复盘
   - 由人工在测试通过后提交代码
   - 使用 `add_session.py` 记录进展

### [X] 不应该做的事

1. [!] **不要**跳过 `.trellis/spec/` 规范阅读
2. [!] **不要**让单个日志文件超过 2000 行
3. **不要**同时开发多个无关任务
4. **不要**在 lint/test 失败时提交代码
5. **不要**在获得新经验后忘记更新规范文档
6. [!] **不要**执行 `git commit` - AI 不应提交代码

---

## 快速参考

### 开发前必须阅读

| 任务类型 | 必读文档 |
|-----------|---------|
| 前端工作 | `frontend/index.md` → 相关文档 |
| 后端工作 | `backend/index.md` → 相关文档 |
| 跨层功能 | `guides/cross-layer-thinking-guide.md` |

### 提交信息约定

```bash
git commit -m "type(scope): description"
```

**Type**：feat、fix、docs、refactor、test、chore
**Scope**：模块名（例如 `auth`、`api`、`ui`）

### 常用命令

```bash
# 会话管理
python3 ./.trellis/scripts/get_context.py    # 获取完整上下文
python3 ./.trellis/scripts/add_session.py    # 记录会话

# 任务管理
python3 ./.trellis/scripts/task.py list      # 列出任务
python3 ./.trellis/scripts/task.py create "<title>" # 创建任务

# 斜杠命令
/trellis:finish-work          # 提交前检查清单
/trellis:break-loop           # 调试后深度复盘
/trellis:check-cross-layer    # 跨层校验
```

---

## 总结

遵循这套工作流可以确保：
- [OK] 多次会话之间保持连续性
- [OK] 代码质量保持一致
- [OK] 进度可追踪
- [OK] 经验可以沉淀到规范中
- [OK] 团队协作过程透明

**核心理念**：先读后写，遵循标准，及时记录，持续沉淀经验
