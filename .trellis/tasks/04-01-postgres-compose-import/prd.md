# 创建 PostgreSQL Docker Compose 导入任务

## 目标
创建一套本地 Docker Compose 配置，用于启动 PostgreSQL 16，并在首次初始化时自动导入 `/Users/leon/Downloads/schema_pincermarket_20260401.sql` 这个 SQL dump。

## 需求
- 通过 Docker Compose 运行 PostgreSQL 16。
- 为导入的 dump 创建默认数据库。
- 在容器首次初始化时自动导入提供的 SQL dump。
- 正确处理 dump 的 UTF-16 little-endian 编码，确保 `psql` 可以稳定执行。
- 使用持久化 Docker volume 保存数据库数据。

## 验收标准
- [ ] `docker-compose.yml` 中定义了可在本地启动的 PostgreSQL 服务。
- [ ] SQL dump 会挂载进容器，并在首次启动时自动导入。
- [ ] 导入流程可以正确处理当前的 UTF-16 dump 文件。
- [ ] 使用说明清楚描述启动方式，以及重置 / 重新导入的行为。

## 技术说明
- 提供的 dump 由 PostgreSQL 16 生成。
- dump 文件编码为 UTF-16 LE，且使用 CRLF 换行。
- 初步检查显示，这份 dump 主要是以 schema 为主的导出，包含 DDL 和约束，没有发现 `COPY` 或 `INSERT` 语句。
