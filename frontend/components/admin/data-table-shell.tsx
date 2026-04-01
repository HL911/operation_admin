import { Panel } from "@/components/ui/panel";
import { StatusBadge } from "@/components/ui/status-badge";
import type {
  AdminTableColumn,
  AdminTableRow,
} from "@/features/admin/config/dashboard-content";

/**
 * DataTableShellProps 描述通用表格骨架所需的列与行数据。
 */
export interface DataTableShellProps {
  /** columns 表示表头定义。 */
  columns: readonly AdminTableColumn[];
  /** rows 表示表格示例数据。 */
  rows: readonly AdminTableRow[];
}

/**
 * DataTableShell 负责展示后台列表页的典型表格布局与状态列。
 */
export function DataTableShell({
  columns,
  rows,
}: DataTableShellProps): React.JSX.Element {
  return (
    <Panel
      title="列表页骨架"
      description="统一展示搜索、筛选、表格行信息与状态列，适合作为 CRUD 列表页的默认起点。"
      eyebrow="列表骨架"
    >
      <div className="overflow-hidden rounded-[24px] border border-border bg-white/65">
        <div className="grid grid-cols-[1.1fr_1.2fr_1fr_1fr_0.7fr] gap-3 border-b border-border px-4 py-3 text-xs font-semibold uppercase tracking-[0.18em] text-muted">
          {columns.map((column) => (
            <span key={column.key}>{column.label}</span>
          ))}
          <span>状态</span>
        </div>
        <div className="divide-y divide-border">
          {rows.map((row) => (
            <div
              key={row.key}
              className="grid grid-cols-[1.1fr_1.2fr_1fr_1fr_0.7fr] gap-3 px-4 py-4 text-sm text-foreground"
            >
              {row.cells.map((cell, index) => (
                <p
                  key={`${row.key}-${index}`}
                  className="leading-6 text-muted first:font-semibold first:text-foreground"
                >
                  {cell}
                </p>
              ))}
              <div>
                <StatusBadge label={row.statusLabel} tone={row.statusTone} />
              </div>
            </div>
          ))}
        </div>
      </div>
    </Panel>
  );
}
