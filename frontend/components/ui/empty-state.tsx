import type { ReactNode } from "react";
import { Panel } from "@/components/ui/panel";

/**
 * EmptyStateProps 描述空态组件的输入内容。
 */
export interface EmptyStateProps {
  /** title 表示空态主标题。 */
  title: string;
  /** description 解释当前为空的原因或后续建议。 */
  description: string;
  /** action 用于插入按钮或说明链接。 */
  action?: ReactNode;
}

/**
 * EmptyState 负责展示暂无数据、暂无筛选结果或待联调阶段的提示。
 */
export function EmptyState({
  action,
  description,
  title,
}: EmptyStateProps): React.JSX.Element {
  return (
    <Panel title={title} description={description} eyebrow="空态">
      <div className="rounded-[24px] border border-dashed border-border bg-white/55 px-5 py-6">
        <p className="text-sm leading-6 text-muted">
          当 PM 文档尚未明确具体字段时，可以先保留当前空态模板，并接入 mock
          contract 进行布局确认。
        </p>
        {action ? <div className="mt-4">{action}</div> : null}
      </div>
    </Panel>
  );
}
