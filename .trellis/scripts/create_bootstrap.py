#!/usr/bin/env python3
"""
为首次使用创建 Bootstrap 引导任务。

该脚本会在第一次初始化 Trellis 后，创建一个引导任务，
帮助用户补齐项目开发规范。

用法：
    python3 create_bootstrap.py [project-type]

参数：
    project-type: frontend | backend | fullstack（默认：fullstack）

前置条件：
    - 必须存在 .trellis/.developer（请先运行 init_developer.py）

生成内容：
    .trellis/tasks/00-bootstrap-guidelines/
        - task.json    # 任务元数据
        - prd.md       # 任务说明与引导
"""

from __future__ import annotations

import json
import sys
from datetime import datetime
from pathlib import Path

from common.paths import (
    DIR_WORKFLOW,
    DIR_SCRIPTS,
    DIR_TASKS,
    get_repo_root,
    get_developer,
    get_tasks_dir,
    set_current_task,
)


# =============================================================================
# Constants
# =============================================================================

TASK_NAME = "00-bootstrap-guidelines"


# =============================================================================
# PRD Content
# =============================================================================

def write_prd_header() -> str:
    """Write PRD header section."""
    return """# 初始化：补齐项目开发规范

## 目标

欢迎使用 Trellis！这是你的第一个任务。

AI Agent 会通过 `.trellis/spec/` 理解你们项目自己的编码约定。
**如果模板为空，AI 很容易写出不符合项目风格的通用代码。**

把这些规范补齐属于一次性投入，但之后每次 AI 会话都会受益。

---

## 你的任务

基于你们**现有代码库**中的真实实现，补齐这些规范文件。
"""


def write_prd_backend_section() -> str:
    """Write PRD backend section."""
    return """

### 后端规范

| 文件 | 需要记录什么 |
|------|--------------|
| `.trellis/spec/backend/directory-structure.md` | 不同类型文件放在哪里（routes、services、utils 等） |
| `.trellis/spec/backend/database-guidelines.md` | ORM、迁移、查询模式、命名约定 |
| `.trellis/spec/backend/error-handling.md` | 错误如何捕获、记录与返回 |
| `.trellis/spec/backend/logging-guidelines.md` | 日志级别、格式、记录内容 |
| `.trellis/spec/backend/quality-guidelines.md` | 代码评审标准、测试要求 |
"""


def write_prd_frontend_section() -> str:
    """Write PRD frontend section."""
    return """

### 前端规范

| 文件 | 需要记录什么 |
|------|--------------|
| `.trellis/spec/frontend/directory-structure.md` | 组件、页面、Hook 的组织方式 |
| `.trellis/spec/frontend/component-guidelines.md` | 组件模式、props 约定 |
| `.trellis/spec/frontend/hook-guidelines.md` | 自定义 Hook 的命名与模式 |
| `.trellis/spec/frontend/state-management.md` | 状态库、使用模式、状态归属 |
| `.trellis/spec/frontend/type-safety.md` | TypeScript 约定、类型组织方式 |
| `.trellis/spec/frontend/quality-guidelines.md` | lint、测试、可访问性要求 |
"""


def write_prd_footer() -> str:
    """Write PRD footer section."""
    return """

### 思考指南（可选）

`.trellis/spec/guides/` 目录里已经有一些通用的思考指南。
如果有需要，你也可以把它们改成更贴合项目的版本。

---

## 如何完善规范

### 原则：记录现实，而不是理想

写下代码库**现在真实在发生的事情**，而不是你希望它以后变成什么样。
AI 需要匹配现有模式，而不是凭空引入新风格。

### 步骤

1. **查看现有代码** - 每种模式至少找 2 到 3 个示例
2. **记录模式** - 描述你观察到的规律
3. **附上文件路径** - 用真实文件做例子
4. **写清反模式** - 团队不推荐哪些做法

---

## 使用 AI 的建议

可以让 AI 协助分析代码库，例如：

- “看看我的代码库，并把你观察到的模式整理出来”
- “分析我的代码结构，并总结出约定”
- “查找错误处理模式，并把它们写成文档”

AI 会读取代码并帮助你完成整理。

---

## 完成检查清单

- [ ] 已补齐与你项目类型相关的规范
- [ ] 每份规范至少包含 2 到 3 个真实代码示例
- [ ] 已记录反模式

完成后：

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
"""


def write_prd(task_dir: Path, project_type: str) -> None:
    """Write prd.md file."""
    content = write_prd_header()

    if project_type == "frontend":
        content += write_prd_frontend_section()
    elif project_type == "backend":
        content += write_prd_backend_section()
    else:  # fullstack
        content += write_prd_backend_section()
        content += write_prd_frontend_section()

    content += write_prd_footer()

    prd_file = task_dir / "prd.md"
    prd_file.write_text(content, encoding="utf-8")


# =============================================================================
# Task JSON
# =============================================================================

def write_task_json(task_dir: Path, developer: str, project_type: str) -> None:
    """Write task.json file."""
    today = datetime.now().strftime("%Y-%m-%d")

    # Generate subtasks and related files based on project type
    if project_type == "frontend":
        subtasks = [
            {"name": "补齐前端规范", "status": "pending"},
            {"name": "补充代码示例", "status": "pending"},
        ]
        related_files = [".trellis/spec/frontend/"]
    elif project_type == "backend":
        subtasks = [
            {"name": "补齐后端规范", "status": "pending"},
            {"name": "补充代码示例", "status": "pending"},
        ]
        related_files = [".trellis/spec/backend/"]
    else:  # fullstack
        subtasks = [
            {"name": "补齐后端规范", "status": "pending"},
            {"name": "补齐前端规范", "status": "pending"},
            {"name": "补充代码示例", "status": "pending"},
        ]
        related_files = [".trellis/spec/backend/", ".trellis/spec/frontend/"]

    task_data = {
        "id": TASK_NAME,
        "name": "补齐开发规范",
        "description": "为 AI 助手补齐项目开发规范",
        "status": "in_progress",
        "dev_type": "docs",
        "priority": "P1",
        "creator": developer,
        "assignee": developer,
        "createdAt": today,
        "completedAt": None,
        "commit": None,
        "subtasks": subtasks,
        "children": [],
        "parent": None,
        "relatedFiles": related_files,
        "notes": f"由 trellis init 创建的首次引导任务（{project_type} 项目）",
        "meta": {},
    }

    task_json = task_dir / "task.json"
    task_json.write_text(json.dumps(task_data, indent=2, ensure_ascii=False), encoding="utf-8")


# =============================================================================
# Main
# =============================================================================

def main() -> int:
    """Main entry point."""
    # Parse project type argument
    project_type = "fullstack"
    if len(sys.argv) > 1:
        project_type = sys.argv[1]

    # Validate project type
    if project_type not in ("frontend", "backend", "fullstack"):
        print(f"未知项目类型：{project_type}，将默认使用 fullstack")
        project_type = "fullstack"

    repo_root = get_repo_root()
    developer = get_developer(repo_root)

    # Check developer initialized
    if not developer:
        print("错误：开发者尚未初始化")
        print(f"请运行：python3 ./{DIR_WORKFLOW}/{DIR_SCRIPTS}/init_developer.py <your-name>")
        return 1

    tasks_dir = get_tasks_dir(repo_root)
    task_dir = tasks_dir / TASK_NAME
    relative_path = f"{DIR_WORKFLOW}/{DIR_TASKS}/{TASK_NAME}"

    # Check if already exists
    if task_dir.exists():
        print(f"Bootstrap 任务已存在：{relative_path}")
        return 0

    # Create task directory
    task_dir.mkdir(parents=True, exist_ok=True)

    # Write files
    write_task_json(task_dir, developer, project_type)
    write_prd(task_dir, project_type)

    # Set as current task
    set_current_task(relative_path, repo_root)

    # Silent output - init command handles user-facing messages
    # Only output the task path for programmatic use
    print(relative_path)
    return 0


if __name__ == "__main__":
    sys.exit(main())
