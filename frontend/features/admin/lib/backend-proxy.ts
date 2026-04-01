import { NextResponse } from "next/server";
import { getRequiredServerEnv } from "@/lib/server-env";

/**
 * ForwardAdminRequestInput 描述代理调用所需的最小输入。
 */
export interface ForwardAdminRequestInput {
  /** request 表示前端 Route Handler 当前收到的原始请求。 */
  request: Request;
  /** segments 表示 `/api/admin/*` 中被捕获的动态路径片段。 */
  segments?: readonly string[];
}

// BACKEND_BASE_URL_ENV_KEY 表示后端代理基址对应的环境变量名。
const BACKEND_BASE_URL_ENV_KEY = "BACKEND_BASE_URL";

// BODY_METHODS 描述需要读取并透传请求体的 HTTP 方法集合。
const BODY_METHODS = new Set(["POST", "PUT", "PATCH", "DELETE"]);

// HOP_BY_HOP_HEADERS 记录不应透传到后端的逐跳请求头。
const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "content-length",
  "host",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
]);

/**
 * resolveBackendBaseUrl 负责解析并校验后端服务基址。
 * 输出始终为带协议的绝对 URL；若配置非法则抛出异常，由上层转换为 JSON 错误响应。
 */
export function resolveBackendBaseUrl(): string {
  // configuredBaseUrl 表示环境变量中声明的后端地址。
  const configuredBaseUrl = getRequiredServerEnv(BACKEND_BASE_URL_ENV_KEY, {
    example: "http://127.0.0.1:8081",
    location: "frontend/.env.local",
  });

  return new URL(configuredBaseUrl).toString();
}

/**
 * buildAdminTargetUrl 负责把前端代理路径转换成真实后端地址。
 * 输入为请求 URL 与动态片段；输出为拼接好查询参数的后端目标地址。
 */
export function buildAdminTargetUrl(
  requestUrl: URL,
  segments?: readonly string[],
): URL {
  // normalizedSegments 用于去除空片段并保持路径拼接稳定。
  const normalizedSegments = (segments ?? []).filter(Boolean);
  // encodedPathSegments 用于保证动态片段在转发时不会破坏路径结构。
  const encodedPathSegments = normalizedSegments.map((segment) => encodeURIComponent(segment));
  // normalizedPath 用于生成 `/admin/v1/*` 的固定代理目标。
  const normalizedPath =
    encodedPathSegments.length > 0
      ? `/admin/v1/${encodedPathSegments.join("/")}`
      : "/admin/v1";
  // targetUrl 指向最终后端资源地址，并保留原始查询串。
  const targetUrl = new URL(normalizedPath, resolveBackendBaseUrl());

  targetUrl.search = requestUrl.search;

  return targetUrl;
}

/**
 * buildForwardHeaders 负责筛掉逐跳请求头，并补充常用的转发上下文。
 * 输入为浏览器请求头与当前请求地址；输出为可安全透传给后端的请求头集合。
 */
export function buildForwardHeaders(sourceHeaders: Headers, requestUrl: URL): Headers {
  // forwardedHeaders 用于承载准备发往后端的新请求头。
  const forwardedHeaders = new Headers();

  sourceHeaders.forEach((value, key) => {
    // normalizedKey 用于统一处理不同大小写形式的请求头名。
    const normalizedKey = key.toLowerCase();

    if (HOP_BY_HOP_HEADERS.has(normalizedKey)) {
      return;
    }

    forwardedHeaders.set(key, value);
  });

  forwardedHeaders.set("x-forwarded-host", requestUrl.host);
  forwardedHeaders.set("x-forwarded-path", requestUrl.pathname);
  forwardedHeaders.set("x-forwarded-proto", requestUrl.protocol.replace(":", ""));

  return forwardedHeaders;
}

/**
 * forwardAdminRequest 负责执行真正的代理调用，并把结果包装成前端可直接返回的响应。
 * 在后端不可达或配置错误时，会统一返回带中文说明的 JSON 错误响应。
 */
export async function forwardAdminRequest(
  input: ForwardAdminRequestInput,
): Promise<Response> {
  try {
    // requestUrl 表示当前代理请求的完整前端地址。
    const requestUrl = new URL(input.request.url);
    // targetUrl 表示实际要访问的后端地址。
    const targetUrl = buildAdminTargetUrl(requestUrl, input.segments);
    // forwardedHeaders 保存已经过滤过的透传请求头。
    const forwardedHeaders = buildForwardHeaders(input.request.headers, requestUrl);
    // requestBody 在需要时读取原始请求体，避免浏览器直连后端。
    const requestBody = BODY_METHODS.has(input.request.method)
      ? await input.request.arrayBuffer()
      : undefined;
    // upstreamResponse 表示后端服务返回的原始结果。
    const upstreamResponse = await fetch(targetUrl, {
      method: input.request.method,
      headers: forwardedHeaders,
      body: requestBody,
      cache: "no-store",
      redirect: "manual",
    });

    return new Response(upstreamResponse.body, {
      status: upstreamResponse.status,
      statusText: upstreamResponse.statusText,
      headers: normalizeResponseHeaders(upstreamResponse.headers),
    });
  } catch (error) {
    return createProxyErrorResponse(
      "后端服务暂时不可达，请检查 BACKEND_BASE_URL 或后端进程状态。",
      502,
      toErrorDetail(error),
    );
  }
}

/**
 * normalizeResponseHeaders 负责清理后端响应中的逐跳头，确保浏览器能够稳定接收。
 */
function normalizeResponseHeaders(sourceHeaders: Headers): Headers {
  // normalizedHeaders 用于承载准备返回给浏览器的响应头。
  const normalizedHeaders = new Headers();

  sourceHeaders.forEach((value, key) => {
    // normalizedKey 用于统一过滤逐跳响应头。
    const normalizedKey = key.toLowerCase();

    if (HOP_BY_HOP_HEADERS.has(normalizedKey)) {
      return;
    }

    normalizedHeaders.set(key, value);
  });

  return normalizedHeaders;
}

/**
 * createProxyErrorResponse 负责生成统一格式的代理错误 JSON 响应。
 * 输出会包含中文错误说明、细节摘要和默认建议，方便前端联调排查。
 */
function createProxyErrorResponse(
  message: string,
  status: number,
  detail: string,
): Response {
  return NextResponse.json(
    {
      code: "ADMIN_PROXY_UNAVAILABLE",
      message,
      detail,
      suggestion: "确认后端是否监听在配置地址，并检查 `/admin/v1/*` 路径是否可访问。",
    },
    {
      status,
    },
  );
}

/**
 * toErrorDetail 负责把未知异常转换成可展示的中文调试细节。
 */
function toErrorDetail(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return "发生了未知代理错误。";
}
