# PincerMarket Schema 数据字典

> 生成日期：2026-04-01
> 来源文件：`/Users/leon/Desktop/my/operation_admin/schema_pincermarket_20260401.sql`
> 说明：SQL dump 中未发现 `COMMENT ON TABLE` / `COMMENT ON COLUMN`。本文件中的“中文名”“字段意思”基于模式名、表名、字段名自动推断，适合用于阅读和建模参考，业务含义仍建议结合代码再校对。

## 总览

- 模式数：11
- 表数：127

## 模式目录

- `connector_hub`：连接器中心（2 张表）
- `core_data`：核心数据（2 张表）
- `decision_producer`：决策产出（4 张表）
- `log_kms`：日志与密钥管理（5 张表）
- `notification_service`：通知服务（6 张表）
- `operator_portal`：运营门户（7 张表）
- `order_gateway`：订单网关（4 张表）
- `settlement_risk`：结算风控（7 张表）
- `signal_server`：信号服务（3 张表）
- `usstocks_runtime`：美股运行时（7 张表）
- `public`：公共模式（80 张表）

## 模式：`connector_hub`

- 中文名：连接器中心
- 表数量：2

### 表：`connector_hub.orders`

- 表中文名：订单
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `exchange_order_id` | `text` | 交易所订单ID | 交易所在外部系统中的订单标识。 |
| `intent_id` | `text` | 意图ID | 关联的意图唯一标识。 |
| `exchange_name` | `text` | 交易所名称 | 交易所名称。 |
| `market` | `text` | 市场 | 交易或业务所属的市场标识。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `result_json` | `text` | 结果JSON | 处理结果的 JSON 表示。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

### 表：`connector_hub.settlement_outbox`

- 表中文名：结算出站队列
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `fill_id` | `text` | 成交ID | 关联的成交唯一标识。 |
| `trace_id` | `text` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `request_id` | `text` | 请求ID | 请求级唯一标识。 |
| `payload_json` | `text` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `attempt_count` | `integer` | 尝试次数 | 已执行的尝试次数。 |
| `last_error` | `text` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `last_attempt_at` | `text` | 最近一次尝试时间 | 最近一次执行尝试的时间。 |
| `next_retry_at` | `text` | 下次重试时间 | 下一次允许重试的时间。 |
| `response_json` | `text` | 响应JSON | 外部响应内容的 JSON 表示。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

## 模式：`core_data`

- 中文名：核心数据
- 表数量：2

### 表：`core_data.events`

- 表中文名：事件
- 字段数量：6

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |
| `tenant_id` | `text` | 租户ID | 租户唯一标识。 |

### 表：`core_data.exports`

- 表中文名：导出
- 字段数量：6

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |
| `tenant_id` | `text` | 租户ID | 租户唯一标识。 |

## 模式：`decision_producer`

- 中文名：决策产出
- 表数量：4

### 表：`decision_producer.policy_states`

- 表中文名：策略状态
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`decision_producer.producer_status`

- 表中文名：产出status
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`decision_producer.publish_events`

- 表中文名：发布事件
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`decision_producer.stream_cursors`

- 表中文名：流游标
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

## 模式：`log_kms`

- 中文名：日志与密钥管理
- 表数量：5

### 表：`log_kms.audit_events`

- 表中文名：审计事件
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`log_kms.evidence`

- 表中文名：证据
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`log_kms.key_events`

- 表中文名：键事件
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`log_kms.privacy_erase_requests`

- 表中文名：隐私擦除请求
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`log_kms.privacy_keys`

- 表中文名：隐私键
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

## 模式：`notification_service`

- 中文名：通知服务
- 表数量：6

### 表：`notification_service.channel_subscriptions`

- 表中文名：渠道订阅
- 字段数量：15

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `source_app` | `text` | 来源应用 | 创建该记录的来源应用。 |
| `external_user_id` | `text` | 外部用户ID | 外部系统中的用户标识。 |
| `channel` | `text` | 渠道 | 渠道。 |
| `endpoint` | `text` | 端点 | 通知或回调目标端点。 |
| `endpoint_hash` | `text` | 端点哈希 | 端点内容的哈希值。 |
| `endpoint_meta_json` | `text` | 端点元数据JSON | 端点附加信息的 JSON 表示。 |
| `product` | `text` | 产品 | 所属产品。 |
| `scene` | `text` | 场景 | 业务场景。 |
| `category` | `text` | 分类 | 所属分类。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |
| `last_subscribed_at` | `text` | 最近一次订阅时间 | 最近一次订阅时间。 |
| `last_unsubscribed_at` | `text` | 最近一次取消订阅时间 | 最近一次取消订阅时间。 |

### 表：`notification_service.deliveries`

- 表中文名：投递
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `message_id` | `text` | 消息ID | 消息唯一标识。 |
| `channel` | `text` | 渠道 | 渠道。 |
| `endpoint` | `text` | 端点 | 通知或回调目标端点。 |
| `endpoint_hash` | `text` | 端点哈希 | 端点内容的哈希值。 |
| `endpoint_meta_json` | `text` | 端点元数据JSON | 端点附加信息的 JSON 表示。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `attempts` | `integer` | 尝试次数 | 已执行的尝试次数。 |
| `next_retry_at` | `text` | 下次重试时间 | 下一次允许重试的时间。 |
| `last_error` | `text` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `provider_message_id` | `text` | 服务商消息ID | 第三方服务商返回的消息标识。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

### 表：`notification_service.notifications`

- 表中文名：通知
- 字段数量：16

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `message_id` | `text` | 消息ID | 消息唯一标识。 |
| `source_app` | `text` | 来源应用 | 创建该记录的来源应用。 |
| `request_id` | `text` | 请求ID | 请求级唯一标识。 |
| `product` | `text` | 产品 | 所属产品。 |
| `scene` | `text` | 场景 | 业务场景。 |
| `event_type` | `text` | 事件type | 事件type。 |
| `severity` | `text` | 严重级别 | 严重级别。 |
| `title` | `text` | 标题 | 展示用标题。 |
| `content_text` | `text` | 纯文本内容 | 纯文本格式的消息内容。 |
| `content_html` | `text` | HTML内容 | HTML 格式的消息内容。 |
| `payload_json` | `text` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `route_ids_json` | `text` | 路由ID列表JSON | 路由 ID 列表的 JSON 表示。 |
| `aggregate_status` | `text` | 聚合状态 | 汇总后的整体状态。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

### 表：`notification_service.wa_bind_requests`

- 表中文名：WhatsApp绑定请求
- 字段数量：14

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `bind_code` | `text` | 绑定码 | 用于完成绑定流程的验证码或绑定码。 |
| `source_app` | `text` | 来源应用 | 创建该记录的来源应用。 |
| `external_user_id` | `text` | 外部用户ID | 外部系统中的用户标识。 |
| `product` | `text` | 产品 | 所属产品。 |
| `category` | `text` | 分类 | 所属分类。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `wa_chat_id` | `text` | WhatsApp会话ID | WhatsApp 会话标识。 |
| `wa_phone` | `text` | WhatsApp手机号 | WhatsApp 绑定手机号。 |
| `requested_payload_json` | `text` | 请求载荷JSON | 发起请求时提交的载荷 JSON。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `expires_at` | `text` | 过期时间 | 记录到期时间。 |
| `confirmed_at` | `text` | 确认时间 | 确认完成时间。 |
| `cancelled_at` | `text` | 取消时间 | 取消发生时间。 |

### 表：`notification_service.wa_delivery_logs`

- 表中文名：WhatsApp投递日志
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `request_id` | `text` | 请求ID | 请求级唯一标识。 |
| `source_app` | `text` | 来源应用 | 创建该记录的来源应用。 |
| `wa_chat_id` | `text` | WhatsApp会话ID | WhatsApp 会话标识。 |
| `product` | `text` | 产品 | 所属产品。 |
| `category` | `text` | 分类 | 所属分类。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `provider_message_id` | `text` | 服务商消息ID | 第三方服务商返回的消息标识。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |

### 表：`notification_service.wa_subscriptions`

- 表中文名：WhatsApp订阅
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `bigint` | 主键ID | 当前表的主键标识。 |
| `source_app` | `text` | 来源应用 | 创建该记录的来源应用。 |
| `external_user_id` | `text` | 外部用户ID | 外部系统中的用户标识。 |
| `wa_chat_id` | `text` | WhatsApp会话ID | WhatsApp 会话标识。 |
| `wa_phone` | `text` | WhatsApp手机号 | WhatsApp 绑定手机号。 |
| `product` | `text` | 产品 | 所属产品。 |
| `category` | `text` | 分类 | 所属分类。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |
| `last_bind_at` | `text` | 最近一次绑定时间 | 最近一次绑定时间。 |
| `last_unsubscribe_at` | `text` | 最近一次unsubscribe时间 | 最近一次unsubscribe时间。 |

## 模式：`operator_portal`

- 中文名：运营门户
- 表数量：7

### 表：`operator_portal.codes`

- 表中文名：编码
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.commands`

- 表中文名：指令
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.idempotency_records`

- 表中文名：幂等记录
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.node_status`

- 表中文名：节点status
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.resource_locks`

- 表中文名：资源锁
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.risk_approvals`

- 表中文名：风控审批
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`operator_portal.sessions`

- 表中文名：会话
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

## 模式：`order_gateway`

- 中文名：订单网关
- 表数量：4

### 表：`order_gateway.intents`

- 表中文名：意图
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `intent_id` | `text` | 意图ID | 关联的意图唯一标识。 |
| `trace_id` | `text` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `request_id` | `text` | 请求ID | 请求级唯一标识。 |
| `idempotency_key` | `text` | 幂等键 | 用于避免重复处理的幂等标识。 |
| `node_id` | `text` | 节点ID | 关联的节点唯一标识。 |
| `lobster_subaccount_id` | `text` | Lobster子账户ID | 关联的Lobster子账户唯一标识。 |
| `target_notional` | `double precision` | 目标名义金额 | 目标名义金额。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `payload_json` | `text` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `connector_json` | `text` | 连接器JSON | 连接器的结构化 JSON 数据。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

### 表：`order_gateway.proposal_batches`

- 表中文名：提案batches
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `batch_id` | `text` | 批次ID | 关联的批次唯一标识。 |
| `bucket_key` | `text` | 分桶键 | 分桶键。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `total_notional` | `double precision` | 总量名义金额 | 总量名义金额。 |
| `execute_budget` | `double precision` | 执行budget | 执行budget。 |
| `executed_notional` | `double precision` | 已执行名义金额 | 已执行名义金额。 |
| `rejected_notional` | `double precision` | 已拒绝名义金额 | 已拒绝名义金额。 |
| `blocked_notional` | `double precision` | 阻断名义金额 | 阻断名义金额。 |
| `child_count` | `integer` | 子项数量 | 子项数量或次数。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

### 表：`order_gateway.proposal_children`

- 表中文名：提案子项
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `child_id` | `text` | 子项ID | 关联的子项唯一标识。 |
| `proposal_id` | `text` | 提案ID | 关联的提案唯一标识。 |
| `batch_id` | `text` | 批次ID | 关联的批次唯一标识。 |
| `intent_id` | `text` | 意图ID | 关联的意图唯一标识。 |
| `allocated_notional` | `double precision` | 已分配名义金额 | 已分配名义金额。 |
| `outcome` | `text` | 结果 | 结果。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `connector_order_id` | `text` | 连接器订单ID | 关联的连接器订单唯一标识。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |

### 表：`order_gateway.proposals`

- 表中文名：提案
- 字段数量：17

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `proposal_id` | `text` | 提案ID | 关联的提案唯一标识。 |
| `idempotency_key` | `text` | 幂等键 | 用于避免重复处理的幂等标识。 |
| `trace_id` | `text` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `request_id` | `text` | 请求ID | 请求级唯一标识。 |
| `bucket_key` | `text` | 分桶键 | 分桶键。 |
| `payload_json` | `text` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `target_notional` | `double precision` | 目标名义金额 | 目标名义金额。 |
| `remaining_notional` | `double precision` | 剩余名义金额 | 剩余名义金额。 |
| `executed_notional` | `double precision` | 已执行名义金额 | 已执行名义金额。 |
| `rejected_notional` | `double precision` | 已拒绝名义金额 | 已拒绝名义金额。 |
| `blocked_notional` | `double precision` | 阻断名义金额 | 阻断名义金额。 |
| `carryover_notional` | `double precision` | 结转名义金额 | 结转名义金额。 |
| `last_batch_id` | `text` | 最近一次批次ID | 关联的最近一次批次唯一标识。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `updated_at` | `text` | 更新时间 | 记录最后更新时间。 |

## 模式：`settlement_risk`

- 中文名：结算风控
- 表数量：7

### 表：`settlement_risk.feedback_outbox`

- 表中文名：反馈出站队列
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.fills`

- 表中文名：成交
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.intent_index`

- 表中文名：意图index
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.payouts`

- 表中文名：打款
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.reconcile_runs`

- 表中文名：reconcileruns
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.risk_action_simulations`

- 表中文名：风控动作simulations
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`settlement_risk.risk_actions`

- 表中文名：风控动作
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

## 模式：`signal_server`

- 中文名：信号服务
- 表数量：3

### 表：`signal_server.signal_acks`

- 表中文名：信号acks
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`signal_server.signal_events`

- 表中文名：信号事件
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

### 表：`signal_server.signal_nodes`

- 表中文名：信号节点
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `item_key` | `text` | 业务键 | 用于定位业务对象的唯一键。 |
| `item_json` | `jsonb` | 业务对象JSON | 业务对象的完整 JSON 内容。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |
| `row_version` | `bigint` | 行版本号 | 用于并发控制或变更跟踪的版本号。 |

## 模式：`usstocks_runtime`

- 中文名：美股运行时
- 表数量：7

### 表：`usstocks_runtime.backtest_runs`

- 表中文名：回测runs
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `run_id` | `text` | runID | 关联的run唯一标识。 |
| `created_at` | `text` | 创建时间 | 记录创建时间。 |
| `mode` | `text` | 模式 | 运行或处理模式。 |
| `symbols_json` | `jsonb` | 标的代码JSON | 标的代码的结构化 JSON 数据。 |
| `config_json` | `jsonb` | 配置JSON | 结构化配置内容。 |
| `metrics_json` | `jsonb` | metricsJSON | metrics的结构化 JSON 数据。 |
| `note` | `text` | 备注 | 补充说明信息。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.backtest_trades`

- 表中文名：回测交易
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `trade_id` | `text` | 交易ID | 关联的交易唯一标识。 |
| `run_id` | `text` | runID | 关联的run唯一标识。 |
| `symbol` | `text` | 标的代码 | 交易标的或证券代码。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `ts_iso` | `text` | ISO时间 | ISO 8601 格式时间。 |
| `side` | `text` | 方向 | 买卖或处理方向。 |
| `qty` | `text` | 数量 | 数量值。 |
| `price` | `text` | 价格 | 价格值。 |
| `fee` | `text` | fee | fee。 |
| `realized_pnl` | `text` | 已实现盈亏 | 已实现盈亏。 |
| `action` | `text` | 动作 | 动作类型。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.order_intents`

- 表中文名：订单意图
- 字段数量：16

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `intent_id` | `text` | 意图ID | 关联的意图唯一标识。 |
| `signal_id` | `text` | 信号ID | 关联的信号唯一标识。 |
| `symbol` | `text` | 标的代码 | 交易标的或证券代码。 |
| `created_ts_ms` | `bigint` | 创建毫秒时间戳 | 创建的毫秒级时间戳。 |
| `created_ts_iso` | `text` | 创建ISO时间 | 创建的 ISO 时间字符串。 |
| `side` | `text` | 方向 | 买卖或处理方向。 |
| `qty` | `text` | 数量 | 数量值。 |
| `notional` | `text` | 名义金额 | 名义金额。 |
| `order_type` | `text` | 订单type | 订单type。 |
| `limit_price` | `text` | 限制价格 | 限制价格。 |
| `time_in_force` | `text` | timeinforce | timeinforce。 |
| `allow_extended_hours` | `integer` | 是否允许扩展时段 | 布尔标记，表示是否允许扩展时段。 |
| `reduce_only` | `integer` | 减仓仅 | 减仓仅。 |
| `client_tag` | `text` | 客户端tag | 客户端tag。 |
| `status` | `text` | 状态 | 当前记录所处的状态。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.portfolio_snapshots`

- 表中文名：组合快照
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `snapshot_id` | `text` | 快照ID | 关联的快照唯一标识。 |
| `run_id` | `text` | runID | 关联的run唯一标识。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `ts_iso` | `text` | ISO时间 | ISO 8601 格式时间。 |
| `equity` | `text` | equity | equity。 |
| `cash` | `text` | 现金 | 现金。 |
| `settled_cash` | `text` | 已结算现金 | 已结算现金。 |
| `unsettled_cash` | `text` | 未结算现金 | 未结算现金。 |
| `buying_power` | `text` | 可买入购买力 | 可买入购买力。 |
| `gross_exposure` | `text` | 毛敞口 | 毛敞口。 |
| `net_exposure` | `text` | 净敞口 | 净敞口。 |
| `position_count` | `integer` | 持仓数量 | 持仓数量或次数。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.positions_snapshots`

- 表中文名：持仓快照
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `snapshot_id` | `text` | 快照ID | 关联的快照唯一标识。 |
| `run_id` | `text` | runID | 关联的run唯一标识。 |
| `symbol` | `text` | 标的代码 | 交易标的或证券代码。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `ts_iso` | `text` | ISO时间 | ISO 8601 格式时间。 |
| `qty` | `text` | 数量 | 数量值。 |
| `avg_entry_price` | `text` | avg入场价格 | avg入场价格。 |
| `market_price` | `text` | 市场价格 | 市场价格。 |
| `market_value` | `text` | 市场值 | 市场值。 |
| `unrealized_pnl` | `text` | 未实现盈亏 | 未实现盈亏。 |
| `realized_pnl` | `text` | 已实现盈亏 | 已实现盈亏。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.predictions`

- 表中文名：预测
- 字段数量：18

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `prediction_id` | `text` | 预测ID | 关联的预测唯一标识。 |
| `run_id` | `text` | runID | 关联的run唯一标识。 |
| `symbol` | `text` | 标的代码 | 交易标的或证券代码。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `ts_iso` | `text` | ISO时间 | ISO 8601 格式时间。 |
| `timeframe` | `text` | 时间框架 | 时间框架。 |
| `session` | `text` | 会话 | 会话。 |
| `regime_state` | `text` | 市场状态状态 | 市场状态状态。 |
| `direction` | `text` | 方向 | 方向。 |
| `score` | `double precision` | 评分 | 模型或规则输出的评分。 |
| `confidence` | `double precision` | 置信度 | 模型输出结果的置信度。 |
| `expected_horizon_bars` | `integer` | 预期观察期K线数 | 预期观察期K线数。 |
| `model_name` | `text` | 模型名称 | 模型名称。 |
| `model_version` | `text` | 模型版本 | 模型版本。 |
| `feature_version` | `text` | 特征版本 | 特征版本。 |
| `feature_ref` | `text` | 特征引用 | 特征引用。 |
| `feature_snapshot_json` | `jsonb` | 特征快照JSON | 特征快照的结构化 JSON 数据。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`usstocks_runtime.signals`

- 表中文名：信号
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `signal_id` | `text` | 信号ID | 关联的信号唯一标识。 |
| `prediction_id` | `text` | 预测ID | 关联的预测唯一标识。 |
| `symbol` | `text` | 标的代码 | 交易标的或证券代码。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `ts_iso` | `text` | ISO时间 | ISO 8601 格式时间。 |
| `action` | `text` | 动作 | 动作类型。 |
| `reason` | `text` | 原因 | 状态变化或处理结果的原因说明。 |
| `score` | `double precision` | 评分 | 模型或规则输出的评分。 |
| `confidence` | `double precision` | 置信度 | 模型输出结果的置信度。 |
| `session` | `text` | 会话 | 会话。 |
| `blocked_by` | `text` | 阻断by | 阻断by。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

## 模式：`public`

- 中文名：公共模式
- 表数量：80

### 表：`public.account_links`

- 表中文名：账户关联
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `portal_user_id` | `integer` | 门户用户ID | 关联的门户用户唯一标识。 |
| `lobster_user_id` | `character varying(64)` | Lobster用户ID | 关联的Lobster用户唯一标识。 |
| `link_status` | `character varying(24)` | 关联状态 | 关联状态。 |
| `linked_at` | `timestamp without time zone` | linked时间 | linked时间。 |
| `last_sync_at` | `timestamp without time zone` | 最近一次sync时间 | 最近一次sync时间。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.activation_revenue_records`

- 表中文名：激活收益记录
- 字段数量：18

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `revenue_id` | `character varying(96)` | 收益ID | 关联的收益唯一标识。 |
| `follower_user_id` | `integer` | 跟单用户用户ID | 关联的跟单用户用户唯一标识。 |
| `leader_user_id` | `integer` | 带单节点用户ID | 关联的带单节点用户唯一标识。 |
| `leader_node_id` | `character varying(96)` | 带单节点节点ID | 关联的带单节点节点唯一标识。 |
| `source_ref_id` | `character varying(120)` | 来源引用ID | 来源引用ID。 |
| `share_ratio_dec` | `character varying(64)` | shareratio十进制定点值 | shareratio的十进制定点数值。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `escrow_wallet_address` | `character varying(96)` | 托管钱包地址 | 托管钱包地址。 |
| `escrow_wallet_ciphertext` | `text` | 托管钱包ciphertext | 托管钱包ciphertext。 |
| `payout_ref_id` | `character varying(120)` | 打款引用ID | 打款引用ID。 |
| `payout_tx_hash` | `character varying(120)` | 打款tx哈希 | 打款tx哈希值。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |
| `entitlement_key` | `character varying(220)` | 权益键 | 权益键。 |

### 表：`public.agent_group`

- 表中文名：智能体分组
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `name` | `character varying(64)` | 名称 | 展示名称。 |
| `description` | `character varying(255)` | 描述 | 详细描述说明。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.agent_group_member`

- 表中文名：智能体分组成员
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `group_id` | `integer` | 分组ID | 关联的分组唯一标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `agent_id` | `character varying(64)` | 智能体ID | 关联的智能体唯一标识。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.alerts`

- 表中文名：告警
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `level` | `character varying(8)` | 级别 | 分级或严重程度级别。 |
| `title` | `character varying(255)` | 标题 | 展示用标题。 |
| `detail` | `text` | 详情 | 详细说明内容。 |
| `status` | `character varying(16)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.asset_node`

- 表中文名：资产节点
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `name` | `character varying(128)` | 名称 | 展示名称。 |
| `owner_user_id` | `character varying(64)` | 所有者用户ID | 关联的所有者用户唯一标识。 |
| `mac_fingerprint` | `character varying(128)` | MAC指纹 | MAC指纹。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `openclaw_version` | `character varying(32)` | OpenClaw版本 | OpenClaw版本。 |
| `openclaw_state` | `character varying(32)` | OpenClaw状态 | OpenClaw状态。 |
| `openclaw_uptime_sec` | `integer` | OpenClaw运行时长sec | OpenClaw运行时长sec。 |
| `agent_capacity` | `integer` | 智能体容量 | 智能体容量。 |
| `last_seen_at` | `timestamp with time zone` | 最近一次见到时间 | 最近一次见到时间。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.audit_command`

- 表中文名：审计指令
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `actor` | `character varying(64)` | actor | actor。 |
| `source` | `character varying(32)` | 来源 | 来源。 |
| `command_tpl` | `character varying(128)` | 指令tpl | 指令tpl。 |
| `params_json` | `text` | paramsJSON | params的结构化 JSON 数据。 |
| `exit_code` | `integer` | exit编码 | exit编码。 |
| `trace_id` | `character varying(64)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.audit_logs`

- 表中文名：审计日志
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `ts` | `timestamp without time zone` | 时间戳 | 通用时间戳字段。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `action` | `character varying(80)` | 动作 | 动作类型。 |
| `target_type` | `character varying(60)` | 目标type | 目标type。 |
| `target_id` | `character varying(120)` | 目标ID | 关联的目标唯一标识。 |
| `trace_id` | `character varying(120)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `request_id` | `character varying(120)` | 请求ID | 请求级唯一标识。 |
| `payload_json` | `json` | 载荷JSON | 请求或事件的原始载荷 JSON。 |

### 表：`public.audit_support_session`

- 表中文名：审计支持会话
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `session_id` | `character varying(64)` | 会话ID | 关联的会话唯一标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `perm_level` | `character varying(32)` | perm级别 | perm级别。 |
| `status` | `character varying(16)` | 状态 | 当前记录所处的状态。 |
| `expired_at` | `timestamp with time zone` | 失效时间 | 记录失效时间。 |
| `revoked_at` | `timestamp with time zone` | revoked时间 | revoked时间。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.biz_cost_config`

- 表中文名：业务成本config
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `power_price_per_kwh` | `double precision` | 购买力价格每千瓦时 | 购买力价格每千瓦时。 |
| `token_fx_rate` | `double precision` | 令牌汇率rate | 令牌汇率rate。 |
| `depreciation_daily` | `double precision` | 折旧每日 | 折旧每日。 |
| `updated_by` | `character varying(64)` | 更新by | 更新by。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.biz_pnl_minute`

- 表中文名：业务盈亏分钟
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `minute_ts` | `timestamp with time zone` | 分钟时间戳 | 分钟的时间戳。 |
| `gross_pnl` | `double precision` | 毛盈亏 | 毛盈亏。 |
| `token_cost` | `double precision` | 令牌成本 | 令牌成本。 |
| `power_cost` | `double precision` | 购买力成本 | 购买力成本。 |
| `depreciation_cost` | `double precision` | 折旧成本 | 折旧成本。 |
| `net_pnl` | `double precision` | 净盈亏 | 净盈亏。 |
| `efficiency_score` | `double precision` | 效率评分 | 效率评分。 |

### 表：`public.broker_credential_envelopes`

- 表中文名：券商凭证封装
- 字段数量：6

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `mac_address` | `character varying(64)` | MAC地址 | MAC地址。 |
| `key_fingerprint` | `character varying(120)` | 键指纹 | 键指纹。 |
| `encrypted_payload` | `text` | encrypted载荷 | encrypted载荷。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.claim_sessions`

- 表中文名：认领会话
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `claim_session_id` | `character varying(80)` | 认领会话ID | 关联的认领会话唯一标识。 |
| `slot_id` | `integer` | 槽位ID | 关联的槽位唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `qr_content` | `character varying(500)` | qr内容 | qr内容。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.coupons`

- 表中文名：优惠券
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `code` | `character varying(80)` | 编码 | 编码。 |
| `discount_rate` | `double precision` | discountrate | discountrate。 |
| `status` | `character varying(20)` | 状态 | 当前记录所处的状态。 |
| `max_uses` | `integer` | 最大uses | 最大uses。 |
| `used_count` | `integer` | used数量 | used数量或次数。 |
| `sales_agent_id` | `integer` | 销售智能体ID | 关联的销售智能体唯一标识。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.custody_wallets`

- 表中文名：托管钱包
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `address` | `character varying(96)` | 地址 | 地址字符串。 |
| `private_key_ciphertext` | `text` | private键ciphertext | private键ciphertext。 |
| `key_ref` | `character varying(120)` | 键引用 | 键引用。 |
| `active` | `boolean` | active | 布尔标记，表示active是否成立。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.deployment_download_logs`

- 表中文名：部署下载日志
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `script_id` | `integer` | 脚本ID | 关联的脚本唯一标识。 |
| `order_no` | `character varying(64)` | 订单no | 订单no。 |
| `ip_address` | `character varying(64)` | IP地址 | 来源 IP 地址。 |
| `user_agent` | `character varying(255)` | 用户代理 | 客户端 User-Agent。 |
| `download_result` | `character varying(24)` | 下载结果 | 下载结果。 |
| `trace_id` | `character varying(80)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.deployment_scripts`

- 表中文名：部署脚本
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `script_code` | `character varying(80)` | 脚本编码 | 脚本编码。 |
| `worker_release_id` | `integer` | Worker发布ID | 关联的Worker发布唯一标识。 |
| `target_platform` | `character varying(40)` | 目标平台 | 目标平台。 |
| `content_ref` | `character varying(255)` | 内容引用 | 内容引用。 |
| `checksum_sha256` | `character varying(128)` | 校验和SHA256 | 校验和SHA256。 |
| `signature_ref` | `character varying(255)` | 签名引用 | 签名引用。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.deposit_intents`

- 表中文名：充值意图
- 字段数量：14

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `intent_id` | `character varying(80)` | 意图ID | 关联的意图唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `claim_session_id` | `character varying(80)` | 认领会话ID | 关联的认领会话唯一标识。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `asset_code` | `character varying(24)` | 资产编码 | 资产编码。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `receive_address` | `character varying(96)` | 接收地址 | 接收地址。 |
| `source_type` | `character varying(20)` | 来源type | 来源type。 |
| `source_ref` | `character varying(120)` | 来源引用 | 来源引用。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.deposit_transactions`

- 表中文名：充值交易
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `intent_id` | `character varying(80)` | 意图ID | 关联的意图唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `tx_hash` | `character varying(120)` | tx哈希 | tx哈希值。 |
| `log_index` | `integer` | 日志index | 日志index。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `confirmations` | `integer` | 确认数 | 确认数。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.edge_intent_ledger`

- 表中文名：边缘意图台账
- 字段数量：31

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `intent_id` | `character varying(128)` | 意图ID | 关联的意图唯一标识。 |
| `trace_id` | `character varying(120)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `request_id` | `character varying(120)` | 请求ID | 请求级唯一标识。 |
| `idempotency_key` | `character varying(160)` | 幂等键 | 用于避免重复处理的幂等标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `owner_node_id` | `character varying(96)` | 所有者节点ID | 关联的所有者节点唯一标识。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `strategy_id` | `character varying(128)` | 策略ID | 关联的策略唯一标识。 |
| `model_id` | `character varying(128)` | 模型ID | 关联的模型唯一标识。 |
| `source_kind` | `character varying(64)` | 来源kind | 来源kind。 |
| `origin_engine` | `character varying(128)` | 来源引擎 | 来源引擎。 |
| `exchange` | `character varying(64)` | 交易所 | 交易所。 |
| `market` | `character varying(128)` | 市场 | 交易或业务所属的市场标识。 |
| `side` | `character varying(16)` | 方向 | 买卖或处理方向。 |
| `action` | `character varying(16)` | 动作 | 动作类型。 |
| `signal_state` | `character varying(16)` | 信号状态 | 信号状态。 |
| `allocation_required` | `boolean` | 分配必需 | 布尔标记，表示分配必需是否成立。 |
| `amount_hint_dec` | `character varying(64)` | 金额提示十进制定点值 | 金额提示的十进制定点数值。 |
| `amount_usdc_hint_dec` | `character varying(64)` | 金额usdc提示十进制定点值 | 金额usdc提示的十进制定点数值。 |
| `allocated_amount_dec` | `character varying(64)` | 已分配金额十进制定点值 | 已分配金额的十进制定点数值。 |
| `allocated_amount_usdc_dec` | `character varying(64)` | 已分配金额usdc十进制定点值 | 已分配金额usdc的十进制定点数值。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `reason` | `character varying(255)` | 原因 | 状态变化或处理结果的原因说明。 |
| `gateway_intent_id` | `character varying(128)` | 网关意图ID | 关联的网关意图唯一标识。 |
| `gateway_status` | `character varying(32)` | 网关状态 | 网关状态。 |
| `payload_json` | `json` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `allocator_meta_json` | `json` | 分配器元数据JSON | 分配器元数据的结构化 JSON 数据。 |
| `gateway_result_json` | `json` | 网关结果JSON | 网关结果的结构化 JSON 数据。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.email_codes`

- 表中文名：邮箱编码
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `email` | `character varying(255)` | 邮箱 | 邮箱地址。 |
| `biz_type` | `character varying(40)` | 业务type | 业务type。 |
| `code_hash` | `character varying(128)` | 编码哈希 | 编码哈希值。 |
| `nonce` | `character varying(64)` | 随机数 | 随机数。 |
| `expire_at` | `timestamp without time zone` | 过期时间 | 记录到期时间。 |
| `consumed` | `boolean` | 已消耗 | 布尔标记，表示已消耗是否成立。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.email_verification_codes`

- 表中文名：邮箱验证编码
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `biz_type` | `character varying(40)` | 业务type | 业务type。 |
| `email` | `character varying(255)` | 邮箱 | 邮箱地址。 |
| `code_hash` | `character varying(128)` | 编码哈希 | 编码哈希值。 |
| `code_nonce` | `character varying(64)` | 编码随机数 | 编码随机数。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.follower_key_records`

- 表中文名：跟单用户键记录
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `key_id` | `character varying(80)` | 键ID | 关联的键唯一标识。 |
| `key_material_sha256` | `character varying(128)` | 键materialSHA256 | 键materialSHA256。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `bound_at_ts` | `numeric(18,3)` | 已绑定at时间戳 | 已绑定at的时间戳。 |
| `key_source` | `character varying(32)` | 键来源 | 键来源。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.follower_node_bindings`

- 表中文名：跟单用户节点绑定
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `dashboard_url` | `character varying(255)` | 看板链接 | 看板链接。 |
| `mode` | `character varying(24)` | 模式 | 运行或处理模式。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `follower_key_id` | `character varying(80)` | 跟单用户键ID | 关联的跟单用户键唯一标识。 |
| `bound_at_ts` | `numeric(18,3)` | 已绑定at时间戳 | 已绑定at的时间戳。 |
| `unbind_cooldown_until_ts` | `numeric(18,3)` | 解绑冷却期截至时间戳 | 解绑冷却期截至的时间戳。 |
| `unbind_reason` | `character varying(255)` | 解绑原因 | 解绑原因。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.group_campaigns`

- 表中文名：分组活动
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `name` | `character varying(120)` | 名称 | 展示名称。 |
| `discount_rate` | `double precision` | discountrate | discountrate。 |
| `duration_days` | `integer` | 持续时长days | 持续时长days。 |
| `target_people` | `integer` | 目标人数 | 目标人数。 |
| `current_people` | `integer` | 当前人数 | 当前人数。 |
| `is_active` | `boolean` | 是否active | 布尔标记，表示是否active。 |
| `start_at` | `timestamp without time zone` | start时间 | start时间。 |
| `end_at` | `timestamp without time zone` | end时间 | end时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.hot_wallet_route_records`

- 表中文名：热钱包路由记录
- 字段数量：18

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `route_id` | `character varying(96)` | 路由ID | 关联的路由唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `direction` | `character varying(32)` | 方向 | 方向。 |
| `trigger` | `character varying(96)` | 触发方式 | 触发方式。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `from_address` | `character varying(96)` | 来源地址 | 来源地址。 |
| `to_address` | `character varying(96)` | 目标地址 | 目标地址。 |
| `tx_hash` | `character varying(120)` | tx哈希 | tx哈希值。 |
| `error_code` | `character varying(80)` | 错误编码 | 错误编码。 |
| `error_message` | `character varying(255)` | 错误消息 | 错误消息。 |
| `ref_id` | `character varying(120)` | refID | 关联的ref唯一标识。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.install_task_events`

- 表中文名：安装任务事件
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `install_task_id` | `integer` | 安装任务ID | 关联的安装任务唯一标识。 |
| `event_type` | `character varying(32)` | 事件type | 事件type。 |
| `event_payload_json` | `text` | 事件载荷JSON | 事件载荷的结构化 JSON 数据。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.install_tasks`

- 表中文名：安装任务
- 字段数量：15

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `task_no` | `character varying(80)` | 任务no | 任务no。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `order_no` | `character varying(64)` | 订单no | 订单no。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `worker_release_id` | `integer` | Worker发布ID | 关联的Worker发布唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `install_token_hash` | `character varying(128)` | 安装令牌哈希 | 安装令牌哈希值。 |
| `issued_at` | `timestamp without time zone` | 签发时间 | 签发时间。 |
| `finished_at` | `timestamp without time zone` | 完成时间 | 完成时间。 |
| `last_error` | `character varying(255)` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `trace_id` | `character varying(80)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.key_application_requests`

- 表中文名：键申请请求
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `request_no` | `character varying(80)` | 请求no | 请求no。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `order_no` | `character varying(64)` | 订单no | 订单no。 |
| `key_type` | `character varying(24)` | 键type | 键type。 |
| `target_node_id` | `character varying(96)` | 目标节点ID | 关联的目标节点唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `review_note` | `character varying(255)` | 审核备注 | 审核备注。 |
| `lobster_case_ref` | `character varying(120)` | Lobster案件引用 | Lobster案件引用。 |
| `trace_id` | `character varying(80)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.ledger_accounts`

- 表中文名：台账账户
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `account_code` | `character varying(64)` | 账户编码 | 账户编码。 |
| `asset_code` | `character varying(24)` | 资产编码 | 资产编码。 |
| `balance_dec` | `character varying(64)` | 余额十进制定点值 | 余额的十进制定点数值。 |
| `balance_atomic` | `bigint` | 余额原子单位值 | 余额的原子单位数值。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.ledger_entries`

- 表中文名：台账入场
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `idempotency_key` | `character varying(160)` | 幂等键 | 用于避免重复处理的幂等标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `asset_code` | `character varying(24)` | 资产编码 | 资产编码。 |
| `entry_type` | `character varying(40)` | 入场type | 入场type。 |
| `dr_account` | `character varying(64)` | 借方账户 | 借方账户。 |
| `cr_account` | `character varying(64)` | 贷方账户 | 贷方账户。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `ref_id` | `character varying(120)` | refID | 关联的ref唯一标识。 |
| `memo` | `character varying(255)` | 备注 | 备注。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.lobster_slots`

- 表中文名：Lobster槽位
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `slot_code` | `character varying(64)` | 槽位编码 | 槽位编码。 |
| `title` | `character varying(120)` | 标题 | 展示用标题。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `current_user_id` | `integer` | 当前用户ID | 关联的当前用户唯一标识。 |
| `reserved_until_ts` | `numeric(18,3)` | 预留截至时间戳 | 预留截至的时间戳。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |
| `owner_user_id` | `integer` | 所有者用户ID | 关联的所有者用户唯一标识。 |
| `slot_no` | `integer` | 槽位no | 槽位no。 |

### 表：`public.log_archive_catalog`

- 表中文名：日志归档catalog
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `archive_id` | `character varying(80)` | 归档ID | 关联的归档唯一标识。 |
| `service` | `character varying(64)` | 服务 | 服务。 |
| `log_type` | `character varying(64)` | 日志type | 日志type。 |
| `date_from` | `character varying(20)` | 日期来源 | 日期来源。 |
| `date_to` | `character varying(20)` | 日期目标 | 日期目标。 |
| `object_path` | `character varying(300)` | 对象path | 对象path。 |
| `sha256` | `character varying(128)` | SHA256 | SHA256。 |
| `record_count` | `integer` | 记录数量 | 记录数量或次数。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.log_local_action_chain`

- 表中文名：日志本地动作链
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `event_id` | `integer` | 事件ID | 关联的事件唯一标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `prev_event_hash` | `character varying(128)` | 前一个事件哈希 | 前一个事件哈希值。 |
| `event_hash` | `character varying(128)` | 事件哈希 | 事件哈希值。 |
| `signature` | `character varying(256)` | 签名 | 签名。 |
| `chain_ts` | `timestamp with time zone` | 链时间戳 | 链的时间戳。 |

### 表：`public.log_local_action_event`

- 表中文名：日志本地动作事件
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `rule_id` | `character varying(64)` | 规则ID | 关联的规则唯一标识。 |
| `matched_pattern` | `character varying(255)` | 匹配pattern | 匹配pattern。 |
| `action_taken` | `character varying(128)` | 动作执行 | 动作执行。 |
| `target_agent` | `character varying(64)` | 目标智能体 | 目标智能体。 |
| `event_ts` | `timestamp with time zone` | 事件时间戳 | 事件的时间戳。 |

### 表：`public.master_settlement_events`

- 表中文名：主结算事件
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `source` | `character varying(64)` | 来源 | 来源。 |
| `source_event_id` | `character varying(128)` | 来源事件ID | 关联的来源事件唯一标识。 |
| `market` | `character varying(80)` | 市场 | 交易或业务所属的市场标识。 |
| `settle_tag` | `character varying(16)` | 结算tag | 结算tag。 |
| `settle_pnl` | `character varying(64)` | 结算盈亏 | 结算盈亏。 |
| `pnl_final` | `character varying(64)` | 盈亏最终 | 盈亏最终。 |
| `payload_json` | `json` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `ts_ms` | `bigint` | 毫秒时间戳 | 毫秒级时间戳。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.mirror_configs`

- 表中文名：镜像configs
- 字段数量：24

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `enabled` | `boolean` | 启用 | 布尔标记，表示启用是否成立。 |
| `leverage` | `character varying(64)` | 杠杆 | 杠杆。 |
| `risk_factor` | `character varying(64)` | 风控系数 | 风控系数。 |
| `user_cap_dec` | `character varying(64)` | 用户上限十进制定点值 | 用户上限的十进制定点数值。 |
| `run_since_ts_ms` | `bigint` | run自毫秒时间戳 | run自的毫秒级时间戳。 |
| `stop_requested_ts_ms` | `bigint` | 停止请求时毫秒时间戳 | 停止请求时的毫秒级时间戳。 |
| `stop_pending_settlement` | `boolean` | 停止待处理结算 | 布尔标记，表示停止待处理结算是否成立。 |
| `risk_circuit_open` | `boolean` | 风控熔断open | 布尔标记，表示风控熔断open是否成立。 |
| `risk_circuit_opened_ts_ms` | `bigint` | 风控熔断打开毫秒时间戳 | 风控熔断打开的毫秒级时间戳。 |
| `risk_circuit_hold_until_ts_ms` | `bigint` | 风控熔断保持截至毫秒时间戳 | 风控熔断保持截至的毫秒级时间戳。 |
| `risk_loss_streak` | `integer` | 风控亏损连击 | 风控亏损连击。 |
| `risk_circuit_reason` | `character varying(64)` | 风控熔断原因 | 风控熔断原因。 |
| `daily_target_min_pct` | `character varying(64)` | 每日目标最小百分比 | 每日目标最小百分比数值。 |
| `daily_target_max_pct` | `character varying(64)` | 每日目标最大百分比 | 每日目标最大百分比数值。 |
| `daily_loss_limit_pct` | `character varying(64)` | 每日亏损限制百分比 | 每日亏损限制百分比数值。 |
| `stop_mode` | `character varying(32)` | 停止mode | 停止mode。 |
| `rule_version` | `character varying(32)` | 规则版本 | 规则版本。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |
| `hot_wallet_routed` | `boolean` | 热钱包routed | 布尔标记，表示热钱包routed是否成立。 |
| `hot_wallet_routed_ts_ms` | `bigint` | 热钱包routed毫秒时间戳 | 热钱包routed的毫秒级时间戳。 |
| `hot_wallet_route_ref` | `character varying(120)` | 热钱包路由引用 | 热钱包路由引用。 |

### 表：`public.mirror_daily_stats`

- 表中文名：镜像每日stats
- 字段数量：15

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `trade_date` | `character varying(20)` | 交易日期 | 交易日期。 |
| `baseline_principal_dec` | `character varying(64)` | 基线本金十进制定点值 | 基线本金的十进制定点数值。 |
| `target_pct_dec` | `character varying(64)` | 目标pct十进制定点值 | 目标pct的十进制定点数值。 |
| `target_pnl_dec` | `character varying(64)` | 目标盈亏十进制定点值 | 目标盈亏的十进制定点数值。 |
| `loss_limit_pnl_dec` | `character varying(64)` | 亏损限制盈亏十进制定点值 | 亏损限制盈亏的十进制定点数值。 |
| `realized_pnl_dec` | `character varying(64)` | 已实现盈亏十进制定点值 | 已实现盈亏的十进制定点数值。 |
| `consumed_token_dec` | `character varying(64)` | 已消耗令牌十进制定点值 | 已消耗令牌的十进制定点数值。 |
| `stopped` | `boolean` | 停止 | 布尔标记，表示停止是否成立。 |
| `stop_reason` | `character varying(64)` | 停止原因 | 停止原因。 |
| `stop_event_id` | `character varying(128)` | 停止事件ID | 关联的停止事件唯一标识。 |
| `rule_version` | `character varying(32)` | 规则版本 | 规则版本。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.mirror_positions`

- 表中文名：镜像持仓
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `position_id` | `character varying(80)` | 持仓ID | 关联的持仓唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `market` | `character varying(80)` | 市场 | 交易或业务所属的市场标识。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `exposure_dec` | `character varying(64)` | 敞口十进制定点值 | 敞口的十进制定点数值。 |
| `entry_ref` | `character varying(128)` | 入场引用 | 入场引用。 |
| `closed_ref` | `character varying(128)` | 平仓引用 | 平仓引用。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.mirror_settlement_batches`

- 表中文名：镜像结算batches
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `batch_id` | `character varying(80)` | 批次ID | 关联的批次唯一标识。 |
| `source_event_id` | `character varying(128)` | 来源事件ID | 关联的来源事件唯一标识。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `processed_users` | `integer` | 已处理用户 | 已处理用户。 |
| `success_users` | `integer` | 成功用户 | 成功用户。 |
| `failed_users` | `integer` | 失败用户 | 失败用户。 |
| `rule_version` | `character varying(32)` | 规则版本 | 规则版本。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.mirror_settlement_items`

- 表中文名：镜像结算条目
- 字段数量：22

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `batch_id` | `character varying(80)` | 批次ID | 关联的批次唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `market` | `character varying(80)` | 市场 | 交易或业务所属的市场标识。 |
| `exposure_dec` | `character varying(64)` | 敞口十进制定点值 | 敞口的十进制定点数值。 |
| `principal_dec` | `character varying(64)` | 本金十进制定点值 | 本金的十进制定点数值。 |
| `gross_pnl_dec` | `character varying(64)` | 毛盈亏十进制定点值 | 毛盈亏的十进制定点数值。 |
| `fee_dec` | `character varying(64)` | fee十进制定点值 | fee的十进制定点数值。 |
| `net_pnl_dec` | `character varying(64)` | 净盈亏十进制定点值 | 净盈亏的十进制定点数值。 |
| `token_cost_dec` | `character varying(64)` | 令牌成本十进制定点值 | 令牌成本的十进制定点数值。 |
| `daily_trade_date` | `character varying(20)` | 每日交易日期 | 每日交易日期。 |
| `daily_realized_pnl_dec` | `character varying(64)` | 每日已实现盈亏十进制定点值 | 每日已实现盈亏的十进制定点数值。 |
| `daily_target_pnl_dec` | `character varying(64)` | 每日目标盈亏十进制定点值 | 每日目标盈亏的十进制定点数值。 |
| `daily_stop_reason` | `character varying(64)` | 每日停止原因 | 每日停止原因。 |
| `rule_version` | `character varying(32)` | 规则版本 | 规则版本。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |
| `fee_trade_dec` | `character varying(64)` | fee交易十进制定点值 | fee交易的十进制定点数值。 |
| `fee_platform_dec` | `character varying(64)` | fee平台十进制定点值 | fee平台的十进制定点数值。 |
| `fee_leader_dec` | `character varying(64)` | fee带单节点十进制定点值 | fee带单节点的十进制定点数值。 |

### 表：`public.model_contract_index`

- 表中文名：模型contractindex
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `model_id` | `character varying(128)` | 模型ID | 关联的模型唯一标识。 |
| `model_ver` | `character varying(64)` | 模型ver | 模型ver。 |
| `file_hash` | `character varying(128)` | 文件哈希 | 文件哈希值。 |
| `schema_hash` | `character varying(128)` | 结构哈希 | 结构哈希值。 |
| `runtime_compat` | `character varying(64)` | 运行时兼容性 | 运行时兼容性。 |
| `status` | `character varying(32)` | 状态 | 当前记录所处的状态。 |
| `indexed_at` | `timestamp with time zone` | 已索引时间 | 已索引时间。 |

### 表：`public.model_license_key_records`

- 表中文名：模型许可证键记录
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `key_id` | `character varying(80)` | 键ID | 关联的键唯一标识。 |
| `key_material_sha256` | `character varying(128)` | 键materialSHA256 | 键materialSHA256。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `model_id` | `character varying(128)` | 模型ID | 关联的模型唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `issued_by` | `character varying(96)` | 签发by | 签发by。 |
| `key_source` | `character varying(32)` | 键来源 | 键来源。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.monitor_rule_state`

- 表中文名：监控规则状态
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `rule_key` | `character varying(64)` | 规则键 | 规则键。 |
| `title` | `character varying(128)` | 标题 | 展示用标题。 |
| `state` | `character varying(16)` | 状态 | 状态。 |
| `severity` | `character varying(16)` | 严重级别 | 严重级别。 |
| `threshold_value` | `double precision` | 阈值值 | 阈值值。 |
| `current_value` | `double precision` | 当前值 | 当前值。 |
| `unit` | `character varying(16)` | 单位 | 数值单位。 |
| `runbook_url` | `character varying(255)` | 操作手册链接 | 操作手册链接。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.node_key_records`

- 表中文名：节点键记录
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `owner_user_id` | `integer` | 所有者用户ID | 关联的所有者用户唯一标识。 |
| `key_material_sha256` | `character varying(128)` | 键materialSHA256 | 键materialSHA256。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `activated_at_ts` | `numeric(18,3)` | activatedat时间戳 | activatedat的时间戳。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `key_source` | `character varying(32)` | 键来源 | 键来源。 |
| `dashboard_url` | `character varying(255)` | 看板链接 | 看板链接。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.node_runtime_status`

- 表中文名：节点运行时status
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `mac_address` | `character varying(64)` | MAC地址 | MAC地址。 |
| `order_id` | `integer` | 订单ID | 关联的订单唯一标识。 |
| `model_version` | `character varying(80)` | 模型版本 | 模型版本。 |
| `npu_usage_pct` | `double precision` | npuusage百分比 | npuusage百分比数值。 |
| `gpu_usage_pct` | `double precision` | gpuusage百分比 | gpuusage百分比数值。 |
| `infer_latency_ms` | `double precision` | 推理延迟毫秒 | 推理延迟毫秒。 |
| `status` | `character varying(30)` | 状态 | 当前记录所处的状态。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.node_workers`

- 表中文名：节点Worker
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `worker_id` | `character varying(80)` | WorkerID | 关联的Worker唯一标识。 |
| `model_id` | `character varying(128)` | 模型ID | 关联的模型唯一标识。 |
| `source_kind` | `character varying(64)` | 来源kind | 来源kind。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `config_json` | `json` | 配置JSON | 结构化配置内容。 |
| `last_error` | `character varying(255)` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.notify_binding_verifications`

- 表中文名：通知绑定验证
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `channel` | `character varying(24)` | 渠道 | 渠道。 |
| `target` | `character varying(255)` | 目标 | 目标。 |
| `code_hash` | `character varying(128)` | 编码哈希 | 编码哈希值。 |
| `code_nonce` | `character varying(64)` | 编码随机数 | 编码随机数。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.notify_dispatch_logs`

- 表中文名：通知分发日志
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `request_id` | `character varying(120)` | 请求ID | 请求级唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `event_type` | `character varying(60)` | 事件type | 事件type。 |
| `channel` | `character varying(24)` | 渠道 | 渠道。 |
| `title` | `character varying(160)` | 标题 | 展示用标题。 |
| `payload_json` | `json` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `response_json` | `json` | 响应JSON | 外部响应内容的 JSON 表示。 |
| `success` | `boolean` | 成功 | 布尔标记，表示成功是否成立。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.order_timelines`

- 表中文名：订单时间线
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `order_id` | `integer` | 订单ID | 关联的订单唯一标识。 |
| `status` | `character varying(30)` | 状态 | 当前记录所处的状态。 |
| `note` | `character varying(255)` | 备注 | 补充说明信息。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.orders`

- 表中文名：订单
- 字段数量：30

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `order_no` | `character varying(64)` | 订单no | 订单no。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `plan_id` | `integer` | 方案ID | 关联的方案唯一标识。 |
| `quantity` | `integer` | quantity | quantity。 |
| `unit_price_usdt` | `numeric(18,6)` | 单位价格usdt | 单位价格usdt。 |
| `discount_rate` | `double precision` | discountrate | discountrate。 |
| `total_usdt` | `numeric(18,6)` | 总量usdt | 总量usdt。 |
| `purchase_mode` | `character varying(20)` | 购买mode | 购买mode。 |
| `group_campaign_id` | `integer` | 分组活动ID | 关联的分组活动唯一标识。 |
| `shipping_unlock_at` | `timestamp without time zone` | 配送解锁时间 | 配送解锁时间。 |
| `payment_chain` | `character varying(20)` | 支付链 | 支付链。 |
| `payment_token` | `character varying(20)` | 支付令牌 | 支付令牌。 |
| `payment_mode` | `character varying(32)` | 支付mode | 支付mode。 |
| `payment_tx_hash` | `character varying(128)` | 支付tx哈希 | 支付tx哈希值。 |
| `payer_wallet_address` | `character varying(128)` | 付款方钱包地址 | 付款方钱包地址。 |
| `payment_intent_at` | `timestamp without time zone` | 支付意图时间 | 支付意图时间。 |
| `payment_status` | `character varying(20)` | 支付状态 | 支付状态。 |
| `status` | `character varying(30)` | 状态 | 当前记录所处的状态。 |
| `coupon_id` | `integer` | 优惠券ID | 关联的优惠券唯一标识。 |
| `sales_agent_id` | `integer` | 销售智能体ID | 关联的销售智能体唯一标识。 |
| `shipping_tracking_no` | `character varying(120)` | 配送trackingno | 配送trackingno。 |
| `shipping_carrier` | `character varying(120)` | 配送承运商 | 配送承运商。 |
| `buyer_note` | `character varying(255)` | 买家备注 | 买家备注。 |
| `shipping_contact_name` | `character varying(80)` | 配送联系人名称 | 配送联系人名称。 |
| `shipping_contact_phone` | `character varying(40)` | 配送联系人手机号 | 配送联系人手机号。 |
| `shipping_region` | `character varying(120)` | 配送地区 | 配送地区。 |
| `shipping_address` | `character varying(255)` | 配送地址 | 配送地址。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.payment_accounts`

- 表中文名：支付账户
- 字段数量：6

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `chain` | `character varying(20)` | 链 | 链。 |
| `token_symbol` | `character varying(20)` | 令牌标的代码 | 令牌标的代码。 |
| `receiving_address` | `character varying(200)` | receiving地址 | receiving地址。 |
| `is_active` | `boolean` | 是否active | 布尔标记，表示是否active。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.product_plans`

- 表中文名：产品方案
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `sku` | `character varying(80)` | sku | sku。 |
| `name` | `character varying(120)` | 名称 | 展示名称。 |
| `unit_price_usdt` | `numeric(18,6)` | 单位价格usdt | 单位价格usdt。 |
| `macmini_count` | `integer` | macmini数量 | macmini数量或次数。 |
| `available_stock` | `integer` | availablestock | availablestock。 |
| `is_active` | `boolean` | 是否active | 布尔标记，表示是否active。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.referral_edges`

- 表中文名：推荐边缘
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `edge_id` | `character varying(96)` | 边缘ID | 关联的边缘唯一标识。 |
| `src_user_id` | `integer` | src用户ID | 关联的src用户唯一标识。 |
| `dst_user_id` | `integer` | dst用户ID | 关联的dst用户唯一标识。 |
| `relation_type` | `character varying(32)` | 关系type | 关系type。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `enabled` | `boolean` | 启用 | 布尔标记，表示启用是否成立。 |
| `effective_from_ts` | `bigint` | 生效来源时间戳 | 生效来源的时间戳。 |
| `effective_to_ts` | `bigint` | 生效目标时间戳 | 生效目标的时间戳。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.runtime_agent_profile`

- 表中文名：运行时智能体配置档案
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `agent_id` | `character varying(64)` | 智能体ID | 关联的智能体唯一标识。 |
| `win_rate` | `double precision` | winrate | winrate。 |
| `pnl_total` | `double precision` | 盈亏总量 | 盈亏总量。 |
| `command_total` | `integer` | 指令总量 | 指令总量。 |
| `last_command` | `character varying(128)` | 最近一次指令 | 最近一次指令。 |
| `bind_status` | `character varying(24)` | 绑定状态 | 绑定状态。 |
| `bind_user_ref` | `character varying(128)` | 绑定用户引用 | 绑定用户引用。 |
| `bind_session_id` | `character varying(64)` | 绑定会话ID | 关联的绑定会话唯一标识。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.runtime_agent_snapshot`

- 表中文名：运行时智能体快照
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `agent_id` | `character varying(64)` | 智能体ID | 关联的智能体唯一标识。 |
| `group_name` | `character varying(64)` | 分组名称 | 分组名称。 |
| `state` | `character varying(32)` | 状态 | 状态。 |
| `strategy_version` | `character varying(32)` | 策略版本 | 策略版本。 |
| `model_version` | `character varying(64)` | 模型版本 | 模型版本。 |
| `queue_length` | `integer` | queuelength | queuelength。 |
| `avg_latency_ms` | `integer` | avg延迟毫秒 | avg延迟毫秒。 |
| `last_error` | `character varying(255)` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `token_total` | `double precision` | 令牌总量 | 令牌总量。 |
| `snapshot_ts` | `timestamp with time zone` | 快照时间戳 | 快照的时间戳。 |
| `report_mode` | `character varying(16)` | 报告mode | 报告mode。 |

### 表：`public.runtime_kv`

- 表中文名：运行时kv
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `key` | `character varying(120)` | 键 | 键。 |
| `value` | `character varying(512)` | 值 | 通用值字段。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.runtime_lobster_binding_session`

- 表中文名：运行时Lobster绑定会话
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `binding_session_id` | `character varying(64)` | 绑定会话ID | 关联的绑定会话唯一标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `agent_id` | `character varying(64)` | 智能体ID | 关联的智能体唯一标识。 |
| `claim_session_id` | `character varying(64)` | 认领会话ID | 关联的认领会话唯一标识。 |
| `user_ref` | `character varying(128)` | 用户引用 | 用户引用。 |
| `binding_key_hash` | `character varying(128)` | 绑定键哈希 | 绑定键哈希值。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `expire_at` | `timestamp with time zone` | 过期时间 | 记录到期时间。 |
| `verified_at` | `timestamp with time zone` | 已验证时间 | 已验证时间。 |
| `bound_at` | `timestamp with time zone` | 已绑定时间 | 已绑定时间。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.runtime_lobster_claim_session`

- 表中文名：运行时Lobster认领会话
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `claim_session_id` | `character varying(64)` | 认领会话ID | 关联的认领会话唯一标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `agent_id` | `character varying(64)` | 智能体ID | 关联的智能体唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `qr_content` | `character varying(255)` | qr内容 | qr内容。 |
| `expire_at` | `timestamp with time zone` | 过期时间 | 记录到期时间。 |
| `scanned_at` | `timestamp with time zone` | 扫描时间 | 扫描时间。 |
| `created_at` | `timestamp with time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.runtime_subsystem_connection`

- 表中文名：运行时子系统connection
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `node_id` | `character varying(64)` | 节点ID | 关联的节点唯一标识。 |
| `system_key` | `character varying(64)` | 系统键 | 系统键。 |
| `display_name` | `character varying(128)` | 显示名称 | 显示名称。 |
| `protocol` | `character varying(16)` | 协议 | 协议。 |
| `base_url` | `character varying(255)` | base链接 | base链接。 |
| `dashboard_url` | `character varying(255)` | 看板链接 | 看板链接。 |
| `capabilities_json` | `text` | 能力JSON | 能力的结构化 JSON 数据。 |
| `status` | `character varying(16)` | 状态 | 当前记录所处的状态。 |
| `health_score` | `double precision` | 健康度评分 | 健康度评分。 |
| `last_error` | `character varying(255)` | 最近一次错误 | 最近一次处理失败时的错误信息。 |
| `last_checked_at` | `timestamp with time zone` | 最近一次检查时间 | 最近一次检查时间。 |
| `updated_at` | `timestamp with time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.sales_agents`

- 表中文名：销售智能体
- 字段数量：6

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `name` | `character varying(80)` | 名称 | 展示名称。 |
| `email` | `character varying(255)` | 邮箱 | 邮箱地址。 |
| `code_prefix` | `character varying(40)` | 编码prefix | 编码prefix。 |
| `is_active` | `boolean` | 是否active | 布尔标记，表示是否active。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

### 表：`public.schema_migrations`

- 表中文名：结构迁移
- 字段数量：4

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `version` | `character varying(160)` | version | version。 |
| `filename` | `character varying(260)` | filename | filename。 |
| `checksum` | `character varying(128)` | 校验和 | 校验和。 |
| `applied_at` | `character varying(40)` | applied时间 | applied时间。 |

### 表：`public.settlement_subsidy_day_states`

- 表中文名：结算补贴day状态
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `trade_date` | `character varying(20)` | 交易日期 | 交易日期。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `used_budget_dec` | `character varying(64)` | usedbudget十进制定点值 | usedbudget的十进制定点数值。 |
| `shadow_budget_accrued_dec` | `character varying(64)` | 影子budgetaccrued十进制定点值 | 影子budgetaccrued的十进制定点数值。 |
| `last_shortfall_dec` | `character varying(64)` | 最近一次缺口十进制定点值 | 最近一次缺口的十进制定点数值。 |
| `last_event_id` | `character varying(128)` | 最近一次事件ID | 关联的最近一次事件唯一标识。 |
| `policy_version` | `character varying(32)` | 策略版本 | 策略版本。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.settlement_subsidy_items`

- 表中文名：结算补贴条目
- 字段数量：18

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `allocation_id` | `character varying(96)` | 分配ID | 关联的分配唯一标识。 |
| `source_event_id` | `character varying(128)` | 来源事件ID | 关联的来源事件唯一标识。 |
| `trade_date` | `character varying(20)` | 交易日期 | 交易日期。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `target_pnl_dec` | `character varying(64)` | 目标盈亏十进制定点值 | 目标盈亏的十进制定点数值。 |
| `trade_net_pnl_dec` | `character varying(64)` | 交易净盈亏十进制定点值 | 交易净盈亏的十进制定点数值。 |
| `paid_before_dec` | `character varying(64)` | 已支付之前十进制定点值 | 已支付之前的十进制定点数值。 |
| `gap_dec` | `character varying(64)` | 差额十进制定点值 | 差额的十进制定点数值。 |
| `subsidy_dec` | `character varying(64)` | 补贴十进制定点值 | 补贴的十进制定点数值。 |
| `weight_dec` | `character varying(64)` | 权重十进制定点值 | 权重的十进制定点数值。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `policy_version` | `character varying(32)` | 策略版本 | 策略版本。 |
| `memo` | `character varying(255)` | 备注 | 备注。 |
| `payload_json` | `json` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.siwe_nonces`

- 表中文名：SIWE随机数
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `nonce` | `character varying(80)` | 随机数 | 随机数。 |
| `wallet_address` | `character varying(96)` | 钱包地址 | 钱包地址。 |
| `claim_session_id` | `character varying(80)` | 认领会话ID | 关联的认领会话唯一标识。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `consumed` | `boolean` | 已消耗 | 布尔标记，表示已消耗是否成立。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.unified_key_registry`

- 表中文名：统一键注册表
- 字段数量：12

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `key_ref` | `character varying(120)` | 键引用 | 键引用。 |
| `key_type` | `character varying(40)` | 键type | 键type。 |
| `key_source` | `character varying(40)` | 键来源 | 键来源。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `owner_user_id` | `integer` | 所有者用户ID | 关联的所有者用户唯一标识。 |
| `follower_user_id` | `integer` | 跟单用户用户ID | 关联的跟单用户用户唯一标识。 |
| `leader_node_id` | `character varying(96)` | 带单节点节点ID | 关联的带单节点节点唯一标识。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.user_auth_identities`

- 表中文名：用户auth身份
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `identity_type` | `character varying(16)` | 身份type | 身份type。 |
| `identity_value` | `character varying(255)` | 身份值 | 身份值。 |
| `verified` | `boolean` | 已验证 | 布尔标记，表示已验证是否成立。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.user_lobster_binding_requests`

- 表中文名：用户Lobster绑定请求
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `request_no` | `character varying(80)` | 请求no | 请求no。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `node_id` | `character varying(96)` | 节点ID | 关联的节点唯一标识。 |
| `follower_key_id` | `character varying(80)` | 跟单用户键ID | 关联的跟单用户键唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `lobster_binding_ref` | `character varying(120)` | Lobster绑定引用 | Lobster绑定引用。 |
| `trace_id` | `character varying(80)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.user_lobster_profiles`

- 表中文名：用户Lobster配置档案
- 字段数量：7

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `slot_cap` | `integer` | 槽位上限 | 槽位上限。 |
| `tier` | `character varying(24)` | 等级 | 等级。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.user_notify_bindings`

- 表中文名：用户通知绑定
- 字段数量：8

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `channel` | `character varying(24)` | 渠道 | 渠道。 |
| `endpoint` | `character varying(255)` | 端点 | 通知或回调目标端点。 |
| `endpoint_meta_json` | `json` | 端点元数据JSON | 端点附加信息的 JSON 表示。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.user_sessions`

- 表中文名：用户会话
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `session_token` | `character varying(256)` | 会话令牌 | 会话令牌。 |
| `expire_at_ts` | `numeric(18,3)` | 过期at时间戳 | 过期at的时间戳。 |
| `ip_address` | `character varying(64)` | IP地址 | 来源 IP 地址。 |
| `user_agent` | `character varying(255)` | 用户代理 | 客户端 User-Agent。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |
| `token_hash` | `character varying(128)` | 令牌哈希 | 令牌哈希值。 |
| `is_admin` | `boolean` | 是否管理员 | 布尔标记，表示是否管理员。 |
| `expire_at` | `timestamp without time zone` | 过期时间 | 记录到期时间。 |

### 表：`public.users`

- 表中文名：用户
- 字段数量：9

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `external_user_id` | `character varying(64)` | 外部用户ID | 外部系统中的用户标识。 |
| `display_name` | `character varying(120)` | 显示名称 | 显示名称。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `custody_wallet_id` | `integer` | 托管钱包ID | 关联的托管钱包唯一标识。 |
| `preferred_channel` | `character varying(32)` | 偏好渠道 | 偏好渠道。 |
| `locale` | `character varying(16)` | 区域 | 区域。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.wallet_aggregation_profiles`

- 表中文名：钱包归集配置档案
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `profile_id` | `character varying(96)` | 配置档案ID | 关联的配置档案唯一标识。 |
| `settlement_pool_id` | `character varying(96)` | 结算池ID | 关联的结算池唯一标识。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `strategy` | `character varying(32)` | 策略 | 策略。 |
| `util_ratio_dec` | `character varying(64)` | 利用率ratio十进制定点值 | 利用率ratio的十进制定点数值。 |
| `cap_per_binding_dec` | `character varying(64)` | 上限每绑定十进制定点值 | 上限每绑定的十进制定点数值。 |
| `risk_cap_dec` | `character varying(64)` | 风控上限十进制定点值 | 风控上限的十进制定点数值。 |
| `meta_json` | `json` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.withdrawal_requests`

- 表中文名：提现请求
- 字段数量：15

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `withdrawal_id` | `character varying(80)` | 提现ID | 关联的提现唯一标识。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `idempotency_key` | `character varying(120)` | 幂等键 | 用于避免重复处理的幂等标识。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `asset_code` | `character varying(24)` | 资产编码 | 资产编码。 |
| `amount_dec` | `character varying(64)` | 金额十进制定点值 | 金额的十进制定点数值。 |
| `amount_atomic` | `bigint` | 金额原子单位值 | 金额的原子单位数值。 |
| `to_address` | `character varying(96)` | 目标地址 | 目标地址。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `risk_snapshot_json` | `json` | 风控快照JSON | 风控快照的结构化 JSON 数据。 |
| `gas_snapshot_json` | `json` | Gas快照JSON | Gas快照的结构化 JSON 数据。 |
| `reviewed_by` | `character varying(80)` | 审核by | 审核by。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.withdrawal_transactions`

- 表中文名：提现交易
- 字段数量：11

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `withdrawal_id` | `character varying(80)` | 提现ID | 关联的提现唯一标识。 |
| `tx_hash` | `character varying(120)` | tx哈希 | tx哈希值。 |
| `chain_id` | `integer` | 链ID | 关联的链唯一标识。 |
| `nonce` | `integer` | 随机数 | 随机数。 |
| `max_fee_per_gas` | `character varying(64)` | 最大fee每Gas | 最大fee每Gas。 |
| `max_priority_fee_per_gas` | `character varying(64)` | 最大优先fee每Gas | 最大优先fee每Gas。 |
| `confirmations` | `integer` | 确认数 | 确认数。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.worker_checkpoints`

- 表中文名：Worker检查点
- 字段数量：5

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `worker_name` | `character varying(40)` | Worker名称 | Worker名称。 |
| `chain` | `character varying(20)` | 链 | 链。 |
| `last_scanned_block` | `integer` | 最近一次扫描区块 | 最近一次扫描区块。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.worker_products`

- 表中文名：Worker产品
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `worker_code` | `character varying(80)` | Worker编码 | Worker编码。 |
| `name` | `character varying(120)` | 名称 | 展示名称。 |
| `slug` | `character varying(120)` | 标识 | 标识。 |
| `category` | `character varying(64)` | 分类 | 所属分类。 |
| `summary` | `character varying(255)` | 摘要 | 摘要。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `vendor_code` | `character varying(64)` | 供应商编码 | 供应商编码。 |
| `default_price_usdt` | `numeric(18,6)` | 默认价格usdt | 默认价格usdt。 |
| `sort_order` | `integer` | 排序订单 | 排序订单。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.worker_release_assets`

- 表中文名：Worker发布资产
- 字段数量：10

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `release_id` | `integer` | 发布ID | 关联的发布唯一标识。 |
| `asset_type` | `character varying(24)` | 资产type | 资产type。 |
| `asset_name` | `character varying(120)` | 资产名称 | 资产名称。 |
| `asset_url` | `character varying(255)` | 资产链接 | 资产链接。 |
| `checksum_sha256` | `character varying(128)` | 校验和SHA256 | 校验和SHA256。 |
| `size_bytes` | `integer` | sizebytes | sizebytes。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.worker_releases`

- 表中文名：Worker发布
- 字段数量：16

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `worker_product_id` | `integer` | Worker产品ID | 关联的Worker产品唯一标识。 |
| `release_version` | `character varying(40)` | 发布版本 | 发布版本。 |
| `release_channel` | `character varying(24)` | 发布渠道 | 发布渠道。 |
| `target_os` | `character varying(24)` | 目标操作系统 | 目标操作系统。 |
| `target_arch` | `character varying(24)` | 目标架构 | 目标架构。 |
| `package_url` | `character varying(255)` | 包链接 | 包链接。 |
| `checksum_sha256` | `character varying(128)` | 校验和SHA256 | 校验和SHA256。 |
| `signature_ref` | `character varying(255)` | 签名引用 | 签名引用。 |
| `compatibility_note` | `character varying(255)` | 兼容说明备注 | 兼容说明备注。 |
| `release_note` | `text` | 发布备注 | 发布备注。 |
| `status` | `character varying(24)` | 状态 | 当前记录所处的状态。 |
| `published_at` | `timestamp without time zone` | 发布时间 | 发布时间。 |
| `meta_json` | `text` | 元数据JSON | 附加元数据的 JSON 表示。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |
| `updated_at` | `timestamp without time zone` | 更新时间 | 记录最后更新时间。 |

### 表：`public.worker_runtime_logs`

- 表中文名：Worker运行时日志
- 字段数量：13

| 字段名 | 字段类型 | 字段中文名 | 字段意思 |
| --- | --- | --- | --- |
| `id` | `integer` | 主键ID | 当前表的主键标识。 |
| `source_node_id` | `character varying(96)` | 来源节点ID | 关联的来源节点唯一标识。 |
| `worker_code` | `character varying(80)` | Worker编码 | Worker编码。 |
| `worker_version` | `character varying(40)` | Worker版本 | Worker版本。 |
| `level` | `character varying(16)` | 级别 | 分级或严重程度级别。 |
| `event_type` | `character varying(40)` | 事件type | 事件type。 |
| `message` | `character varying(255)` | 消息 | 消息。 |
| `order_no` | `character varying(64)` | 订单no | 订单no。 |
| `user_id` | `integer` | 用户ID | 关联的用户唯一标识。 |
| `trace_id` | `character varying(80)` | 链路追踪ID | 用于串联跨服务调用链路的追踪标识。 |
| `payload_json` | `text` | 载荷JSON | 请求或事件的原始载荷 JSON。 |
| `reported_at` | `timestamp without time zone` | 上报时间 | 上报时间。 |
| `created_at` | `timestamp without time zone` | 创建时间 | 记录创建时间。 |

