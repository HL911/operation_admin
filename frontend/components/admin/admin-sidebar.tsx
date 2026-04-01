import type { ComponentType } from "react";
import Link from "next/link";
import { LayoutDashboard, ListFilter, NotebookPen, ScanSearch, Workflow } from "lucide-react";
import { StatusBadge } from "@/components/ui/status-badge";
import type { AdminNavigationItem } from "@/features/admin/config/dashboard-content";
import { cn } from "@/lib/utils";

/**
 * AdminSidebarProps 描述后台侧边栏需要的导航信息。
 */
export interface AdminSidebarProps {
  /** activeHref 表示当前页面的激活导航地址。 */
  activeHref: string;
  /** items 表示要渲染的导航项列表。 */
  items: readonly AdminNavigationItem[];
}

// navigationIconMap 用于把配置中的导航键映射到固定图标。
const navigationIconMap = {
  overview: LayoutDashboard,
  "table-shell": ListFilter,
  "detail-shell": ScanSearch,
  "form-shell": NotebookPen,
  "v0-playbook": Workflow,
} satisfies Record<string, ComponentType<{ className?: string }>>;

/**
 * AdminSidebar 负责渲染后台左侧导航、项目定位与登录挂载说明。
 */
export function AdminSidebar({
  activeHref,
  items,
}: AdminSidebarProps): React.JSX.Element {
  return (
    <aside className="sticky top-0 hidden h-screen w-[var(--sidebar-width)] shrink-0 border-r border-border bg-[rgba(255,250,242,0.88)] px-5 py-6 backdrop-blur-xl xl:flex xl:flex-col">
      <div className="rounded-[26px] border border-border bg-white/75 px-4 py-4 shadow-[var(--shadow)]">
        <p className="text-[11px] font-semibold uppercase tracking-[0.24em] text-muted">
          运营后台前端
        </p>
        <h1 className="mt-3 text-2xl font-semibold tracking-tight text-foreground">
          可复用后台壳层
        </h1>
        <p className="mt-3 text-sm leading-6 text-muted">
          当前阶段优先沉淀结构、状态与代理边界，具体业务页等待 PM 文档落位后接入。
        </p>
      </div>

      <nav className="mt-6 flex-1 space-y-3">
        {items.map((item) => {
          // Icon 表示当前导航项对应的视觉图标组件。
          const Icon = navigationIconMap[item.key] ?? LayoutDashboard;
          // isActive 用于标识当前导航项是否高亮。
          const isActive = item.href === activeHref;

          return (
            <Link
              key={item.key}
              href={item.href}
              className={cn(
                "block rounded-[24px] border px-4 py-4 transition-all duration-200 hover:-translate-y-0.5",
                isActive
                  ? "border-accent/25 bg-accent-soft shadow-[0_18px_40px_rgba(15,118,110,0.18)]"
                  : "border-transparent bg-white/50 hover:border-border hover:bg-white/80",
              )}
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-start gap-3">
                  <span className="mt-0.5 rounded-2xl bg-white/80 p-2 text-accent shadow-sm">
                    <Icon className="size-4" />
                  </span>
                  <div>
                    <p className="text-sm font-semibold text-foreground">{item.label}</p>
                    <p className="mt-1 text-xs leading-5 text-muted">{item.description}</p>
                  </div>
                </div>
                <StatusBadge
                  label={item.status === "ready" ? "已准备" : "规划中"}
                  tone={item.status === "ready" ? "success" : "warning"}
                />
              </div>
            </Link>
          );
        })}
      </nav>

      <div
        id="auth-slot"
        className="rounded-[24px] border border-dashed border-border bg-white/60 px-4 py-4"
      >
        <p className="text-sm font-semibold text-foreground">登录与权限挂载位</p>
        <p className="mt-2 text-xs leading-5 text-muted">
          当前仅保留入口位置。等 PM 或架构方案明确后，再补接真实登录、菜单权限和角色控制。
        </p>
      </div>
    </aside>
  );
}
