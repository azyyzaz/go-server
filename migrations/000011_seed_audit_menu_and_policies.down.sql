DELETE FROM casbin_rule
WHERE ptype = 'p'
  AND v0 = 'admin'
  AND v1 IN (
      '/api/v1/system/audits/operation-logs',
      '/api/v1/system/audits/login-logs'
  );

DELETE rm FROM role_menus rm
JOIN roles r ON r.id = rm.role_id
JOIN menus m ON m.id = rm.menu_id
WHERE r.code = 'admin'
  AND m.path = '/system/audits';

DELETE FROM menus
WHERE path = '/system/audits';
