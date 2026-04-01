import { Panel } from "@/components/ui/panel";
import type { AdminDetailSection } from "@/features/admin/config/dashboard-content";

/**
 * DetailPanelProps 描述详情骨架需要的信息分组。
 */
export interface DetailPanelProps {
  /** sections 表示详情面板中的分组列表。 */
  sections: readonly AdminDetailSection[];
}

/**
 * DetailPanel 负责展示详情页常见的区块化信息编排。
 */
export function DetailPanel({
  sections,
}: DetailPanelProps): React.JSX.Element {
  return (
    <Panel
      title="详情页骨架"
      description="通过摘要字段、补充说明和边界提示，快速形成适合运营后台的详情布局。"
      eyebrow="详情骨架"
    >
      <div className="space-y-4">
        {sections.map((section) => (
          <section
            key={section.title}
            className="rounded-[24px] border border-border bg-white/68 px-4 py-4"
          >
            <div className="mb-4">
              <h3 className="text-sm font-semibold text-foreground">{section.title}</h3>
              <p className="mt-1 text-xs leading-5 text-muted">{section.description}</p>
            </div>
            <dl className="grid gap-4 sm:grid-cols-2">
              {section.items.map((item) => (
                <div
                  key={`${section.title}-${item.label}`}
                  className="rounded-[20px] bg-surface-alt/70 px-4 py-3"
                >
                  <dt className="text-xs font-medium uppercase tracking-[0.18em] text-muted">
                    {item.label}
                  </dt>
                  <dd className="mt-2 text-sm font-medium leading-6 text-foreground">
                    {item.value}
                  </dd>
                  {item.hint ? (
                    <p className="mt-1 text-xs leading-5 text-muted">{item.hint}</p>
                  ) : null}
                </div>
              ))}
            </dl>
          </section>
        ))}
      </div>
    </Panel>
  );
}
