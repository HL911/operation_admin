import { Panel } from "@/components/ui/panel";
import type { AdminWorkflowStep } from "@/features/admin/config/dashboard-content";

/**
 * WorkflowTimelineProps 描述 `v0` 流程时间线所需的步骤列表。
 */
export interface WorkflowTimelineProps {
  /** steps 表示从需求拆解到验收修正的标准步骤。 */
  steps: readonly AdminWorkflowStep[];
}

/**
 * WorkflowTimeline 负责展示 `v0` 从 PM 文档到仓库落地的标准节奏。
 */
export function WorkflowTimeline({
  steps,
}: WorkflowTimelineProps): React.JSX.Element {
  return (
    <Panel
      title="v0 工作流"
      description="把视觉探索、组件回收和联调边界放到同一条流水线上，避免设计和实现脱节。"
      eyebrow="v0 Playbook"
    >
      <ol className="space-y-4">
        {steps.map((step, index) => (
          <li key={step.title} className="grid grid-cols-[auto_1fr] gap-4">
            <div className="flex flex-col items-center">
              <span className="flex size-10 items-center justify-center rounded-2xl bg-accent text-sm font-semibold text-accent-foreground">
                {index + 1}
              </span>
              {index < steps.length - 1 ? (
                <span className="mt-2 h-full w-px bg-border" aria-hidden="true" />
              ) : null}
            </div>
            <div className="rounded-[24px] border border-border bg-white/70 px-4 py-4">
              <h3 className="text-sm font-semibold text-foreground">{step.title}</h3>
              <p className="mt-2 text-sm leading-6 text-muted">{step.description}</p>
              <p className="mt-3 text-xs leading-5 text-foreground">
                <span className="font-semibold">阶段产出：</span>
                {step.deliverable}
              </p>
            </div>
          </li>
        ))}
      </ol>
    </Panel>
  );
}
