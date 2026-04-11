DELETE FROM casbin_rule WHERE ptype = 'g' AND v0 = 'admin' AND v1 = 'admin';
DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 LIKE '/api/v1/system/users%';
DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 = 'admin' AND v1 LIKE '/api/v1/system/roles%';
DELETE FROM user_roles
WHERE user_id = (SELECT id FROM users WHERE username = 'admin')
  AND role_id = (SELECT id FROM roles WHERE code = 'admin');
DELETE FROM roles WHERE code = 'admin';
