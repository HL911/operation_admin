"use client";

import { useEffect, useId } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * DialogProps 描述通用模态框组件需要的输入属性。
 */
export interface DialogProps {
  /** open 表示当前模态框是否处于打开状态。 */
  open: boolean;
  /** onOpenChange 用于在遮罩、关闭按钮或按下 Esc 时同步开关状态。 */
  onOpenChange: (open: boolean) => void;
  /** title 表示模态框主标题。 */
  title: string;
  /** description 用于展示模态框的补充说明。 */
  description?: string;
  /** children 表示模态框主体区域内容。 */
  children: React.ReactNode;
  /** footer 用于渲染底部操作区。 */
  footer?: React.ReactNode;
  /** className 允许调用方覆盖弹窗面板的局部样式。 */
  className?: string;
}

/**
 * Dialog 负责提供居中显示、遮罩关闭与 Esc 关闭能力的基础模态框外壳。
 */
export function Dialog({
  children,
  className,
  description,
  footer,
  onOpenChange,
  open,
  title,
}: DialogProps): React.JSX.Element | null {
  // titleId 用于把标题和对话框语义正确关联。
  const titleId = useId();
  // descriptionId 用于把描述文本和对话框语义正确关联。
  const descriptionId = useId();

  useEffect(() => {
    if (!open) {
      return;
    }

    // originalOverflow 用于在模态框关闭后恢复页面原始滚动状态。
    const originalOverflow = document.body.style.overflow;

    document.body.style.overflow = "hidden";

    /**
     * handleKeyDown 负责在用户按下 Escape 时关闭模态框。
     */
    function handleKeyDown(event: KeyboardEvent): void {
      if (event.key === "Escape") {
        onOpenChange(false);
      }
    }

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      document.body.style.overflow = originalOverflow;
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [onOpenChange, open]);

  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6">
      <button
        type="button"
        aria-label="关闭模态框"
        className="absolute inset-0 bg-black/72 backdrop-blur-md"
        onClick={() => onOpenChange(false)}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        aria-describedby={description ? descriptionId : undefined}
        className={cn(
          "relative z-10 w-full max-w-2xl overflow-hidden rounded-[28px] border border-white/10 bg-[rgba(10,14,20,0.94)] p-6 shadow-[0_32px_120px_rgba(0,0,0,0.48)] backdrop-blur-2xl",
          className,
        )}
      >
        <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-cyan-300/30 to-transparent" />
        <div className="flex items-start justify-between gap-4">
          <div className="space-y-2">
            <h2 id={titleId} className="text-xl font-semibold tracking-[0.02em] text-zinc-50">
              {title}
            </h2>
            {description ? (
              <p id={descriptionId} className="max-w-xl text-sm leading-6 text-zinc-400">
                {description}
              </p>
            ) : null}
          </div>
          <button
            type="button"
            aria-label="关闭"
            className="inline-flex size-10 items-center justify-center rounded-2xl border border-white/10 bg-white/[0.04] text-zinc-300 transition-colors hover:bg-white/[0.08] hover:text-zinc-50"
            onClick={() => onOpenChange(false)}
          >
            <X className="size-4" />
          </button>
        </div>

        <div className="mt-6">{children}</div>

        {footer ? (
          <div className="mt-6 flex flex-col-reverse gap-3 border-t border-white/10 pt-5 sm:flex-row sm:justify-end">
            {footer}
          </div>
        ) : null}
      </div>
    </div>
  );
}
