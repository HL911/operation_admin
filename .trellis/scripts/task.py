#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Task Management Script for Multi-Agent Pipeline.

Usage:
    python3 task.py create "<title>" [--slug <name>] [--assignee <dev>] [--priority P0|P1|P2|P3] [--parent <dir>]
    python3 task.py init-context <dir> <type>   # Initialize jsonl files
    python3 task.py add-context <dir> <file> <path> [reason] # Add jsonl entry
    python3 task.py validate <dir>              # Validate jsonl files
    python3 task.py list-context <dir>          # List jsonl entries
    python3 task.py start <dir>                 # Set as current task
    python3 task.py finish                      # Clear current task
    python3 task.py set-branch <dir> <branch>   # Set git branch
    python3 task.py set-base-branch <dir> <branch>  # Set PR target branch
    python3 task.py set-scope <dir> <scope>     # Set scope for PR title
    python3 task.py create-pr [dir] [--dry-run] # Create PR from task
    python3 task.py archive <task-name>         # Archive completed task
    python3 task.py list                        # List active tasks
    python3 task.py list-archive [month]        # List archived tasks
    python3 task.py add-subtask <parent-dir> <child-dir>     # Link child to parent
    python3 task.py remove-subtask <parent-dir> <child-dir>  # Unlink child from parent
"""

from __future__ import annotations

import sys

# IMPORTANT: Force stdout to use UTF-8 on Windows
# This fixes UnicodeEncodeError when outputting non-ASCII characters
if sys.platform == "win32":
    import io as _io
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")  # type: ignore[union-attr]
    elif hasattr(sys.stdout, "detach"):
        sys.stdout = _io.TextIOWrapper(sys.stdout.detach(), encoding="utf-8", errors="replace")  # type: ignore[union-attr]

import argparse
import json
import re
import sys
from datetime import datetime
from pathlib import Path

from common.cli_adapter import get_cli_adapter_auto
from common.git_context import _run_git_command
from common.paths import (
    DIR_WORKFLOW,
    DIR_TASKS,
    DIR_SPEC,
    DIR_ARCHIVE,
    FILE_TASK_JSON,
    get_repo_root,
    get_developer,
    get_tasks_dir,
    get_current_task,
    set_current_task,
    clear_current_task,
    generate_task_date_prefix,
)
from common.task_utils import (
    find_task_by_name,
    archive_task_complete,
)
from common.config import get_hooks


# =============================================================================
# Colors
# =============================================================================

class Colors:
    RED = "\033[0;31m"
    GREEN = "\033[0;32m"
    YELLOW = "\033[1;33m"
    BLUE = "\033[0;34m"
    CYAN = "\033[0;36m"
    NC = "\033[0m"


def colored(text: str, color: str) -> str:
    """Apply color to text."""
    return f"{color}{text}{Colors.NC}"


STATUS_LABELS = {
    "planning": "规划中",
    "in_progress": "进行中",
    "review": "评审中",
    "completed": "已完成",
    "done": "已完成",
    "blocked": "阻塞中",
    "unknown": "未知",
}


def display_status(status: str) -> str:
    """将内部状态值转换为中文显示文本。"""
    return STATUS_LABELS.get(status, status)


# =============================================================================
# Lifecycle Hooks
# =============================================================================

def _run_hooks(event: str, task_json_path: Path, repo_root: Path) -> None:
    """Run lifecycle hooks for an event.

    Args:
        event: Event name (e.g. "after_create").
        task_json_path: Absolute path to the task's task.json.
        repo_root: Repository root for cwd and config lookup.
    """
    import os
    import subprocess

    commands = get_hooks(event, repo_root)
    if not commands:
        return

    env = {**os.environ, "TASK_JSON_PATH": str(task_json_path)}

    for cmd in commands:
        try:
            result = subprocess.run(
                cmd,
                shell=True,
                cwd=repo_root,
                env=env,
                capture_output=True,
                text=True,
                encoding="utf-8",
                errors="replace",
            )
            if result.returncode != 0:
                print(
                    colored(f"[WARN] Hook failed ({event}): {cmd}", Colors.YELLOW),
                    file=sys.stderr,
                )
                if result.stderr.strip():
                    print(f"  {result.stderr.strip()}", file=sys.stderr)
        except Exception as e:
            print(
                colored(f"[WARN] Hook error ({event}): {cmd} — {e}", Colors.YELLOW),
                file=sys.stderr,
            )


# =============================================================================
# Helper Functions
# =============================================================================

def _read_json_file(path: Path) -> dict | None:
    """Read and parse a JSON file."""
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (FileNotFoundError, json.JSONDecodeError, OSError):
        return None


def _write_json_file(path: Path, data: dict) -> bool:
    """Write dict to JSON file."""
    try:
        path.write_text(json.dumps(data, indent=2, ensure_ascii=False), encoding="utf-8")
        return True
    except (OSError, IOError):
        return False


def _slugify(title: str) -> str:
    """Convert title to slug (only works with ASCII)."""
    result = title.lower()
    result = re.sub(r"[^a-z0-9]", "-", result)
    result = re.sub(r"-+", "-", result)
    result = result.strip("-")
    return result


def _resolve_task_dir(target_dir: str, repo_root: Path) -> Path:
    """Resolve task directory to absolute path.

    Supports:
    - Absolute path: /path/to/task
    - Relative path: .trellis/tasks/01-31-my-task
    - Task name: my-task (uses find_task_by_name for lookup)
    """
    if not target_dir:
        return Path()

    # Absolute path
    if target_dir.startswith("/"):
        return Path(target_dir)

    # Relative path (contains path separator or starts with .trellis)
    if "/" in target_dir or target_dir.startswith(".trellis"):
        return repo_root / target_dir

    # Task name - try to find in tasks directory
    tasks_dir = get_tasks_dir(repo_root)
    found = find_task_by_name(target_dir, tasks_dir)
    if found:
        return found

    # Fallback to treating as relative path
    return repo_root / target_dir


# =============================================================================
# JSONL Default Content Generators
# =============================================================================

def get_implement_base() -> list[dict]:
    """Get base implement context entries."""
    return [
        {"file": f"{DIR_WORKFLOW}/workflow.md", "reason": "Project workflow and conventions"},
    ]


def get_implement_backend() -> list[dict]:
    """Get backend implement context entries."""
    return [
        {"file": f"{DIR_WORKFLOW}/{DIR_SPEC}/backend/index.md", "reason": "Backend development guide"},
    ]


def get_implement_frontend() -> list[dict]:
    """Get frontend implement context entries."""
    return [
        {"file": f"{DIR_WORKFLOW}/{DIR_SPEC}/frontend/index.md", "reason": "Frontend development guide"},
    ]


def get_check_context(dev_type: str, repo_root: Path) -> list[dict]:
    """Get check context entries."""
    adapter = get_cli_adapter_auto(repo_root)

    entries = [
        {"file": adapter.get_trellis_command_path("finish-work"), "reason": "Finish work checklist"},
    ]

    if dev_type in ("backend", "fullstack"):
        entries.append({"file": adapter.get_trellis_command_path("check-backend"), "reason": "Backend check spec"})
    if dev_type in ("frontend", "fullstack"):
        entries.append({"file": adapter.get_trellis_command_path("check-frontend"), "reason": "Frontend check spec"})

    return entries


def get_debug_context(dev_type: str, repo_root: Path) -> list[dict]:
    """Get debug context entries."""
    adapter = get_cli_adapter_auto(repo_root)

    entries: list[dict] = []

    if dev_type in ("backend", "fullstack"):
        entries.append({"file": adapter.get_trellis_command_path("check-backend"), "reason": "Backend check spec"})
    if dev_type in ("frontend", "fullstack"):
        entries.append({"file": adapter.get_trellis_command_path("check-frontend"), "reason": "Frontend check spec"})

    return entries


def _write_jsonl(path: Path, entries: list[dict]) -> None:
    """Write entries to JSONL file."""
    lines = [json.dumps(entry, ensure_ascii=False) for entry in entries]
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


# =============================================================================
# Task Operations
# =============================================================================

def ensure_tasks_dir(repo_root: Path) -> Path:
    """Ensure tasks directory exists."""
    tasks_dir = get_tasks_dir(repo_root)
    archive_dir = tasks_dir / "archive"

    if not tasks_dir.exists():
        tasks_dir.mkdir(parents=True)
        print(colored(f"已创建任务目录：{tasks_dir}", Colors.GREEN), file=sys.stderr)

    if not archive_dir.exists():
        archive_dir.mkdir(parents=True)

    return tasks_dir


# =============================================================================
# Command: create
# =============================================================================

def cmd_create(args: argparse.Namespace) -> int:
    """Create a new task."""
    repo_root = get_repo_root()

    if not args.title:
        print(colored("错误：标题不能为空", Colors.RED), file=sys.stderr)
        return 1

    # Default assignee to current developer
    assignee = args.assignee
    if not assignee:
        assignee = get_developer(repo_root)
        if not assignee:
            print(colored("错误：尚未设置开发者。请先运行 init_developer.py，或使用 --assignee", Colors.RED), file=sys.stderr)
            return 1

    ensure_tasks_dir(repo_root)

    # Get current developer as creator
    creator = get_developer(repo_root) or assignee

    # Generate slug if not provided
    slug = args.slug or _slugify(args.title)
    if not slug:
        print(colored("错误：无法根据标题生成 slug", Colors.RED), file=sys.stderr)
        return 1

    # Create task directory with MM-DD-slug format
    tasks_dir = get_tasks_dir(repo_root)
    date_prefix = generate_task_date_prefix()
    dir_name = f"{date_prefix}-{slug}"
    task_dir = tasks_dir / dir_name
    task_json_path = task_dir / FILE_TASK_JSON

    if task_dir.exists():
        print(colored(f"警告：任务目录已存在：{dir_name}", Colors.YELLOW), file=sys.stderr)
    else:
        task_dir.mkdir(parents=True)

    today = datetime.now().strftime("%Y-%m-%d")

    # Record current branch as base_branch (PR target)
    _, branch_out, _ = _run_git_command(["branch", "--show-current"], cwd=repo_root)
    current_branch = branch_out.strip() or "main"

    task_data = {
        "id": slug,
        "name": slug,
        "title": args.title,
        "description": args.description or "",
        "status": "planning",
        "dev_type": None,
        "scope": None,
        "priority": args.priority,
        "creator": creator,
        "assignee": assignee,
        "createdAt": today,
        "completedAt": None,
        "branch": None,
        "base_branch": current_branch,
        "worktree_path": None,
        "current_phase": 0,
        "next_action": [
            {"phase": 1, "action": "implement"},
            {"phase": 2, "action": "check"},
            {"phase": 3, "action": "finish"},
            {"phase": 4, "action": "create-pr"},
        ],
        "commit": None,
        "pr_url": None,
        "subtasks": [],
        "children": [],
        "parent": None,
        "relatedFiles": [],
        "notes": "",
        "meta": {},
    }

    _write_json_file(task_json_path, task_data)

    # Handle --parent: establish bidirectional link
    if args.parent:
        parent_dir = _resolve_task_dir(args.parent, repo_root)
        parent_json_path = parent_dir / FILE_TASK_JSON
        if not parent_json_path.is_file():
            print(colored(f"警告：未找到父任务 task.json：{args.parent}", Colors.YELLOW), file=sys.stderr)
        else:
            parent_data = _read_json_file(parent_json_path)
            if parent_data:
                # Add child to parent's children list
                parent_children = parent_data.get("children", [])
                if dir_name not in parent_children:
                    parent_children.append(dir_name)
                    parent_data["children"] = parent_children
                    _write_json_file(parent_json_path, parent_data)

                # Set parent in child's task.json
                task_data["parent"] = parent_dir.name
                _write_json_file(task_json_path, task_data)

                print(colored(f"已关联为子任务：{parent_dir.name}", Colors.GREEN), file=sys.stderr)

    print(colored(f"已创建任务：{dir_name}", Colors.GREEN), file=sys.stderr)
    print("", file=sys.stderr)
    print(colored("下一步：", Colors.BLUE), file=sys.stderr)
    print("  1. 创建包含需求的 prd.md", file=sys.stderr)
    print("  2. 运行：python3 task.py init-context <dir> <dev_type>", file=sys.stderr)
    print("  3. 运行：python3 task.py start <dir>", file=sys.stderr)
    print("", file=sys.stderr)

    # Output relative path for script chaining
    print(f"{DIR_WORKFLOW}/{DIR_TASKS}/{dir_name}")

    _run_hooks("after_create", task_json_path, repo_root)
    return 0


# =============================================================================
# Command: init-context
# =============================================================================

def cmd_init_context(args: argparse.Namespace) -> int:
    """Initialize JSONL context files for a task."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)
    dev_type = args.type

    if not dev_type:
        print(colored("错误：缺少参数", Colors.RED))
        print("用法：python3 task.py init-context <task-dir> <dev_type>")
        print("  dev_type：backend | frontend | fullstack | test | docs")
        return 1

    if not target_dir.is_dir():
        print(colored(f"错误：未找到目录：{target_dir}", Colors.RED))
        return 1

    print(colored("=== 正在初始化 Agent 上下文文件 ===", Colors.BLUE))
    print(f"目标目录：{target_dir}")
    print(f"开发类型：{dev_type}")
    print()

    # implement.jsonl
    print(colored("正在创建 implement.jsonl...", Colors.CYAN))
    implement_entries = get_implement_base()
    if dev_type in ("backend", "test"):
        implement_entries.extend(get_implement_backend())
    elif dev_type == "frontend":
        implement_entries.extend(get_implement_frontend())
    elif dev_type == "fullstack":
        implement_entries.extend(get_implement_backend())
        implement_entries.extend(get_implement_frontend())

    implement_file = target_dir / "implement.jsonl"
    _write_jsonl(implement_file, implement_entries)
    print(f"  {colored('✓', Colors.GREEN)} {len(implement_entries)} entries")

    # check.jsonl
    print(colored("正在创建 check.jsonl...", Colors.CYAN))
    check_entries = get_check_context(dev_type, repo_root)
    check_file = target_dir / "check.jsonl"
    _write_jsonl(check_file, check_entries)
    print(f"  {colored('✓', Colors.GREEN)} {len(check_entries)} entries")

    # debug.jsonl
    print(colored("正在创建 debug.jsonl...", Colors.CYAN))
    debug_entries = get_debug_context(dev_type, repo_root)
    debug_file = target_dir / "debug.jsonl"
    _write_jsonl(debug_file, debug_entries)
    print(f"  {colored('✓', Colors.GREEN)} {len(debug_entries)} entries")

    print()
    print(colored("✓ 所有上下文文件已创建", Colors.GREEN))
    print()
    print(colored("下一步：", Colors.BLUE))
    print("  1. 添加任务专属规范：python3 task.py add-context <dir> <jsonl> <path>")
    print("  2. 设为当前任务：python3 task.py start <dir>")

    return 0


# =============================================================================
# Command: add-context
# =============================================================================

def cmd_add_context(args: argparse.Namespace) -> int:
    """Add entry to JSONL context file."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)

    jsonl_name = args.file
    path = args.path
    reason = args.reason or "手动添加"

    if not target_dir.is_dir():
        print(colored(f"错误：未找到目录：{target_dir}", Colors.RED))
        return 1

    # Support shorthand
    if not jsonl_name.endswith(".jsonl"):
        jsonl_name = f"{jsonl_name}.jsonl"

    jsonl_file = target_dir / jsonl_name
    full_path = repo_root / path

    entry_type = "file"
    if full_path.is_dir():
        entry_type = "directory"
        if not path.endswith("/"):
            path = f"{path}/"
    elif not full_path.is_file():
        print(colored(f"错误：未找到路径：{path}", Colors.RED))
        return 1

    # Check if already exists
    if jsonl_file.is_file():
        content = jsonl_file.read_text(encoding="utf-8")
        if f'"{path}"' in content:
            print(colored(f"警告：{path} 的条目已存在", Colors.YELLOW))
            return 0

    # Add entry
    entry: dict
    if entry_type == "directory":
        entry = {"file": path, "type": "directory", "reason": reason}
    else:
        entry = {"file": path, "reason": reason}

    with jsonl_file.open("a", encoding="utf-8") as f:
        f.write(json.dumps(entry, ensure_ascii=False) + "\n")

    display_type = "目录" if entry_type == "directory" else "文件"
    print(colored(f"已添加{display_type}：{path}", Colors.GREEN))
    return 0


# =============================================================================
# Command: validate
# =============================================================================

def cmd_validate(args: argparse.Namespace) -> int:
    """Validate JSONL context files."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)

    if not target_dir.is_dir():
        print(colored("错误：必须提供任务目录", Colors.RED))
        return 1

    print(colored("=== 正在校验上下文文件 ===", Colors.BLUE))
    print(f"目标目录：{target_dir}")
    print()

    total_errors = 0
    for jsonl_name in ["implement.jsonl", "check.jsonl", "debug.jsonl"]:
        jsonl_file = target_dir / jsonl_name
        errors = _validate_jsonl(jsonl_file, repo_root)
        total_errors += errors

    print()
    if total_errors == 0:
        print(colored("✓ 所有校验均已通过", Colors.GREEN))
        return 0
    else:
        print(colored(f"✗ 校验失败（{total_errors} 个错误）", Colors.RED))
        return 1


def _validate_jsonl(jsonl_file: Path, repo_root: Path) -> int:
    """Validate a single JSONL file."""
    file_name = jsonl_file.name
    errors = 0

    if not jsonl_file.is_file():
        print(f"  {colored(f'{file_name}: 未找到（已跳过）', Colors.YELLOW)}")
        return 0

    line_num = 0
    for line in jsonl_file.read_text(encoding="utf-8").splitlines():
        line_num += 1
        if not line.strip():
            continue

        try:
            data = json.loads(line)
        except json.JSONDecodeError:
            print(f"  {colored(f'{file_name}:{line_num}: JSON 无效', Colors.RED)}")
            errors += 1
            continue

        file_path = data.get("file")
        entry_type = data.get("type", "file")

        if not file_path:
            print(f"  {colored(f'{file_name}:{line_num}: 缺少 file 字段', Colors.RED)}")
            errors += 1
            continue

        full_path = repo_root / file_path
        if entry_type == "directory":
            if not full_path.is_dir():
                print(f"  {colored(f'{file_name}:{line_num}: 未找到目录：{file_path}', Colors.RED)}")
                errors += 1
        else:
            if not full_path.is_file():
                print(f"  {colored(f'{file_name}:{line_num}: 未找到文件：{file_path}', Colors.RED)}")
                errors += 1

    if errors == 0:
        print(f"  {colored(f'{file_name}: ✓（{line_num} 条）', Colors.GREEN)}")
    else:
        print(f"  {colored(f'{file_name}: ✗（{errors} 个错误）', Colors.RED)}")

    return errors


# =============================================================================
# Command: list-context
# =============================================================================

def cmd_list_context(args: argparse.Namespace) -> int:
    """List JSONL context entries."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)

    if not target_dir.is_dir():
        print(colored("错误：必须提供任务目录", Colors.RED))
        return 1

    print(colored("=== 上下文文件 ===", Colors.BLUE))
    print()

    for jsonl_name in ["implement.jsonl", "check.jsonl", "debug.jsonl"]:
        jsonl_file = target_dir / jsonl_name
        if not jsonl_file.is_file():
            continue

        print(colored(f"[{jsonl_name}]", Colors.CYAN))

        count = 0
        for line in jsonl_file.read_text(encoding="utf-8").splitlines():
            if not line.strip():
                continue

            try:
                data = json.loads(line)
            except json.JSONDecodeError:
                continue

            count += 1
            file_path = data.get("file", "?")
            entry_type = data.get("type", "file")
            reason = data.get("reason", "-")

            if entry_type == "directory":
                print(f"  {colored(f'{count}.', Colors.GREEN)} [DIR] {file_path}")
            else:
                print(f"  {colored(f'{count}.', Colors.GREEN)} {file_path}")
            print(f"     {colored('→', Colors.YELLOW)} {reason}")

        print()

    return 0


# =============================================================================
# Command: start / finish
# =============================================================================

def cmd_start(args: argparse.Namespace) -> int:
    """Set current task."""
    repo_root = get_repo_root()
    task_input = args.dir

    if not task_input:
        print(colored("错误：必须提供任务目录或任务名", Colors.RED))
        return 1

    # Resolve task directory (supports task name, relative path, or absolute path)
    full_path = _resolve_task_dir(task_input, repo_root)

    if not full_path.is_dir():
        print(colored(f"错误：未找到任务：{task_input}", Colors.RED))
        print("提示：可使用任务名（如 `my-task`）或完整路径（如 `.trellis/tasks/01-31-my-task`）")
        return 1

    # Convert to relative path for storage
    try:
        task_dir = str(full_path.relative_to(repo_root))
    except ValueError:
        task_dir = str(full_path)

    if set_current_task(task_dir, repo_root):
        print(colored(f"✓ 当前任务已设置为：{task_dir}", Colors.GREEN))
        print()
        print(colored("Hook 现在会从该任务的 jsonl 文件中注入上下文。", Colors.BLUE))

        task_json_path = full_path / FILE_TASK_JSON
        _run_hooks("after_start", task_json_path, repo_root)
        return 0
    else:
        print(colored("错误：设置当前任务失败", Colors.RED))
        return 1


def cmd_finish(args: argparse.Namespace) -> int:
    """Clear current task."""
    repo_root = get_repo_root()
    current = get_current_task(repo_root)

    if not current:
        print(colored("当前没有设置任务", Colors.YELLOW))
        return 0

    # Resolve task.json path before clearing
    task_json_path = repo_root / current / FILE_TASK_JSON

    clear_current_task(repo_root)
    print(colored(f"✓ 已清除当前任务（之前为：{current}）", Colors.GREEN))

    if task_json_path.is_file():
        _run_hooks("after_finish", task_json_path, repo_root)
    return 0


# =============================================================================
# Command: archive
# =============================================================================

def cmd_archive(args: argparse.Namespace) -> int:
    """Archive completed task."""
    repo_root = get_repo_root()
    task_name = args.name

    if not task_name:
        print(colored("错误：任务名不能为空", Colors.RED), file=sys.stderr)
        return 1

    tasks_dir = get_tasks_dir(repo_root)

    # Find task directory
    task_dir = find_task_by_name(task_name, tasks_dir)

    if not task_dir or not task_dir.is_dir():
        print(colored(f"错误：未找到任务：{task_name}", Colors.RED), file=sys.stderr)
        print("活动任务：", file=sys.stderr)
        cmd_list(argparse.Namespace(mine=False, status=None))
        return 1

    dir_name = task_dir.name
    task_json_path = task_dir / FILE_TASK_JSON

    # Update status before archiving
    today = datetime.now().strftime("%Y-%m-%d")
    if task_json_path.is_file():
        data = _read_json_file(task_json_path)
        if data:
            data["status"] = "completed"
            data["completedAt"] = today
            _write_json_file(task_json_path, data)

            # Handle subtask relationships on archive
            task_parent = data.get("parent")
            task_children = data.get("children", [])

            # If this is a child, remove from parent's children list
            if task_parent:
                parent_dir = find_task_by_name(task_parent, tasks_dir)
                if parent_dir:
                    parent_json = parent_dir / FILE_TASK_JSON
                    if parent_json.is_file():
                        parent_data = _read_json_file(parent_json)
                        if parent_data:
                            parent_children = parent_data.get("children", [])
                            if dir_name in parent_children:
                                parent_children.remove(dir_name)
                                parent_data["children"] = parent_children
                                _write_json_file(parent_json, parent_data)

            # If this is a parent, clear parent field in all children
            if task_children:
                for child_name in task_children:
                    child_dir_path = find_task_by_name(child_name, tasks_dir)
                    if child_dir_path:
                        child_json = child_dir_path / FILE_TASK_JSON
                        if child_json.is_file():
                            child_data = _read_json_file(child_json)
                            if child_data:
                                child_data["parent"] = None
                                _write_json_file(child_json, child_data)

    # Clear if current task
    current = get_current_task(repo_root)
    if current and dir_name in current:
        clear_current_task(repo_root)

    # Archive
    result = archive_task_complete(task_dir, repo_root)
    if "archived_to" in result:
        archive_dest = Path(result["archived_to"])
        year_month = archive_dest.parent.name
        print(colored(f"Archived: {dir_name} -> archive/{year_month}/", Colors.GREEN), file=sys.stderr)

        # Auto-commit unless --no-commit
        if not getattr(args, "no_commit", False):
            _auto_commit_archive(dir_name, repo_root)

        # Return the archive path
        print(f"{DIR_WORKFLOW}/{DIR_TASKS}/{DIR_ARCHIVE}/{year_month}/{dir_name}")

        # Run hooks with the archived path
        archived_json = archive_dest / FILE_TASK_JSON
        _run_hooks("after_archive", archived_json, repo_root)
        return 0

    return 1


def _auto_commit_archive(task_name: str, repo_root: Path) -> None:
    """Stage .trellis/tasks/ changes and commit after archive."""
    tasks_rel = f"{DIR_WORKFLOW}/{DIR_TASKS}"
    _run_git_command(["add", "-A", tasks_rel], cwd=repo_root)

    # Check if there are staged changes
    rc, _, _ = _run_git_command(
        ["diff", "--cached", "--quiet", "--", tasks_rel], cwd=repo_root
    )
    if rc == 0:
        print("[OK] No task changes to commit.", file=sys.stderr)
        return

    commit_msg = f"chore(task): archive {task_name}"
    rc, _, err = _run_git_command(["commit", "-m", commit_msg], cwd=repo_root)
    if rc == 0:
        print(f"[OK] Auto-committed: {commit_msg}", file=sys.stderr)
    else:
        print(f"[WARN] Auto-commit failed: {err.strip()}", file=sys.stderr)


# =============================================================================
# Command: add-subtask
# =============================================================================

def cmd_add_subtask(args: argparse.Namespace) -> int:
    """Link a child task to a parent task."""
    repo_root = get_repo_root()

    parent_dir = _resolve_task_dir(args.parent_dir, repo_root)
    child_dir = _resolve_task_dir(args.child_dir, repo_root)

    parent_json_path = parent_dir / FILE_TASK_JSON
    child_json_path = child_dir / FILE_TASK_JSON

    if not parent_json_path.is_file():
        print(colored(f"错误：未找到父任务 task.json：{args.parent_dir}", Colors.RED), file=sys.stderr)
        return 1

    if not child_json_path.is_file():
        print(colored(f"错误：未找到子任务 task.json：{args.child_dir}", Colors.RED), file=sys.stderr)
        return 1

    parent_data = _read_json_file(parent_json_path)
    child_data = _read_json_file(child_json_path)

    if not parent_data or not child_data:
        print(colored("错误：读取 task.json 失败", Colors.RED), file=sys.stderr)
        return 1

    # Check if child already has a parent
    existing_parent = child_data.get("parent")
    if existing_parent:
        print(colored(f"错误：子任务已经有父任务：{existing_parent}", Colors.RED), file=sys.stderr)
        return 1

    # Add child to parent's children list
    parent_children = parent_data.get("children", [])
    child_dir_name = child_dir.name
    if child_dir_name not in parent_children:
        parent_children.append(child_dir_name)
        parent_data["children"] = parent_children

    # Set parent in child's task.json
    child_data["parent"] = parent_dir.name

    # Write both
    _write_json_file(parent_json_path, parent_data)
    _write_json_file(child_json_path, child_data)

    print(colored(f"已关联：{child_dir.name} -> {parent_dir.name}", Colors.GREEN), file=sys.stderr)
    return 0


# =============================================================================
# Command: remove-subtask
# =============================================================================

def cmd_remove_subtask(args: argparse.Namespace) -> int:
    """Unlink a child task from a parent task."""
    repo_root = get_repo_root()

    parent_dir = _resolve_task_dir(args.parent_dir, repo_root)
    child_dir = _resolve_task_dir(args.child_dir, repo_root)

    parent_json_path = parent_dir / FILE_TASK_JSON
    child_json_path = child_dir / FILE_TASK_JSON

    if not parent_json_path.is_file():
        print(colored(f"错误：未找到父任务 task.json：{args.parent_dir}", Colors.RED), file=sys.stderr)
        return 1

    if not child_json_path.is_file():
        print(colored(f"错误：未找到子任务 task.json：{args.child_dir}", Colors.RED), file=sys.stderr)
        return 1

    parent_data = _read_json_file(parent_json_path)
    child_data = _read_json_file(child_json_path)

    if not parent_data or not child_data:
        print(colored("错误：读取 task.json 失败", Colors.RED), file=sys.stderr)
        return 1

    # Remove child from parent's children list
    parent_children = parent_data.get("children", [])
    child_dir_name = child_dir.name
    if child_dir_name in parent_children:
        parent_children.remove(child_dir_name)
        parent_data["children"] = parent_children

    # Clear parent in child's task.json
    child_data["parent"] = None

    # Write both
    _write_json_file(parent_json_path, parent_data)
    _write_json_file(child_json_path, child_data)

    print(colored(f"Unlinked: {child_dir.name} from {parent_dir.name}", Colors.GREEN), file=sys.stderr)
    return 0


# =============================================================================
# Command: list
# =============================================================================

def _get_children_progress(children: list[str], tasks_dir: Path) -> str:
    """Get children progress summary like '[2/3 done]'."""
    if not children:
        return ""
    done_count = 0
    total = len(children)
    for child_name in children:
        child_dir = tasks_dir / child_name
        child_json = child_dir / FILE_TASK_JSON
        if child_json.is_file():
            data = _read_json_file(child_json)
            if data:
                status = data.get("status", "")
                if status in ("completed", "done"):
                    done_count += 1
    return f" [{done_count}/{total} done]"


def cmd_list(args: argparse.Namespace) -> int:
    """List active tasks."""
    repo_root = get_repo_root()
    tasks_dir = get_tasks_dir(repo_root)
    current_task = get_current_task(repo_root)
    developer = get_developer(repo_root)
    filter_mine = args.mine
    filter_status = args.status

    if filter_mine:
        if not developer:
            print(colored("错误：尚未设置开发者。请先运行 init_developer.py", Colors.RED), file=sys.stderr)
            return 1
        print(colored(f"我的任务（负责人：{developer}）：", Colors.BLUE))
    else:
        print(colored("所有活动任务：", Colors.BLUE))
    print()

    # First pass: collect all task data and identify parent/child relationships
    all_tasks: dict[str, dict] = {}
    if tasks_dir.is_dir():
        for d in sorted(tasks_dir.iterdir()):
            if not d.is_dir() or d.name == "archive":
                continue

            dir_name = d.name
            task_json = d / FILE_TASK_JSON
            status = "unknown"
            assignee = "-"
            children: list[str] = []
            parent: str | None = None

            if task_json.is_file():
                data = _read_json_file(task_json)
                if data:
                    status = data.get("status", "unknown")
                    assignee = data.get("assignee", "-")
                    children = data.get("children", [])
                    parent = data.get("parent")

            all_tasks[dir_name] = {
                "status": status,
                "assignee": assignee,
                "children": children,
                "parent": parent,
            }

    # Second pass: display tasks hierarchically
    count = 0

    def _print_task(dir_name: str, indent: int = 0) -> None:
        nonlocal count
        info = all_tasks[dir_name]
        status = info["status"]
        assignee = info["assignee"]
        children = info["children"]

        # Apply --mine filter
        if filter_mine and assignee != developer:
            return

        # Apply --status filter
        if filter_status and status != filter_status:
            return

        relative_path = f"{DIR_WORKFLOW}/{DIR_TASKS}/{dir_name}"
        marker = ""
        if relative_path == current_task:
            marker = f" {colored('<- 当前', Colors.GREEN)}"

        # Children progress
        progress = _get_children_progress(children, tasks_dir) if children else ""

        prefix = "  " * indent + "  - "

        if filter_mine:
            print(f"{prefix}{dir_name}/ ({display_status(status)}){progress}{marker}")
        else:
            print(f"{prefix}{dir_name}/ ({display_status(status)}){progress} [{colored(assignee, Colors.CYAN)}]{marker}")
        count += 1

        # Print children indented
        for child_name in children:
            if child_name in all_tasks:
                _print_task(child_name, indent + 1)

    # Display only top-level tasks (those without a parent)
    for dir_name in sorted(all_tasks.keys()):
        info = all_tasks[dir_name]
        if not info["parent"]:
            _print_task(dir_name)

    if count == 0:
        if filter_mine:
            print("  （暂无分配给你的任务）")
        else:
            print("  （暂无活动任务）")

    print()
    print(f"总计：{count} 个任务")
    return 0


# =============================================================================
# Command: list-archive
# =============================================================================

def cmd_list_archive(args: argparse.Namespace) -> int:
    """List archived tasks."""
    repo_root = get_repo_root()
    tasks_dir = get_tasks_dir(repo_root)
    archive_dir = tasks_dir / "archive"
    month = args.month

    print(colored("已归档任务：", Colors.BLUE))
    print()

    if month:
        month_dir = archive_dir / month
        if month_dir.is_dir():
            print(f"[{month}]")
            for d in sorted(month_dir.iterdir()):
                if d.is_dir():
                    print(f"  - {d.name}/")
        else:
            print(f"  {month} 没有归档任务")
    else:
        if archive_dir.is_dir():
            for month_dir in sorted(archive_dir.iterdir()):
                if month_dir.is_dir():
                    month_name = month_dir.name
                    count = sum(1 for d in month_dir.iterdir() if d.is_dir())
                    print(f"[{month_name}] - {count} task(s)")

    return 0


# =============================================================================
# Command: set-branch
# =============================================================================

def cmd_set_branch(args: argparse.Namespace) -> int:
    """Set git branch for task."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)
    branch = args.branch

    if not branch:
        print(colored("错误：缺少参数", Colors.RED))
        print("用法：python3 task.py set-branch <task-dir> <branch-name>")
        return 1

    task_json = target_dir / FILE_TASK_JSON
    if not task_json.is_file():
        print(colored(f"错误：未在 {target_dir} 找到 task.json", Colors.RED))
        return 1

    data = _read_json_file(task_json)
    if not data:
        return 1

    data["branch"] = branch
    _write_json_file(task_json, data)

    print(colored(f"✓ 已设置分支：{branch}", Colors.GREEN))
    print()
    print(colored("现在你可以启动多 Agent 流水线：", Colors.BLUE))
    print(f"  python3 ./.trellis/scripts/multi_agent/start.py {args.dir}")
    return 0


# =============================================================================
# Command: set-base-branch
# =============================================================================

def cmd_set_base_branch(args: argparse.Namespace) -> int:
    """Set the base branch (PR target) for task."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)
    base_branch = args.base_branch

    if not base_branch:
        print(colored("错误：缺少参数", Colors.RED))
        print("用法：python3 task.py set-base-branch <task-dir> <base-branch>")
        print("示例：python3 task.py set-base-branch <dir> develop")
        print()
        print("这个命令用于设置 PR 的目标分支，也就是功能最终要合入的分支。")
        return 1

    task_json = target_dir / FILE_TASK_JSON
    if not task_json.is_file():
        print(colored(f"错误：未在 {target_dir} 找到 task.json", Colors.RED))
        return 1

    data = _read_json_file(task_json)
    if not data:
        return 1

    data["base_branch"] = base_branch
    _write_json_file(task_json, data)

    print(colored(f"✓ 已设置基础分支：{base_branch}", Colors.GREEN))
    print(f"  PR 目标分支：{base_branch}")
    return 0


# =============================================================================
# Command: set-scope
# =============================================================================

def cmd_set_scope(args: argparse.Namespace) -> int:
    """Set scope for PR title."""
    repo_root = get_repo_root()
    target_dir = _resolve_task_dir(args.dir, repo_root)
    scope = args.scope

    if not scope:
        print(colored("错误：缺少参数", Colors.RED))
        print("用法：python3 task.py set-scope <task-dir> <scope>")
        return 1

    task_json = target_dir / FILE_TASK_JSON
    if not task_json.is_file():
        print(colored(f"错误：未在 {target_dir} 找到 task.json", Colors.RED))
        return 1

    data = _read_json_file(task_json)
    if not data:
        return 1

    data["scope"] = scope
    _write_json_file(task_json, data)

    print(colored(f"✓ 已设置 scope：{scope}", Colors.GREEN))
    return 0


# =============================================================================
# Command: create-pr (delegates to multi-agent script)
# =============================================================================

def cmd_create_pr(args: argparse.Namespace) -> int:
    """Create PR from task - delegates to multi_agent/create_pr.py."""
    import subprocess
    script_dir = Path(__file__).parent
    create_pr_script = script_dir / "multi_agent" / "create_pr.py"

    cmd = [sys.executable, str(create_pr_script)]
    if args.dir:
        cmd.append(args.dir)
    if args.dry_run:
        cmd.append("--dry-run")

    result = subprocess.run(cmd)
    return result.returncode


# =============================================================================
# Help
# =============================================================================

def show_usage() -> None:
    """Show usage help."""
    print("""多 Agent 流水线任务管理脚本

用法：
  python3 task.py create <title>                     创建新任务目录
  python3 task.py create <title> --parent <dir>      创建父任务的子任务
  python3 task.py init-context <dir> <dev_type>      初始化 jsonl 文件
  python3 task.py add-context <dir> <jsonl> <path> [reason]  向 jsonl 添加条目
  python3 task.py validate <dir>                     校验 jsonl 文件
  python3 task.py list-context <dir>                 列出 jsonl 条目
  python3 task.py start <dir>                        设为当前任务
  python3 task.py finish                             清除当前任务
  python3 task.py set-branch <dir> <branch>          为多 Agent 流程设置 git 分支
  python3 task.py set-scope <dir> <scope>            设置 PR 标题的 scope
  python3 task.py create-pr [dir] [--dry-run]        从任务创建 PR
  python3 task.py archive <task-name>                归档已完成任务
  python3 task.py add-subtask <parent> <child>       关联子任务到父任务
  python3 task.py remove-subtask <parent> <child>    取消父子任务关联
  python3 task.py list [--mine] [--status <status>]  列出任务
  python3 task.py list-archive [YYYY-MM]             列出已归档任务

参数：
  dev_type: backend | frontend | fullstack | test | docs

列表选项：
  --mine, -m           仅显示分配给当前开发者的任务
  --status, -s <s>     按状态过滤（planning、in_progress、review、completed）

示例：
  python3 task.py create "新增登录功能" --slug add-login
  python3 task.py create "子任务" --slug child --parent .trellis/tasks/01-21-parent
  python3 task.py init-context .trellis/tasks/01-21-add-login backend
  python3 task.py add-context <dir> implement .trellis/spec/backend/auth.md "认证规范"
  python3 task.py set-branch <dir> task/add-login
  python3 task.py start .trellis/tasks/01-21-add-login
  python3 task.py create-pr                          # 使用当前任务
  python3 task.py create-pr <dir> --dry-run          # 仅预览，不实际修改
  python3 task.py finish
  python3 task.py archive add-login
  python3 task.py add-subtask parent-task child-task  # 关联已有任务
  python3 task.py remove-subtask parent-task child-task
  python3 task.py list                               # 列出全部活动任务
  python3 task.py list --mine                        # 仅列出我的任务
  python3 task.py list --mine --status in_progress   # 列出我进行中的任务
""")


# =============================================================================
# Main Entry
# =============================================================================

def main() -> int:
    """CLI entry point."""
    parser = argparse.ArgumentParser(
        description="多 Agent 流水线任务管理脚本",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    subparsers = parser.add_subparsers(dest="command", help="命令")

    # create
    p_create = subparsers.add_parser("create", help="创建新任务")
    p_create.add_argument("title", help="任务标题")
    p_create.add_argument("--slug", "-s", help="任务 slug")
    p_create.add_argument("--assignee", "-a", help="任务负责人")
    p_create.add_argument("--priority", "-p", default="P2", help="优先级（P0-P3）")
    p_create.add_argument("--description", "-d", help="任务描述")
    p_create.add_argument("--parent", help="父任务目录（用于建立子任务关联）")

    # init-context
    p_init = subparsers.add_parser("init-context", help="初始化上下文文件")
    p_init.add_argument("dir", help="任务目录")
    p_init.add_argument("type", help="开发类型：backend|frontend|fullstack|test|docs")

    # add-context
    p_add = subparsers.add_parser("add-context", help="添加上下文条目")
    p_add.add_argument("dir", help="任务目录")
    p_add.add_argument("file", help="JSONL 文件（implement|check|debug）")
    p_add.add_argument("path", help="要添加的文件路径")
    p_add.add_argument("reason", nargs="?", help="添加原因")

    # validate
    p_validate = subparsers.add_parser("validate", help="校验上下文文件")
    p_validate.add_argument("dir", help="任务目录")

    # list-context
    p_listctx = subparsers.add_parser("list-context", help="列出上下文条目")
    p_listctx.add_argument("dir", help="任务目录")

    # start
    p_start = subparsers.add_parser("start", help="设置当前任务")
    p_start.add_argument("dir", help="任务目录")

    # finish
    subparsers.add_parser("finish", help="清除当前任务")

    # set-branch
    p_branch = subparsers.add_parser("set-branch", help="设置 git 分支")
    p_branch.add_argument("dir", help="任务目录")
    p_branch.add_argument("branch", help="分支名")

    # set-base-branch
    p_base = subparsers.add_parser("set-base-branch", help="设置 PR 目标分支")
    p_base.add_argument("dir", help="任务目录")
    p_base.add_argument("base_branch", help="基础分支名（PR 目标）")

    # set-scope
    p_scope = subparsers.add_parser("set-scope", help="设置 scope")
    p_scope.add_argument("dir", help="任务目录")
    p_scope.add_argument("scope", help="scope 名称")

    # create-pr
    p_pr = subparsers.add_parser("create-pr", help="创建 PR")
    p_pr.add_argument("dir", nargs="?", help="任务目录")
    p_pr.add_argument("--dry-run", action="store_true", help="演练模式")

    # archive
    p_archive = subparsers.add_parser("archive", help="归档任务")
    p_archive.add_argument("name", help="任务名")
    p_archive.add_argument("--no-commit", action="store_true", help="归档后跳过自动 git commit")

    # list
    p_list = subparsers.add_parser("list", help="列出任务")
    p_list.add_argument("--mine", "-m", action="store_true", help="仅显示我的任务")
    p_list.add_argument("--status", "-s", help="按状态过滤")

    # add-subtask
    p_addsub = subparsers.add_parser("add-subtask", help="关联子任务到父任务")
    p_addsub.add_argument("parent_dir", help="父任务目录")
    p_addsub.add_argument("child_dir", help="子任务目录")

    # remove-subtask
    p_rmsub = subparsers.add_parser("remove-subtask", help="取消父子任务关联")
    p_rmsub.add_argument("parent_dir", help="父任务目录")
    p_rmsub.add_argument("child_dir", help="子任务目录")

    # list-archive
    p_listarch = subparsers.add_parser("list-archive", help="列出已归档任务")
    p_listarch.add_argument("month", nargs="?", help="月份（YYYY-MM）")

    args = parser.parse_args()

    if not args.command:
        show_usage()
        return 1

    commands = {
        "create": cmd_create,
        "init-context": cmd_init_context,
        "add-context": cmd_add_context,
        "validate": cmd_validate,
        "list-context": cmd_list_context,
        "start": cmd_start,
        "finish": cmd_finish,
        "set-branch": cmd_set_branch,
        "set-base-branch": cmd_set_base_branch,
        "set-scope": cmd_set_scope,
        "create-pr": cmd_create_pr,
        "archive": cmd_archive,
        "add-subtask": cmd_add_subtask,
        "remove-subtask": cmd_remove_subtask,
        "list": cmd_list,
        "list-archive": cmd_list_archive,
    }

    if args.command in commands:
        return commands[args.command](args)
    else:
        show_usage()
        return 1


if __name__ == "__main__":
    sys.exit(main())
