import { Search } from "lucide-react";
import type { AdminFilterChip } from "@/features/admin/config/dashboard-content";
import { cn } from "@/lib/utils";

/**
 * FilterToolbarProps 描述列表筛选骨架需要的静态筛选数据。
 */
export interface FilterToolbarProps {
  /** chips 表示筛选条中的筛选项集合。 */
  chips: readonly AdminFilterChip[];
  /** searchPlaceholder 表示搜索输入框的提示文案。 */
  searchPlaceholder: string;
}

/**
 * FilterToolbar 负责展示后台列表页常见的搜索框与快捷筛选条。
 */
export function FilterToolbar({
  chips,
  searchPlaceholder,
}: FilterToolbarProps): React.JSX.Element {
  return (
    <div className="rounded-[28px] border border-border bg-white/68 p-4 shadow-[var(--shadow)]">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <label className="flex items-center gap-3 rounded-full border border-border bg-surface-strong px-4 py-3 text-sm text-muted lg:min-w-[22rem]">
          <Search className="size-4 text-accent" />
          <input
            aria-label="示例搜索框"
            className="w-full bg-transparent text-sm text-foreground outline-none placeholder:text-muted"
            placeholder={searchPlaceholder}
            readOnly
          />
        </label>

        <div className="flex flex-wrap gap-2">
          {chips.map((chip) => (
            <span
              key={chip.label}
              className={cn(
                "rounded-full border px-3 py-2 text-xs font-medium",
                chip.active
                  ? "border-accent/20 bg-accent-soft text-accent"
                  : "border-border bg-white/80 text-muted",
              )}
            >
              {chip.label}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
