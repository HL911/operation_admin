"use client";

import { useEffect, useRef, useState } from "react";
import type { ComponentType, ReactNode } from "react";
import {
  ChevronDown,
  Globe2,
  KeyRound,
  LogOut,
  Settings2,
  Shrimp,
  UserCircle2,
  Users2,
} from "lucide-react";
import {
  followerUsersConsoleLanguageOptions,
  followerUsersConsoleNavItems,
  followerUsersConsolePageTitle,
  followerUsersConsoleProductName,
  followerUsersConsoleUserActions,
  type FollowerUsersConsoleLanguageOption,
  type FollowerUsersConsoleNavKey,
  type FollowerUsersConsoleUserAction,
} from "@/features/follower-users/config/follower-user-content";
import { cn } from "@/lib/utils";

/**
 * FollowerUsersConsoleShellProps 描述运营后台控制台壳层所需的输入属性。
 */
export interface FollowerUsersConsoleShellProps {
  /** children 表示右侧主内容区需要渲染的页面主体。 */
  children: ReactNode;
}

// navIconMap 用于把侧边栏菜单 key 映射到对应的细线图标。
const navIconMap: Record<FollowerUsersConsoleNavKey, ComponentType<{ className?: string }>> = {
  "little-follower-query": Shrimp,
  "big-follower-query": Users2,
  "big-follower-key-query": KeyRound,
};

// userActionIconMap 用于把头像菜单动作映射到对应图标。
const userActionIconMap: Record<
  FollowerUsersConsoleUserAction["key"],
  ComponentType<{ className?: string }>
> = {
  profile: UserCircle2,
  settings: Settings2,
  logout: LogOut,
};

/**
 * renderLanguageOptionClassName 负责根据语言是否激活返回对应的样式类名。
 */
function renderLanguageOptionClassName(isActive: boolean): string {
  return isActive
    ? "border-cyan-300/20 bg-cyan-300/12 text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.12)]"
    : "border-white/10 bg-white/[0.03] text-zinc-400 hover:bg-white/[0.06] hover:text-zinc-200";
}

/**
 * renderUserActionClassName 负责根据动作语义返回下拉菜单项样式。
 */
function renderUserActionClassName(
  tone: FollowerUsersConsoleUserAction["tone"],
): string {
  return tone === "danger"
    ? "text-red-100 hover:bg-red-400/[0.08] hover:text-red-50"
    : "text-zinc-300 hover:bg-white/[0.05] hover:text-zinc-50";
}

/**
 * FollowerUsersConsoleShell 负责渲染固定侧边栏、固定顶栏与主内容区的后台壳层。
 */
export function FollowerUsersConsoleShell({
  children,
}: FollowerUsersConsoleShellProps): React.JSX.Element {
  // currentLanguage 保存当前顶部语言切换器中的选中语言。
  const [currentLanguage, setCurrentLanguage] = useState<FollowerUsersConsoleLanguageOption["key"]>("zh");
  // isUserMenuOpen 表示当前头像下拉菜单是否处于展开状态。
  const [isUserMenuOpen, setIsUserMenuOpen] = useState<boolean>(false);
  // userMenuRef 用于判断点击是否发生在用户菜单区域外部。
  const userMenuRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    /**
     * handlePointerDown 负责在点击菜单区域外部时关闭头像下拉菜单。
     */
    function handlePointerDown(event: MouseEvent): void {
      if (!userMenuRef.current) {
        return;
      }

      if (!userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false);
      }
    }

    window.addEventListener("mousedown", handlePointerDown);

    return () => {
      window.removeEventListener("mousedown", handlePointerDown);
    };
  }, []);

  return (
    <div className="min-h-screen bg-[#04070a] text-zinc-50">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-[16.5rem] border-r border-white/8 bg-[rgba(7,10,14,0.82)] px-4 py-5 backdrop-blur-xl lg:block">
        <div className="rounded-[24px] border border-white/10 bg-white/[0.03] px-4 py-4">
          <p className="text-[11px] font-medium uppercase tracking-[0.24em] text-zinc-500">
            {followerUsersConsoleProductName}
          </p>
          <p className="mt-3 text-lg font-semibold tracking-[0.02em] text-zinc-50">
            查询控制台
          </p>
        </div>

        <nav className="mt-6 space-y-2">
          {followerUsersConsoleNavItems.map((item) => {
            // Icon 表示当前侧边栏菜单对应的图标组件。
            const Icon = navIconMap[item.key];

            return (
              <button
                key={item.key}
                type="button"
                aria-current={item.active ? "page" : undefined}
                className={cn(
                  "group flex w-full items-center gap-3 rounded-[20px] border px-4 py-3 text-left transition-all duration-200",
                  item.active
                    ? "border-cyan-300/18 bg-cyan-300/[0.10] text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.10)]"
                    : "border-transparent bg-transparent text-zinc-400 hover:border-white/10 hover:bg-white/[0.04] hover:text-zinc-100",
                )}
              >
                <span
                  className={cn(
                    "inline-flex size-9 items-center justify-center rounded-2xl border transition-colors",
                    item.active
                      ? "border-cyan-300/15 bg-cyan-300/12 text-cyan-100"
                      : "border-white/10 bg-white/[0.04] text-zinc-400 group-hover:text-zinc-100",
                  )}
                >
                  <Icon className="size-4" />
                </span>
                <span className="text-sm font-medium">{item.label}</span>
              </button>
            );
          })}
        </nav>
      </aside>

      <header className="fixed left-0 right-0 top-0 z-20 border-b border-white/8 bg-[rgba(4,7,10,0.82)] backdrop-blur-xl lg:left-[16.5rem]">
        <div className="flex h-[4.5rem] items-center justify-between px-4 sm:px-6">
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-zinc-100">
              {followerUsersConsoleProductName}
            </p>
            <p className="mt-1 truncate text-xs uppercase tracking-[0.22em] text-zinc-500">
              {followerUsersConsolePageTitle}
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="hidden items-center gap-1 rounded-full border border-white/10 bg-white/[0.03] p-1 sm:flex">
              <span className="inline-flex size-8 items-center justify-center rounded-full text-zinc-500">
                <Globe2 className="size-4" />
              </span>
              {followerUsersConsoleLanguageOptions.map((option) => (
                <button
                  key={option.key}
                  type="button"
                  className={cn(
                    "rounded-full border px-3 py-1.5 text-xs font-medium transition-colors",
                    renderLanguageOptionClassName(currentLanguage === option.key),
                  )}
                  onClick={() => setCurrentLanguage(option.key)}
                >
                  {option.label}
                </button>
              ))}
            </div>

            <div ref={userMenuRef} className="relative">
              <button
                type="button"
                aria-expanded={isUserMenuOpen}
                aria-haspopup="menu"
                className="flex items-center gap-3 rounded-full border border-white/10 bg-white/[0.03] px-2.5 py-2 text-left transition-colors hover:bg-white/[0.06]"
                onClick={() => setIsUserMenuOpen((currentOpen) => !currentOpen)}
              >
                <span className="inline-flex size-9 items-center justify-center rounded-full border border-cyan-300/15 bg-cyan-300/12 text-sm font-semibold text-cyan-50">
                  OP
                </span>
                <span className="hidden sm:block">
                  <span className="block text-sm font-medium text-zinc-100">运营账号</span>
                  <span className="block text-xs text-zinc-500">
                    {currentLanguage === "zh" ? "已登录" : "Signed In"}
                  </span>
                </span>
                <ChevronDown className="size-4 text-zinc-500" />
              </button>

              {isUserMenuOpen ? (
                <div className="absolute right-0 top-[calc(100%+0.75rem)] w-52 rounded-[22px] border border-white/10 bg-[rgba(10,14,19,0.96)] p-2 shadow-[0_24px_80px_rgba(0,0,0,0.42)] backdrop-blur-2xl">
                  {followerUsersConsoleUserActions.map((action) => {
                    // Icon 表示当前头像菜单动作对应的图标。
                    const Icon = userActionIconMap[action.key];

                    return (
                      <button
                        key={action.key}
                        type="button"
                        className={cn(
                          "flex w-full items-center gap-3 rounded-2xl px-3 py-2.5 text-sm transition-colors",
                          renderUserActionClassName(action.tone),
                        )}
                        onClick={() => setIsUserMenuOpen(false)}
                      >
                        <Icon className="size-4" />
                        <span>{action.label}</span>
                      </button>
                    );
                  })}
                </div>
              ) : null}
            </div>
          </div>
        </div>
      </header>

      <main className="px-4 pb-6 pt-[5.5rem] sm:px-6 lg:pl-[18rem] lg:pr-6">
        <div className="mx-auto max-w-[1360px]">
          <nav className="mb-5 flex gap-2 overflow-x-auto pb-1 lg:hidden">
            {followerUsersConsoleNavItems.map((item) => {
              // Icon 表示当前移动菜单项对应的图标组件。
              const Icon = navIconMap[item.key];

              return (
                <button
                  key={item.key}
                  type="button"
                  aria-current={item.active ? "page" : undefined}
                  className={cn(
                    "inline-flex shrink-0 items-center gap-2 rounded-full border px-4 py-2 text-sm transition-colors",
                    item.active
                      ? "border-cyan-300/18 bg-cyan-300/[0.10] text-cyan-50"
                      : "border-white/10 bg-white/[0.03] text-zinc-400",
                  )}
                >
                  <Icon className="size-4" />
                  <span>{item.label}</span>
                </button>
              );
            })}
          </nav>

          {children}
        </div>
      </main>
    </div>
  );
}
