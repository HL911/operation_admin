import Link from "next/link";
import { Activity, Blocks, Sparkles } from "lucide-react";
import { AdminBreadcrumbs } from "@/components/admin/admin-breadcrumbs";
import { buttonVariants } from "@/components/ui/button";
import { StatusBadge } from "@/components/ui/status-badge";
import type {
  AdminBreadcrumbItem,
  AdminQuickAction,
} from "@/features/admin/config/dashboard-content";
import { cn } from "@/lib/utils";

/**
 * AdminTopbarProps 描述顶部栏需要的标题、路径与快捷操作信息。
 */
export interface AdminTopbarProps {
  /** breadcrumbItems 表示当前页面的面包屑路径。 */
  breadcrumbItems: readonly AdminBreadcrumbItem[];
  /** pageTitle 表示页面主标题。 */
  pageTitle: string;
  /** pageDescription 表示页面摘要说明。 */
  pageDescription: string;
  /** quickActions 表示顶部快捷操作按钮。 */
  quickActions: readonly AdminQuickAction[];
}

/**
 * AdminTopbar 负责渲染后台主区域顶部的摘要、状态与快捷操作。
 */
export function AdminTopbar({
  breadcrumbItems,
  pageDescription,
  pageTitle,
  quickActions,
}: AdminTopbarProps): React.JSX.Element {
  return (
    <header className="rounded-[32px] border border-border bg-surface px-5 py-5 shadow-[var(--shadow)] backdrop-blur-xl sm:px-6">
      <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-4">
          <AdminBreadcrumbs items={breadcrumbItems} />
          <div className="space-y-3">
            <div className="flex flex-wrap items-center gap-3">
              <h2 className="text-3xl font-semibold tracking-tight text-foreground">
                {pageTitle}
              </h2>
              <StatusBadge label="骨架阶段" tone="success" />
            </div>
            <p className="max-w-3xl text-sm leading-7 text-muted">{pageDescription}</p>
          </div>
        </div>

        <div className="grid gap-3 sm:min-w-[20rem]">
          <div className="grid gap-3 rounded-[24px] border border-border bg-white/70 p-4">
            <div className="flex items-center gap-3">
              <span className="rounded-2xl bg-accent-soft p-2 text-accent">
                <Activity className="size-4" />
              </span>
              <div>
                <p className="text-sm font-semibold text-foreground">代理边界已固定</p>
                <p className="text-xs leading-5 text-muted">前端统一走 `/api/admin/*`。</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <span className="rounded-2xl bg-warning-soft p-2 text-warning">
                <Blocks className="size-4" />
              </span>
              <div>
                <p className="text-sm font-semibold text-foreground">业务页待 PM 文档</p>
                <p className="text-xs leading-5 text-muted">当前不预判具体页面范围。</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <span className="rounded-2xl bg-success-soft p-2 text-success">
                <Sparkles className="size-4" />
              </span>
              <div>
                <p className="text-sm font-semibold text-foreground">v0 节奏已收口</p>
                <p className="text-xs leading-5 text-muted">生成初稿后统一回仓库拆分组件。</p>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            {quickActions.map((action) => (
              <Link
                key={action.label}
                href={action.href}
                className={cn(
                  buttonVariants({
                    variant: action.tone === "primary" ? "primary" : "secondary",
                    size: "sm",
                  }),
                  "w-fit",
                )}
              >
                {action.label}
              </Link>
            ))}
          </div>
        </div>
      </div>
    </header>
  );
}
