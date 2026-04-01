import type {
  FollowerUserApiEnvelope,
  FollowerUserCreatePayload,
  FollowerUserDeletePayload,
  FollowerUserDeleteResult,
  FollowerUserListData,
  FollowerUserListQuery,
  FollowerUserRecord,
  FollowerUserUpdatePayload,
} from "@/features/follower-users/types";

// FOLLOWER_USERS_API_PREFIX 统一声明小龙虾用户通过前端代理访问的路径前缀。
const FOLLOWER_USERS_API_PREFIX = "/api/admin/follower-users";

/**
 * fetchFollowerUserList 负责按接口文档定义的筛选条件查询列表。
 */
export async function fetchFollowerUserList(
  query: FollowerUserListQuery,
): Promise<FollowerUserListData> {
  // searchParams 用于按需拼接非空筛选条件。
  const searchParams = new URLSearchParams();

  appendSearchParam(searchParams, "pageNum", String(query.pageNum));
  appendSearchParam(searchParams, "pageSize", String(query.pageSize));
  appendSearchParam(searchParams, "userId", query.userId);
  appendSearchParam(searchParams, "accountStatus", query.accountStatus);
  appendSearchParam(searchParams, "strategyStatus", query.strategyStatus);
  appendSearchParam(searchParams, "bindingStatus", query.bindingStatus);
  appendSearchParam(searchParams, "responsibilityDomain", query.responsibilityDomain);
  appendSearchParam(searchParams, "updatedFrom", toRequestDateTime(query.updatedFrom));
  appendSearchParam(searchParams, "updatedTo", toRequestDateTime(query.updatedTo));

  return requestFollowerUserApi<FollowerUserListData>(
    `${FOLLOWER_USERS_API_PREFIX}?${searchParams.toString()}`,
    {
      method: "GET",
      errorMessage: "查询小龙虾用户列表失败。",
    },
  );
}

/**
 * fetchFollowerUserDetail 负责根据用户 ID 查询详情。
 */
export async function fetchFollowerUserDetail(userId: string): Promise<FollowerUserRecord> {
  return requestFollowerUserApi<FollowerUserRecord>(
    `${FOLLOWER_USERS_API_PREFIX}/${encodeURIComponent(userId)}`,
    {
      method: "GET",
      errorMessage: "查询小龙虾用户详情失败。",
    },
  );
}

/**
 * createFollowerUser 负责提交新建请求。
 */
export async function createFollowerUser(
  payload: FollowerUserCreatePayload,
): Promise<FollowerUserRecord> {
  return requestFollowerUserApi<FollowerUserRecord>(FOLLOWER_USERS_API_PREFIX, {
    method: "POST",
    body: JSON.stringify(payload),
    errorMessage: "创建小龙虾用户失败。",
  });
}

/**
 * updateFollowerUser 负责提交编辑请求。
 */
export async function updateFollowerUser(
  userId: string,
  payload: FollowerUserUpdatePayload,
): Promise<FollowerUserRecord> {
  return requestFollowerUserApi<FollowerUserRecord>(
    `${FOLLOWER_USERS_API_PREFIX}/${encodeURIComponent(userId)}`,
    {
      method: "PATCH",
      body: JSON.stringify(payload),
      errorMessage: "更新小龙虾用户失败。",
    },
  );
}

/**
 * deleteFollowerUser 负责提交逻辑删除请求。
 */
export async function deleteFollowerUser(
  userId: string,
  payload: FollowerUserDeletePayload,
): Promise<FollowerUserDeleteResult> {
  return requestFollowerUserApi<FollowerUserDeleteResult>(
    `${FOLLOWER_USERS_API_PREFIX}/${encodeURIComponent(userId)}`,
    {
      method: "DELETE",
      body: JSON.stringify(payload),
      errorMessage: "删除小龙虾用户失败。",
    },
  );
}

/**
 * requestFollowerUserApi 负责统一处理成功响应、错误响应和 JSON 解析。
 */
async function requestFollowerUserApi<TData>(
  input: string,
  options: RequestInit & { errorMessage: string },
): Promise<TData> {
  // requestHeaders 用于为 JSON 请求补齐统一请求头。
  const requestHeaders = new Headers(options.headers);

  if (options.body && !requestHeaders.has("Content-Type")) {
    requestHeaders.set("Content-Type", "application/json");
  }

  // response 表示前端代理层返回的原始 HTTP 响应。
  const response = await fetch(input, {
    ...options,
    headers: requestHeaders,
    cache: "no-store",
  });
  // responseBody 表示被解析后的接口返回体。
  const responseBody = await parseResponseBody(response);

  if (!response.ok) {
    throw new Error(extractErrorMessage(responseBody) || options.errorMessage);
  }

  if (!isApiEnvelope<TData>(responseBody)) {
    throw new Error("接口返回格式不符合预期。");
  }

  return responseBody.data;
}

/**
 * parseResponseBody 负责根据响应内容尝试解析 JSON 或读取文本。
 */
async function parseResponseBody(response: Response): Promise<unknown> {
  // rawText 用于保留响应原文，兼容 JSON 和文本两类场景。
  const rawText = await response.text();

  if (!rawText) {
    return null;
  }

  try {
    return JSON.parse(rawText) as unknown;
  } catch {
    return rawText;
  }
}

/**
 * extractErrorMessage 负责尽量从接口响应体中提取可读的错误提示。
 */
function extractErrorMessage(body: unknown): string | null {
  if (typeof body === "string" && body.trim()) {
    return body;
  }

  if (body && typeof body === "object") {
    // messageValue 用于读取对象中的 `message` 字段。
    const messageValue = Reflect.get(body, "message");

    if (typeof messageValue === "string" && messageValue.trim()) {
      return messageValue;
    }
  }

  return null;
}

/**
 * isApiEnvelope 负责校验响应体是否符合通用包裹结构。
 */
function isApiEnvelope<TData>(body: unknown): body is FollowerUserApiEnvelope<TData> {
  return body !== null && typeof body === "object" && "data" in body;
}

/**
 * appendSearchParam 负责仅在值非空时向查询参数中追加字段。
 */
function appendSearchParam(
  searchParams: URLSearchParams,
  key: string,
  value: string,
): void {
  if (!value.trim()) {
    return;
  }

  searchParams.set(key, value);
}

/**
 * toRequestDateTime 负责把 `datetime-local` 值转换成后端更容易识别的时间格式。
 */
function toRequestDateTime(value: string): string {
  if (!value.trim()) {
    return "";
  }

  // normalizedValue 用于把浏览器中的 `T` 分隔符替换为空格。
  const normalizedValue = value.replace("T", " ");

  return normalizedValue.length === 16 ? `${normalizedValue}:00` : normalizedValue;
}
