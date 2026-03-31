CREATE TYPE account_role AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS "users" (
     "id" SERIAL PRIMARY KEY,
     "username" VARCHAR(255),
     "email" VARCHAR(255),
     "password" VARCHAR(255),
     "role" account_role NOT NULL DEFAULT 'user',
     "created_at" TIMESTAMP DEFAULT now(),
     "is_active" BOOLEAN default true
);
