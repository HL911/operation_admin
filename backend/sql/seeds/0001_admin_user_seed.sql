-- 0001_admin_user_seed.sql
-- 用途：初始化首个后台开发者管理员账号。
-- 执行前请先确认 operator_portal schema、admin_user_seq 序列和 admin_users 表已经由服务启动完成初始化。
-- 注意：请先替换下方的 login_name 与 password_hash 占位值，再执行本脚本。

INSERT INTO operator_portal.admin_users (
    admin_user_id,
    login_name,
    password_hash,
    display_name,
    role_code,
    status,
    created_at,
    updated_at,
    is_deleted,
    row_version
) VALUES (
    'ADM-' || nextval('operator_portal.admin_user_seq'),
    '{{REPLACE_WITH_LOGIN_NAME}}',
    '{{REPLACE_WITH_BCRYPT_HASH}}',
    '系统管理员',
    'admin',
    'active',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    false,
    1
);

-- 示例提醒：
-- 1. login_name 必须满足 ^[a-z0-9._-]{4,32}$ 且全局唯一。
-- 2. password_hash 必须是 bcrypt 生成的哈希值，禁止写入明文密码。
