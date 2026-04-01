# 初始化：补齐项目开发规范

## 目标

欢迎使用 Trellis！这是你的第一个任务。

AI Agent 会通过 `.trellis/spec/` 来理解**你们项目自己的编码约定**。
**如果模板是空的，AI 很容易写出不贴合你们项目风格的通用代码。**

把这些规范补齐属于一次性投入，但之后每次 AI 会话都会受益。

---

## 你的任务

基于你们**现有代码库**中的真实实现，补齐下面这些规范文件。

### 后端规范

| 文件 | 需要记录什么 |
|------|--------------|
| `.trellis/spec/backend/directory-structure.md` | 不同类型文件放在哪里（routes、services、utils 等） |
| `.trellis/spec/backend/database-guidelines.md` | ORM、迁移、查询模式、命名约定 |
| `.trellis/spec/backend/error-handling.md` | 错误如何捕获、记录与返回 |
| `.trellis/spec/backend/logging-guidelines.md` | 日志级别、格式、记录内容 |
| `.trellis/spec/backend/quality-guidelines.md` | 代码评审标准、测试要求 |

### 前端规范

| 文件 | 需要记录什么 |
|------|--------------|
| `.trellis/spec/frontend/directory-structure.md` | 组件、页面、Hook 的组织方式 |
| `.trellis/spec/frontend/component-guidelines.md` | 组件模式、props 约定 |
| `.trellis/spec/frontend/hook-guidelines.md` | 自定义 Hook 的命名与模式 |
| `.trellis/spec/frontend/state-management.md` | 状态库、使用模式、状态归属 |
| `.trellis/spec/frontend/type-safety.md` | TypeScript 约定、类型组织方式 |
| `.trellis/spec/frontend/quality-guidelines.md` | lint、测试、可访问性要求 |

### 思考指南（可选）

`.trellis/spec/guides/` 目录里已经有一些通用的最佳实践思考指南。
如果有需要，你也可以把它们改成更贴合本项目的版本。

---

## 如何完善这些规范

### 第 0 步：优先导入已有规范（推荐）

很多项目原本就已经在别处记录过编码约定。**动手从零写之前，先检查这些地方：**

| 文件 / 目录 | 来源工具 |
|------|------|
| `CLAUDE.md` / `CLAUDE.local.md` | Claude Code |
| `AGENTS.md` | Claude Code |
| `.cursorrules` | Cursor |
| `.cursor/rules/*.mdc` | Cursor（规则目录） |
| `.windsurfrules` | Windsurf |
| `.clinerules` | Cline |
| `.roomodes` | Roo Code |
| `.github/copilot-instructions.md` | GitHub Copilot |
| `.vscode/settings.json` → `github.copilot.chat.codeGeneration.instructions` | VS Code Copilot |
| `CONVENTIONS.md` / `.aider.conf.yml` | aider |
| `CONTRIBUTING.md` | 通用项目约定 |
| `.editorconfig` | 编辑器格式规则 |

如果这些文件存在，先阅读它们，再把与编码约定相关的内容抽取到对应的 `.trellis/spec/` 文件中。这样通常比从零开始写要省力很多。

### 第 1 步：分析代码库

可以让 AI 帮你从真实代码里提炼模式，例如：

- “读取现有配置文件（CLAUDE.md、.cursorrules 等），把编码约定提取到 `.trellis/spec/`”
- “分析我的代码库，并把你看到的模式整理成文档”
- “查找错误处理 / 组件 / API 模式，并把它们写进规范”

### 第 2 步：记录现实，而不是理想

写下代码库**现在真实在发生的事情**，而不是你希望它以后变成什么样。
AI 的目标是匹配现有模式，而不是凭空引入新风格。

- **查看现有代码**：每种模式至少找 2 到 3 个示例
- **附上文件路径**：引用真实文件做例子
- **写清反模式**：团队明确不推荐什么做法

---

## 完成检查清单

- [ ] 已补齐与你项目类型相关的规范
- [ ] 每份规范至少包含 2 到 3 个真实代码示例
- [ ] 已记录反模式

完成后运行：

```bash
python3 ./.trellis/scripts/task.py finish
python3 ./.trellis/scripts/task.py archive 00-bootstrap-guidelines
```

---

## 为什么这件事值得做

完成这个任务后：

1. AI 写出来的代码会更贴合你们项目风格
2. 相关 `/trellis:before-*-dev` 命令会注入真实上下文
3. `/trellis:check-*` 命令会按你们的实际标准做校验
4. 后续开发者（无论人类还是 AI）都会更快上手
