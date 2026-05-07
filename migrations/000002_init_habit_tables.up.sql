CREATE TABLE IF NOT EXISTS tag (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255)
);

CREATE TABLE habit (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(255) UNIQUE,
                         description TEXT,
                         created_at TIMESTAMP DEFAULT NOW(),
                          image_filename VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS habit_tag (
                             id SERIAL PRIMARY KEY,
                             tag_id INTEGER REFERENCES tag(id) ON DELETE CASCADE,
                             habit_id INTEGER REFERENCES habit(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_habit (
                              id SERIAL PRIMARY KEY,
                              "habit_id" INTEGER REFERENCES habit(id) ON DELETE CASCADE,
                              "user_id" INTEGER REFERENCES users(id),
                              "added_at" TIMESTAMP DEFAULT NOW()
);

ALTER TABLE user_habit ADD CONSTRAINT unique_user_habit UNIQUE ("habit_id", "user_id");

ALTER TABLE habit_tag
    ADD CONSTRAINT habit_tag_unique UNIQUE (habit_id, tag_id);