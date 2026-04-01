import type { Metadata } from "next";
import "./globals.css";

// metadata 负责定义整个后台前端的默认标题与说明。
export const metadata: Metadata = {
  title: {
    default: "运营后台前端工作台",
    template: "%s | 运营后台前端工作台",
  },
  description: "用于承接运营后台需求、v0 设计稿和后端代理联调的前端基础工程。",
};

/**
 * RootLayout 负责提供 App Router 所需的根文档结构与全局样式挂载点。
 * 输入为页面子树；输出为设置中文语言环境后的 HTML 外壳。
 */
export default function RootLayout({
  children,
}: Readonly<{
  /** children 表示当前路由要渲染的页面内容。 */
  children: React.ReactNode;
}>): React.JSX.Element {
  return (
    <html lang="zh-CN" className="h-full">
      <body className="admin-ambient min-h-full bg-background text-foreground antialiased">
        {children}
      </body>
    </html>
  );
}
