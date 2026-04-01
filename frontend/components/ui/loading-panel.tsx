import { Panel } from "@/components/ui/panel";

/**
 * LoadingPanelProps 描述加载态面板所需的展示文案。
 */
export interface LoadingPanelProps {
  /** title 表示加载态卡片标题。 */
  title: string;
  /** description 表示该加载态正在准备的内容说明。 */
  description: string;
  /** rows 表示需要渲染的占位条数量。 */
  rows?: number;
}

/**
 * LoadingPanel 负责展示列表、详情或表单等区域的通用加载占位状态。
 */
export function LoadingPanel({
  description,
  rows = 3,
  title,
}: LoadingPanelProps): React.JSX.Element {
  // skeletonRows 用于根据传入数量渲染占位条。
  const skeletonRows = Array.from({ length: rows });

  return (
    <Panel
      title={title}
      description={description}
      eyebrow="加载态"
      className="border-dashed"
    >
      <div className="space-y-3">
        {skeletonRows.map((_, index) => (
          <div
            key={`loading-row-${index}`}
            className="h-12 animate-pulse rounded-2xl border border-border bg-white/70"
          />
        ))}
      </div>
    </Panel>
  );
}
