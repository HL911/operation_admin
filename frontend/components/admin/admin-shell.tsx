import type { ReactNode } from "react";
import { AdminSidebar } from "@/components/admin/admin-sidebar";
import { AdminTopbar } from "@/components/admin/admin-topbar";
import type {
  AdminBreadcrumbItem,
  AdminNavigationItem,
  AdminQuickAction,
} from "@/features/admin/config/dashboard-content";

/**
 * AdminShellProps 描述后台页面壳层的输入属性。
 */
export interface AdminShellProps {
  /** activeHref 用于同步当前激活的导航项。 */
  activeHref: string;
  /** breadcrumbItems 表示顶部路径导航。 */
  breadcrumbItems: readonly AdminBreadcrumbItem[];
  /** children 表示页面主体内容。 */
  children: ReactNode;
  /** navigationItems 表示侧边栏导航结构。 */
  navigationItems: readonly AdminNavigationItem[];
  /** pageDescription 用于概括当前页面承担的职责。 */
  pageDescription: string;
  /** pageTitle 表示当前页面主标题。 */
  pageTitle: string;
  /** quickActions 表示顶部快捷操作按钮。 */
  quickActions: readonly AdminQuickAction[];
}

/**
 * AdminShell 负责拼装后台页面的侧栏、顶部栏和主体内容区域。
 */
export function AdminShell({
  activeHref,
  breadcrumbItems,
  children,
  navigationItems,
  pageDescription,
  pageTitle,
  quickActions,
}: AdminShellProps): React.JSX.Element {
  return (
    <div className="min-h-screen xl:flex">
      <AdminSidebar activeHref={activeHref} items={navigationItems} />
      <div className="flex min-h-screen flex-1 flex-col px-4 py-4 sm:px-6 xl:px-8 xl:py-6">
        <AdminTopbar
          breadcrumbItems={breadcrumbItems}
          pageDescription={pageDescription}
          pageTitle={pageTitle}
          quickActions={quickActions}
        />
        <main className="mt-6 flex-1">{children}</main>
      </div>
    </div>
  );
}
