INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT NULL, '仪表盘', 'menu', '/dashboard', 'dashboard/index', 'dashboard:view', 1, 1, 1
WHERE NOT EXISTS (
    SELECT 1 FROM menus WHERE name = '仪表盘' AND path = '/dashboard'
);

INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r
JOIN menus m ON m.path = '/dashboard'
WHERE r.code = 'admin'
  AND NOT EXISTS (
      SELECT 1 FROM role_menus rm WHERE rm.role_id = r.id AND rm.menu_id = m.id
  );

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/dashboard/*', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/dashboard/*' AND v2 = 'GET'
);
