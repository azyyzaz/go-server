INSERT INTO roles (name, code, remark, status)
SELECT '管理员', 'admin', '系统默认管理员', 1
WHERE NOT EXISTS (
    SELECT 1 FROM roles WHERE code = 'admin'
);

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin'
  AND r.code = 'admin'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id
  );

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'g', 'admin', 'admin', ''
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'g' AND v0 = 'admin' AND v1 = 'admin'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users/*', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users/*' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users/*', 'PUT'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users/*' AND v2 = 'PUT'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users/*', 'DELETE'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users/*' AND v2 = 'DELETE'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/users/*', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/users/*' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/roles', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/roles' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/roles', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/roles' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/roles/*', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/roles/*' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/roles/*', 'PUT'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/roles/*' AND v2 = 'PUT'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/roles/*', 'DELETE'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/roles/*' AND v2 = 'DELETE'
);
