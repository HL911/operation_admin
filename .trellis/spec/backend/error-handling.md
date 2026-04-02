# 错误处理

> 本项目中的错误处理方式。

---

## 概览

当前后端统一通过 `backend/internal/http/response` 包输出 JSON 响应，禁止在 Handler 中直接手写 `map[string]any` 再调用 `ctx.JSON`。

现阶段约定如下：

1. 所有接口响应都必须包含 `code`、`message`、`data`、`timestamp` 四个顶层字段。
2. `timestamp` 统一为 Unix 毫秒时间戳。
3. 成功响应的 `code` 固定为 `200`。
4. 失败响应优先返回 HTTP 状态码；业务自定义错误码使用 `40001`、`40002` 这类数值。
5. 无数据时 `data` 必须显式返回 `null`，不允许省略。
6. 分页列表必须使用统一分页结构：`list`、`total`、`pageNum`、`pageSize`、`pages`。

当前真实实现位置：

- `backend/internal/http/response/json.go`
- `backend/internal/http/handler/health.go`
- `backend/internal/http/router/router.go`

---

## 错误类型

- 通用成功码：`200`
- 参数校验失败：`40001`
- 后台账号已存在：`40010`
- 未登录或登录过期：`401`
- 无权限或缺少二次校验：`403`
- 资源不存在：`404`
- 请求方法不允许：`405`
- 服务器内部错误：`500`

约定：

1. 业务层需要补充更细错误码时，继续沿用数值型业务码，不要改成字符串。
2. 同类错误的 `message` 要稳定，方便前端统一提示和埋点。
3. 对外错误信息要可读，但不能暴露敏感内部实现细节。
4. 后台登录场景中，“账号不存在”和“密码错误”必须统一返回 `401`，避免暴露账号枚举信息。

---

## 错误处理模式

推荐分层方式：

1. `repository` 返回底层错误或包装错误，不负责拼装 HTTP 响应。
2. `service` 负责业务校验、错误分类和业务错误码决策。
3. `handler` 只负责参数绑定、调用 service，并通过 `response` 组件回写统一响应。
4. `middleware` 统一处理 panic、鉴权失败、权限失败等横切错误。

当前约定的响应辅助函数：

- `response.Success`
- `response.SuccessPage`
- `response.ValidationFailed`
- `response.Unauthorized`
- `response.Forbidden`
- `response.NotFound`
- `response.MethodNotAllowed`
- `response.InternalError`

后台鉴权模块当前还约定了以下业务错误语义：

- `401` + `登录名或密码错误`：后台登录名不存在或密码错误
- `403` + `账号已被禁用`：后台账号状态为 disabled
- `40010` + `后台账号已存在`：管理员创建后台账号时 `loginName` 冲突

---

## API 错误响应

### 统一顶层结构

```json
{
  "code": 200,
  "message": "请求成功",
  "data": {},
  "timestamp": 1775104225000
}
```

### 成功响应：对象结果

```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "id": 10001,
    "username": "admin",
    "phone": "13800138000",
    "createTime": "2026-03-31 15:30:25"
  },
  "timestamp": 1775104225000
}
```

### 成功响应：分页列表

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "订单A",
        "status": 1
      },
      {
        "id": 2,
        "name": "订单B",
        "status": 2
      }
    ],
    "total": 128,
    "pageNum": 1,
    "pageSize": 10,
    "pages": 13
  },
  "timestamp": 1775104225000
}
```

### 失败响应：参数校验失败

```json
{
  "code": 40001,
  "message": "参数校验失败：手机号格式错误",
  "data": null,
  "timestamp": 1775104225000
}
```

### 失败响应：登录过期

```json
{
  "code": 401,
  "message": "登录已过期，请重新登录",
  "data": null,
  "timestamp": 1775104225000
}
```

### 失败响应：权限不足或缺少二次校验

```json
{
  "code": 403,
  "message": "当前操作需要二次校验",
  "data": null,
  "timestamp": 1775104225000
}
```

### 失败响应：服务器异常

```json
{
  "code": 500,
  "message": "服务器内部错误，请稍后重试",
  "data": null,
  "timestamp": 1775104225000
}
```

### JSON 请求体常用格式

普通 JSON 提交：

```json
{
  "username": "zhangsan",
  "password": "123456",
  "age": 26,
  "gender": 1
}
```

分页查询请求：

```json
{
  "keyword": "交易记录",
  "status": 0,
  "pageNum": 1,
  "pageSize": 10
}
```

批量操作请求：

```json
{
  "ids": [101, 102, 103],
  "operateType": "delete"
}
```

### 字段规范

- `code`：`200` 表示成功，`4xx` 表示客户端错误，`5xx` 表示服务端错误，自定义业务错误码使用 `40001`、`40002` 等。
- `data`：无数据时统一返回 `null`，不允许省略字段。
- `timestamp`：统一返回 Unix 毫秒时间戳。
- 时间字符串字段：统一使用 `yyyy-MM-dd HH:mm:ss`。
- JSON 字段命名：统一使用小驼峰，例如 `pageNum`、`createTime`、`serviceName`。
- 列表接口：统一返回分页结构，当前 mock 服务默认返回第 1 页数据。

---

## 常见错误

- 直接在 Handler 里调用 `ctx.JSON`，导致不同接口的 `code`、`timestamp`、`data` 结构不一致。
- 成功响应把 `data` 省略掉，前端解析时需要写额外兼容逻辑。
- 分页接口只返回数组，不返回 `total`、`pageNum`、`pageSize`、`pages`。
- 时间字段混用 RFC3339、秒级时间戳和中文时间字符串，导致前端展示和排序行为不一致。
- 错误码使用字符串常量，破坏前后端约定的数值型错误码语义。
- 后台登录接口把“账号不存在”和“密码错误”区分返回，导致前端和攻击者都能枚举后台账号。
