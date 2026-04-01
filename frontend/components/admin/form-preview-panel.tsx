import { Panel } from "@/components/ui/panel";
import type { AdminFormSection } from "@/features/admin/config/dashboard-content";

/**
 * FormPreviewPanelProps 描述表单骨架预览所需的数据分组。
 */
export interface FormPreviewPanelProps {
  /** sections 表示表单中各个字段分组。 */
  sections: readonly AdminFormSection[];
}

/**
 * FormPreviewPanel 负责展示创建与编辑页面推荐采用的字段编排方式。
 */
export function FormPreviewPanel({
  sections,
}: FormPreviewPanelProps): React.JSX.Element {
  return (
    <Panel
      title="表单骨架"
      description="通过字段分组、辅助说明和必填标记，为新增或编辑页面提供统一模板。"
      eyebrow="表单骨架"
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
            <div className="space-y-4">
              {section.fields.map((field) => (
                <label key={`${section.title}-${field.label}`} className="block">
                  <div className="mb-2 flex items-center gap-2 text-sm font-medium text-foreground">
                    <span>{field.label}</span>
                    {field.required ? (
                      <span className="rounded-full bg-danger-soft px-2 py-0.5 text-[10px] uppercase tracking-[0.16em] text-danger">
                        必填
                      </span>
                    ) : null}
                  </div>
                  {field.multiline ? (
                    <textarea
                      className="min-h-28 w-full resize-none rounded-[20px] border border-border bg-surface-strong px-4 py-3 text-sm text-foreground outline-none"
                      placeholder={field.placeholder}
                      readOnly
                    />
                  ) : (
                    <input
                      className="h-12 w-full rounded-full border border-border bg-surface-strong px-4 text-sm text-foreground outline-none"
                      placeholder={field.placeholder}
                      readOnly
                    />
                  )}
                  <p className="mt-2 text-xs leading-5 text-muted">{field.helperText}</p>
                </label>
              ))}
            </div>
          </section>
        ))}
      </div>
    </Panel>
  );
}
