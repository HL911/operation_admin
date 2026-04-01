import type { ReactNode } from "react";
import { Panel } from "@/components/ui/panel";

/**
 * ErrorStateProps 描述错误态组件的输入内容。
 */
export interface ErrorStateProps {
  /** title 表示错误态标题。 */
  title: string;
  /** description 用于解释错误现象和推荐排查方向。 */
  description: string;
  /** action 用于放置重试按钮或外链。 */
  action?: ReactNode;
}

/**
 * ErrorState 负责展示代理失败、请求异常或渲染错误时的统一提示区。
 */
export function ErrorState({
  action,
  description,
  title,
}: ErrorStateProps): React.JSX.Element {
  return (
    <Panel title={title} description={description} eyebrow="错误态" className="border-danger/20">
      <div className="rounded-[24px] border border-danger/20 bg-danger-soft px-5 py-6">
        <p className="text-sm leading-6 text-danger">
          优先检查接口路径、环境变量、鉴权头与联调环境是否可达，再决定是否回退到 mock。
        </p>
        {action ? <div className="mt-4">{action}</div> : null}
      </div>
    </Panel>
  );
}
