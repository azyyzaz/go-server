INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '文件管理', 'menu', '/system/files', 'system/files/index', 'system:file:list', 70, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '文件管理' AND path = '/system/files'
  );

INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r
JOIN menus m ON m.path = '/system/files'
WHERE r.code = 'admin'
  AND NOT EXISTS (
      SELECT 1 FROM role_menus rm WHERE rm.role_id = r.id AND rm.menu_id = m.id
  );

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/files', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/files' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/files/upload', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/files/upload' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/files/*', 'DELETE'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/files/*' AND v2 = 'DELETE'
);
