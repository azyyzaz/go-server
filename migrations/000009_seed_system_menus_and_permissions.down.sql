DELETE FROM casbin_rule
WHERE ptype = 'p'
  AND v0 = 'admin'
  AND v1 IN (
      '/api/v1/system/menus',
      '/api/v1/system/menus/*',
      '/api/v1/system/depts',
      '/api/v1/system/depts/*',
      '/api/v1/system/dicts/types',
      '/api/v1/system/dicts/items',
      '/api/v1/system/dicts/types/*',
      '/api/v1/system/dicts/items/*',
      '/api/v1/system/dicts/lookup/*'
  );

DELETE rm FROM role_menus rm
JOIN roles r ON r.id = rm.role_id
JOIN menus m ON m.id = rm.menu_id
WHERE r.code = 'admin'
  AND m.path IN ('/system', '/system/users', '/system/roles', '/system/menus', '/system/depts', '/system/dicts');

DELETE FROM menus
WHERE path IN ('/system', '/system/users', '/system/roles', '/system/menus', '/system/depts', '/system/dicts');
