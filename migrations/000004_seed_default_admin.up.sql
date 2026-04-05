-- 默认管理员账号：admin / admin123
INSERT INTO users (username, password, name, email, status)
VALUES (
    'admin',
    '$2a$10$p0sKb71jHiKyj3b6gWL6deJMf5cNmZh71Twe0Dx/Ss231t6JAjG1W',
    'Admin',
    'admin@example.com',
    1
);

-- 为 admin 用户分配管理员角色
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin' AND r.code = 'admin';
