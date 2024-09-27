-- 用户表
CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(20) UNIQUE NOT NULL,
    "nickname" VARCHAR(60) UNIQUE NOT NULL,
    "password" TEXT NOT NULL,
    "email" VARCHAR(64) UNIQUE NOT NULL,
    "gender" SMALLINT NOT NULL DEFAULT 1,
    "profile_picture_url" VARCHAR(255) NOT NULL DEFAULT 'default.jpg',
    "is_online" BOOLEAN NOT NULL DEFAULT FALSE,
    "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
)