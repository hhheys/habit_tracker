ALTER TABLE users ADD COLUMN timezone TEXT;
UPDATE users SET timezone = 'Europe/Moscow' WHERE timezone IS NULL;