import Link from "next/link";
import type { AdminBreadcrumbItem } from "@/features/admin/config/dashboard-content";

/**
 * AdminBreadcrumbsProps 描述面包屑组件所需的层级数据。
 */
export interface AdminBreadcrumbsProps {
  /** items 表示当前页面的层级路径。 */
  items: readonly AdminBreadcrumbItem[];
}

/**
 * AdminBreadcrumbs 负责渲染后台顶部的层级导航路径。
 */
export function AdminBreadcrumbs({
  items,
}: AdminBreadcrumbsProps): React.JSX.Element {
  return (
    <nav aria-label="页面路径" className="flex flex-wrap items-center gap-2 text-sm text-muted">
      {items.map((item, index) => (
        <span key={`${item.label}-${index}`} className="flex items-center gap-2">
          {item.href ? (
            <Link className="transition-colors hover:text-foreground" href={item.href}>
              {item.label}
            </Link>
          ) : (
            <span className="text-foreground">{item.label}</span>
          )}
          {index < items.length - 1 ? <span aria-hidden="true">/</span> : null}
        </span>
      ))}
    </nav>
  );
}
