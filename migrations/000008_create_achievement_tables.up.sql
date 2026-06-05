CREATE TABLE IF NOT EXISTS metric (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      metric_key VARCHAR(60) UNIQUE NOT NULL
);

INSERT INTO metric (metric_key)
VALUES ('total_habits'), ('current_streak')
ON CONFLICT (metric_key) DO NOTHING;

CREATE TABLE IF NOT EXISTS user_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    metric_id UUID NOT NULL REFERENCES metric(id),
    value BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, metric_id)
);

CREATE TABLE IF NOT EXISTS user_habit_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_habit_id INTEGER NOT NULL REFERENCES user_habit(id) ON DELETE CASCADE,
    metric_id UUID NOT NULL REFERENCES metric(id),
    value BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_habit_id, metric_id)
);

CREATE INDEX IF NOT EXISTS user_metrics_user_id_idx
    ON user_metrics(user_id);

CREATE TABLE IF NOT EXISTS achievement (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   code VARCHAR(100) UNIQUE NOT NULL,
   title VARCHAR(255) NOT NULL,
   description TEXT,
   enabled BOOLEAN NOT NULL DEFAULT TRUE,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS achievement_condition (
 id BIGSERIAL PRIMARY KEY,
 achievement_id UUID NOT NULL REFERENCES achievement(id) ON DELETE CASCADE,
 metric_scope VARCHAR(10) NOT NULL,
 metric_id UUID NOT NULL REFERENCES metric(id),
 operator VARCHAR(10) NOT NULL,
 value BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS achievement_conditions_metric_id_idx
    ON achievement_condition(metric_id);

CREATE INDEX IF NOT EXISTS achievement_conditions_achievement_id_idx
    ON achievement_condition(achievement_id);

CREATE TABLE IF NOT EXISTS user_achievement (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievement(id),
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS user_achievements_user_id_idx
    ON user_achievement(user_id);
