ALTER TABLE users
    DROP FOREIGN KEY fk_users_dept,
    DROP INDEX idx_users_dept_id,
    DROP COLUMN dept_id;

DROP TABLE IF EXISTS depts;
