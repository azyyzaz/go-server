-- 默认管理员账号：admin / admin123
INSERT INTO users (username, password, name, email, status)
VALUES (
    'admin',
    '$2a$10$p0sKb71jHiKyj3b6gWL6deJMf5cNmZh71Twe0Dx/Ss231t6JAjG1W',
    'Admin',
    'admin@example.com',
    1
);
