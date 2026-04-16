CREATE TABLE IF NOT EXISTS streak (
                          "id" SERIAL PRIMARY KEY,
                          "user_habit_id" INTEGER REFERENCES user_habit(id),
                          "current_streak" INTEGER,
                          "longest_streak" INTEGER,
                          "last_confirmed_date" TIMESTAMP
);

CREATE TABLE daily_confirmation (
                                      "id" SERIAL PRIMARY KEY,
                                      "user_habit_id" INTEGER REFERENCES user_habit(id) ON DELETE CASCADE,
                                      "confirmed_at" TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_user_habit_day
    ON daily_confirmation (user_habit_id, DATE(confirmed_at));

CREATE UNIQUE INDEX unique_streak_user_habit
    ON streak (user_habit_id);