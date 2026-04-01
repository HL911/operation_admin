import { AdminShell } from "@/components/admin/admin-shell";
import { DataTableShell } from "@/components/admin/data-table-shell";
import { DetailPanel } from "@/components/admin/detail-panel";
import { FilterToolbar } from "@/components/admin/filter-toolbar";
import { FormPreviewPanel } from "@/components/admin/form-preview-panel";
import { WorkflowTimeline } from "@/components/admin/workflow-timeline";
import { EmptyState } from "@/components/ui/empty-state";
import { ErrorState } from "@/components/ui/error-state";
import { LoadingPanel } from "@/components/ui/loading-panel";
import { Panel } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import {
  adminBreadcrumbItems,
  adminDetailSections,
  adminFilterChips,
  adminFormSections,
  adminMetricCards,
  adminNavigationItems,
  adminQuickActions,
  adminTableColumns,
  adminTableRows,
  adminWorkflowSteps,
} from "@/features/admin/config/dashboard-content";

/**
 * HomePage 负责展示后台前端骨架的总览页面与可复用模板。
 */
export default function HomePage(): React.JSX.Element {
  return (
    <AdminShell
      activeHref="/"
      breadcrumbItems={adminBreadcrumbItems}
      navigationItems={adminNavigationItems}
      pageTitle="前端接入与 v0 准备工作台"
      pageDescription="这里集中展示运营后台前端的基础布局、代理边界、通用状态组件与 `v0` 落地节奏。等 PM 文档到来后，可以直接从这些骨架出发收敛真实业务页面。"
      quickActions={adminQuickActions}
    >
      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.7fr)_minmax(0,1fr)]">
        <section className="space-y-6">
          <div className="grid gap-4 md:grid-cols-3">
            {adminMetricCards.map((metricCard) => (
              <div
                key={metricCard.title}
                className="rounded-[28px] border border-border bg-surface px-5 py-5 shadow-[var(--shadow)] backdrop-blur-xl"
              >
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-[11px] font-semibold uppercase tracking-[0.24em] text-muted">
                      {metricCard.title}
                    </p>
                    <p className="mt-3 text-3xl font-semibold tracking-tight text-foreground">
                      {metricCard.value}
                    </p>
                  </div>
                  <StatusBadge
                    label={
                      metricCard.tone === "accent"
                        ? "核心"
                        : metricCard.tone === "success"
                          ? "已就绪"
                          : "进行中"
                    }
                    tone={metricCard.tone === "accent" ? "muted" : metricCard.tone}
                  />
                </div>
                <p className="mt-4 text-sm leading-6 text-muted">{metricCard.description}</p>
              </div>
            ))}
          </div>

          <section id="table-shell" className="space-y-6">
            <FilterToolbar
              chips={adminFilterChips}
              searchPlaceholder="示例：按对象名称、外部用户 ID、节点 ID 搜索"
            />
            <DataTableShell columns={adminTableColumns} rows={adminTableRows} />
          </section>

          <div className="grid gap-6 xl:grid-cols-2">
            <section id="detail-shell">
              <DetailPanel sections={adminDetailSections} />
            </section>
            <section id="form-shell">
              <FormPreviewPanel sections={adminFormSections} />
            </section>
          </div>
        </section>

        <section className="space-y-6">
          <Panel
            title="当前准备度"
            description="这部分用来提醒后续功能实现时最不应该偏离的协作约束。"
            eyebrow="执行边界"
          >
            <div className="grid gap-3">
              <div className="rounded-[24px] border border-border bg-white/70 px-4 py-4">
                <p className="text-sm font-semibold text-foreground">页面范围暂不预判</p>
                <p className="mt-2 text-sm leading-6 text-muted">
                  当前只沉淀后台通用结构，不预先把现有 CRUD 文档绑定为首批页面范围。
                </p>
              </div>
              <div className="rounded-[24px] border border-border bg-white/70 px-4 py-4">
                <p className="text-sm font-semibold text-foreground">优先走前端内部代理</p>
                <p className="mt-2 text-sm leading-6 text-muted">
                  后续页面不要直接请求后端地址，统一接入 `/api/admin/*`，这样更利于鉴权和跨域控制。
                </p>
              </div>
              <div className="rounded-[24px] border border-border bg-white/70 px-4 py-4">
                <p className="text-sm font-semibold text-foreground">v0 产物必须二次收敛</p>
                <p className="mt-2 text-sm leading-6 text-muted">
                  初稿回仓库后，需要拆成共享组件、业务组件与接口层，并补齐中文注释与验收态。
                </p>
              </div>
            </div>
          </Panel>

          <section id="v0-playbook">
            <WorkflowTimeline steps={adminWorkflowSteps} />
          </section>

          <LoadingPanel
            title="加载态骨架"
            description="适用于接口返回前的列表、详情或表单预取阶段。"
            rows={3}
          />

          <EmptyState
            title="空态骨架"
            description="适合 PM 文档还未明确字段，或当前筛选条件暂未命中任何数据的场景。"
          />

          <ErrorState
            title="错误态骨架"
            description="适合代理请求失败、鉴权缺失或后端返回非预期状态时统一承载提示。"
          />
        </section>
      </div>
    </AdminShell>
  );
}
