-- +goose Up
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT TRUE NOT NULL;
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP NULL;

ALTER TABLE users MODIFY COLUMN role ENUM('owner', 'admin', 'student') NOT NULL DEFAULT 'student';

-- +goose Down
ALTER TABLE users DROP COLUMN is_active;
ALTER TABLE users DROP COLUMN last_login_at;

UPDATE users SET role = 'admin' WHERE role = 'owner';
ALTER TABLE users MODIFY COLUMN role ENUM('admin', 'student') NOT NULL DEFAULT 'student';