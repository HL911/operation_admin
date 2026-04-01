"use client";

import { Button } from "@/components/ui/button";
import { ErrorState } from "@/components/ui/error-state";

/**
 * ErrorPageProps 描述根级错误边界传入的异常与恢复方法。
 */
interface ErrorPageProps {
  /** error 表示当前路由渲染过程中抛出的异常。 */
  error: Error & { digest?: string };
  /** reset 用于请求 Next.js 重新尝试渲染当前路由。 */
  reset: () => void;
}

/**
 * GlobalErrorPage 负责在首页渲染失败时输出统一的中文错误态。
 */
export default function GlobalErrorPage({
  error,
  reset,
}: ErrorPageProps): React.JSX.Element {
  return (
    <main className="mx-auto flex min-h-screen max-w-3xl items-center px-4 py-12">
      <ErrorState
        title="页面渲染失败"
        description={`当前错误摘要：${error.message || "未知错误"}。可以先重试一次，如果仍失败，再检查最近改动和控制台日志。`}
        action={
          <Button type="button" onClick={reset}>
            重新加载页面
          </Button>
        }
      />
    </main>
  );
}
