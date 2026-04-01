import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

// statusBadgeVariants 负责统一后台状态徽标的颜色与描边语义。
const statusBadgeVariants = cva(
  "inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold tracking-[0.18em] uppercase",
  {
    variants: {
      tone: {
        success: "border-success/20 bg-success-soft text-success",
        warning: "border-warning/20 bg-warning-soft text-warning",
        danger: "border-danger/20 bg-danger-soft text-danger",
        muted: "border-border bg-white/70 text-muted",
      },
    },
    defaultVariants: {
      tone: "muted",
    },
  },
);

/**
 * StatusBadgeProps 描述状态徽标所需的最小输入。
 */
export interface StatusBadgeProps extends VariantProps<typeof statusBadgeVariants> {
  /** label 表示状态徽标展示的中文文案。 */
  label: string;
  /** className 允许调用方追加局部布局样式。 */
  className?: string;
}

/**
 * StatusBadge 负责以统一视觉样式展示业务状态、联调阶段或模块准备度。
 */
export function StatusBadge({
  className,
  label,
  tone,
}: StatusBadgeProps): React.JSX.Element {
  return <span className={cn(statusBadgeVariants({ tone }), className)}>{label}</span>;
}
