import type { FollowerUserTone } from "@/features/follower-users/config/follower-user-content";
import { cn } from "@/lib/utils";

/**
 * FollowerUserStatusPillProps 描述深色工作台中的状态标签输入属性。
 */
export interface FollowerUserStatusPillProps {
  /** label 表示标签展示的中文文本。 */
  label: string;
  /** tone 表示标签采用的颜色语义。 */
  tone: FollowerUserTone;
  /** className 允许调用方追加局部样式。 */
  className?: string;
}

// statusToneClassNameMap 用于统一管理深色主题下的状态标签颜色。
const statusToneClassNameMap: Record<FollowerUserTone, string> = {
  accent: "border-cyan-300/15 bg-cyan-300/10 text-cyan-100",
  success: "border-emerald-300/15 bg-emerald-300/10 text-emerald-100",
  warning: "border-amber-300/18 bg-amber-300/10 text-amber-100",
  muted: "border-white/10 bg-white/[0.04] text-zinc-300",
};

/**
 * FollowerUserStatusPill 负责在深色 Web3 风格页面中展示统一的状态标签。
 */
export function FollowerUserStatusPill({
  className,
  label,
  tone,
}: FollowerUserStatusPillProps): React.JSX.Element {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full border px-2.5 py-1 text-[11px] font-medium tracking-[0.14em]",
        statusToneClassNameMap[tone],
        className,
      )}
    >
      {label}
    </span>
  );
}
