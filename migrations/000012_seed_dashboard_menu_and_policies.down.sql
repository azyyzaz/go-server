DELETE FROM role_menus
WHERE menu_id IN (
    SELECT id FROM menus WHERE path = '/dashboard'
);

DELETE FROM casbin_rule
WHERE ptype = 'p'
  AND v0 = 'admin'
  AND v1 = '/api/v1/dashboard/*'
  AND v2 = 'GET';

DELETE FROM menus
WHERE path = '/dashboard';
