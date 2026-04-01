# 小龙虾用户 Gorm CRUD 落地

## Goal
基于 `docs/operator-follower-user-custom-table-crud.md` 中定义的 `operator_portal.follower_users` 独立表，使用 Gorm 在 `backend/` 中实现完整的运营后台 CRUD 接口，并接入现有统一响应组件。

## Requirements
- 使用 Gorm 为 `operator_portal.follower_users` 建立模型映射。
- 提供小龙虾用户的创建、分页查询、详情查询、更新、删除接口。
- 接口路径与文档保持一致，统一挂载在 `/admin/v1/follower-users`。
- 请求与响应字段命名统一使用小驼峰，时间字符串统一使用 `yyyy-MM-dd HH:mm:ss`。
- 所有接口必须复用现有统一响应组件。
- 删除采用逻辑删除，更新与删除需要校验 `rowVersion`。
- 列表查询需要支持文档中定义的筛选条件和分页结构。
- 代码实现需保留清晰的分层，至少包含模型、仓储、服务、HTTP 处理器与路由注册。

## Acceptance Criteria
- [ ] `backend/` 中存在 `follower_users` 的 Gorm 模型和 CRUD 业务代码。
- [ ] `POST /admin/v1/follower-users` 可以创建记录并返回文档约定字段。
- [ ] `GET /admin/v1/follower-users` 可以按条件分页查询并返回统一分页结构。
- [ ] `GET /admin/v1/follower-users/{userId}` 可以返回详情。
- [ ] `PATCH /admin/v1/follower-users/{userId}` 可以基于 `rowVersion` 更新数据。
- [ ] `DELETE /admin/v1/follower-users/{userId}` 可以执行逻辑删除并返回成功结果。
- [ ] `go test ./...` 与 `go build ./...` 通过。

## Technical Notes
- 数据库表位于 `operator_portal` schema。
- 当前代码库已接入 Gin、Gorm、Viper、Zap 和统一响应组件，可直接复用。
- 由于文档已明确 API 契约，本次实现按“文档即契约”落地，不再额外设计第二套字段结构。
