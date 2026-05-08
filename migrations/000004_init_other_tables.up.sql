CREATE TYPE "goal_status" AS ENUM (
    'active',
    'declined',
    'passed'
    );


CREATE TABLE "goal" (
                        "id" serial PRIMARY KEY,
                        "user_habit_id" int,
                        "description" text,
                        "start_date" date,
                        "end_date" date,
                        "status" goal_status
);

CREATE TABLE "post" (
                        "id" serial PRIMARY KEY,
                        "user_habit_id" int,
                        "title" text,
                        "description" text,
                        "created_at" timestamp
);

CREATE TABLE "comment" (
                           "id" serial PRIMARY KEY,
                           "post_id" int,
                           "user_id" int,
                           "content" text,
                           "created_at" timestamp
);

CREATE TABLE "reaction" (
                            "id" serial PRIMARY KEY,
                            "post_id" int,
                            "user_id" int,
                            "timestamp" timestamp
);

CREATE TABLE "achievement" (
                               "id" serial PRIMARY KEY,
                               "title" varchar,
                               "description" text,
                               "condition" varchar
);

CREATE TABLE "user_achievement" (
                                    "id" serial PRIMARY KEY,
                                    "user_id" int,
                                    "achievement_id" int,
                                    "earned_at" timestamp
);

CREATE TABLE "reminder_type" (
                                 "id" serial PRIMARY KEY,
                                 "name" varchar,
                                 "description" varchar
);

CREATE TABLE "reminder" (
                            "id" serial PRIMARY KEY,
                            "user_habit_id" int,
                            "reminder_type" int,
                            "reminder_time" time,
                            "created_at" timestamp,
                            "is_active" boolean
);

CREATE TABLE "user_settings" (
                                 "id" serial PRIMARY KEY,
                                 "user_id" int,
                                 "notifications_enabled" boolean,
                                 "dark_theme" boolean,
                                 "language" varchar
);

ALTER TABLE "goal" ADD FOREIGN KEY ("user_habit_id") REFERENCES "user_habit" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "post" ADD FOREIGN KEY ("user_habit_id") REFERENCES "user_habit" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "reaction" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "reaction" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_achievement" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_achievement" ADD FOREIGN KEY ("achievement_id") REFERENCES "achievement" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "reminder" ADD FOREIGN KEY ("id") REFERENCES "user_habit" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "reminder" ADD FOREIGN KEY ("reminder_type") REFERENCES "reminder_type" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comment" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comment" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;