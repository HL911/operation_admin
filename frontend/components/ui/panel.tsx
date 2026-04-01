import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

/**
 * PanelProps 描述通用信息面板的输入属性。
 */
export interface PanelProps {
  /** title 表示面板主标题。 */
  title: string;
  /** description 用于补充该面板承担的业务说明。 */
  description?: string;
  /** eyebrow 表示标题上方的小型标签文案。 */
  eyebrow?: string;
  /** actions 用于放置右上角的按钮或说明节点。 */
  actions?: ReactNode;
  /** children 表示面板主体承载的自定义内容。 */
  children: ReactNode;
  /** className 允许页面对外层容器追加布局样式。 */
  className?: string;
}

/**
 * Panel 负责提供统一的后台卡片外壳、标题区和内容区。
 */
export function Panel({
  actions,
  children,
  className,
  description,
  eyebrow,
  title,
}: PanelProps): React.JSX.Element {
  return (
    <section
      className={cn(
        "fade-enter rounded-[28px] border border-border bg-surface px-5 py-5 shadow-[var(--shadow)] backdrop-blur-xl sm:px-6",
        className,
      )}
    >
      <header className="mb-5 flex items-start justify-between gap-4">
        <div className="space-y-2">
          {eyebrow ? (
            <p className="text-[11px] font-semibold uppercase tracking-[0.24em] text-muted">
              {eyebrow}
            </p>
          ) : null}
          <div className="space-y-1">
            <h2 className="text-lg font-semibold tracking-tight text-foreground">{title}</h2>
            {description ? <p className="text-sm leading-6 text-muted">{description}</p> : null}
          </div>
        </div>
        {actions ? <div className="shrink-0">{actions}</div> : null}
      </header>
      {children}
    </section>
  );
}
