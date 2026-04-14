DELETE FROM casbin_rule
WHERE ptype = 'p'
  AND v0 = 'admin'
  AND (
      (v1 = '/api/v1/system/files' AND v2 = 'GET')
      OR (v1 = '/api/v1/system/files/upload' AND v2 = 'POST')
      OR (v1 = '/api/v1/system/files/*' AND v2 = 'DELETE')
  );

DELETE rm
FROM role_menus rm
JOIN menus m ON m.id = rm.menu_id
WHERE m.path = '/system/files';

DELETE FROM menus
WHERE path = '/system/files';
