INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT NULL, '系统管理', 'directory', '/system', 'Layout', '', 100, 1, 1
WHERE NOT EXISTS (
    SELECT 1 FROM menus WHERE name = '系统管理' AND parent_id IS NULL
);

INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '用户管理', 'menu', '/system/users', 'system/users/index', 'system:user:list', 10, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '用户管理' AND path = '/system/users'
  );

INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '角色管理', 'menu', '/system/roles', 'system/roles/index', 'system:role:list', 20, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '角色管理' AND path = '/system/roles'
  );

INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '菜单管理', 'menu', '/system/menus', 'system/menus/index', 'system:menu:list', 30, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '菜单管理' AND path = '/system/menus'
  );

INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '部门管理', 'menu', '/system/depts', 'system/depts/index', 'system:dept:list', 40, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '部门管理' AND path = '/system/depts'
  );

INSERT INTO menus (parent_id, name, type, path, component, permission, sort, visible, status)
SELECT root.id, '字典管理', 'menu', '/system/dicts', 'system/dicts/index', 'system:dict:list', 50, 1, 1
FROM menus root
WHERE root.name = '系统管理'
  AND root.parent_id IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM menus WHERE name = '字典管理' AND path = '/system/dicts'
  );

INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r
JOIN menus m ON m.path IN ('/system', '/system/users', '/system/roles', '/system/menus', '/system/depts', '/system/dicts')
WHERE r.code = 'admin'
  AND NOT EXISTS (
      SELECT 1 FROM role_menus rm WHERE rm.role_id = r.id AND rm.menu_id = m.id
  );

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/menus', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/menus' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/menus', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/menus' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/menus/*', '*'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/menus/*' AND v2 = '*'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/depts', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/depts' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/depts', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/depts' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/depts/*', '*'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/depts/*' AND v2 = '*'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/types', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/types' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/types', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/types' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/items', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/items' AND v2 = 'GET'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/items', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/items' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/types/*', '*'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/types/*' AND v2 = '*'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/items/*', '*'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/items/*' AND v2 = '*'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'admin', '/api/v1/system/dicts/lookup/*', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 = '/api/v1/system/dicts/lookup/*' AND v2 = 'GET'
);
