import { AdminShell } from "@/components/admin/admin-shell";
import { LoadingPanel } from "@/components/ui/loading-panel";
import {
  adminBreadcrumbItems,
  adminNavigationItems,
  adminQuickActions,
} from "@/features/admin/config/dashboard-content";

/**
 * Loading 负责在首页流式渲染期间展示统一的后台加载骨架。
 */
export default function Loading(): React.JSX.Element {
  return (
    <AdminShell
      activeHref="/"
      breadcrumbItems={adminBreadcrumbItems}
      navigationItems={adminNavigationItems}
      pageTitle="前端接入与 v0 准备工作台"
      pageDescription="正在准备后台骨架预览，请稍候。"
      quickActions={adminQuickActions}
    >
      <div className="grid gap-6 xl:grid-cols-2">
        <LoadingPanel title="列表骨架加载中" description="正在准备筛选条与表格结构。" rows={4} />
        <LoadingPanel title="详情骨架加载中" description="正在准备信息分区与说明文案。" rows={3} />
      </div>
    </AdminShell>
  );
}
