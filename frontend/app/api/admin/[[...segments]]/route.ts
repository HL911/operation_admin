import { forwardAdminRequest } from "@/features/admin/lib/backend-proxy";

// runtime 固定为 Node.js，以便代理请求在本地开发环境中具备稳定的转发行为。
export const runtime = "nodejs";

// dynamic 强制关闭静态缓存，确保代理始终按请求实时转发。
export const dynamic = "force-dynamic";

/**
 * AdminProxyRouteContext 描述 `/api/admin/*` 路由处理器的动态参数结构。
 */
interface AdminProxyRouteContext {
  /** params 以异步形式提供可选的动态路径片段。 */
  params: Promise<{
    /** segments 表示 `/api/admin` 之后的所有路径片段。 */
    segments?: string[];
  }>;
}

/**
 * proxyRequest 负责提取动态片段并调用统一代理逻辑。
 */
async function proxyRequest(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  // params 表示由 Next.js 提供的动态路由参数。
  const params = await context.params;

  return forwardAdminRequest({
    request,
    segments: params.segments,
  });
}

/**
 * GET 负责转发只读请求到后端管理接口。
 */
export async function GET(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * POST 负责转发创建类请求到后端管理接口。
 */
export async function POST(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * PUT 负责转发完整更新请求到后端管理接口。
 */
export async function PUT(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * PATCH 负责转发局部更新请求到后端管理接口。
 */
export async function PATCH(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * DELETE 负责转发停用或删除语义请求到后端管理接口。
 */
export async function DELETE(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * HEAD 负责转发探活类请求到后端管理接口。
 */
export async function HEAD(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}

/**
 * OPTIONS 负责转发预检或能力探测请求到后端管理接口。
 */
export async function OPTIONS(
  request: Request,
  context: AdminProxyRouteContext,
): Promise<Response> {
  return proxyRequest(request, context);
}
