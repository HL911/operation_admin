import type { ButtonHTMLAttributes } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

// buttonVariants 负责统一项目内按钮的尺寸、强调色与边框表现。
export const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-full border text-sm font-medium transition-transform duration-200 hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        primary:
          "border-accent bg-accent px-4 py-2.5 text-accent-foreground shadow-[0_16px_36px_rgba(15,118,110,0.24)] hover:bg-[#0b625c]",
        secondary:
          "border-border bg-white/70 px-4 py-2.5 text-foreground backdrop-blur-sm hover:bg-white",
        ghost:
          "border-transparent bg-transparent px-3 py-2 text-muted hover:bg-white/70 hover:text-foreground",
      },
      size: {
        sm: "h-9 text-xs",
        md: "h-10",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "md",
    },
  },
);

/**
 * ButtonProps 描述基础按钮组件支持的属性。
 */
export interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  /** className 用于追加调用方的局部样式。 */
  className?: string;
}

/**
 * Button 负责渲染项目内统一视觉语言的按钮元素。
 * 输入为原生按钮属性与变体选项；输出为可直接参与表单或工具栏的按钮。
 */
export function Button({
  className,
  size,
  variant,
  ...props
}: ButtonProps): React.JSX.Element {
  return <button className={cn(buttonVariants({ variant, size }), className)} {...props} />;
}
