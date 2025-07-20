-- 用户表
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(20) UNIQUE NOT NULL,
    "nickname" VARCHAR(60) UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR(64) UNIQUE NOT NULL,
    "gender" SMALLINT NOT NULL DEFAULT 3,
    "avatar_url" VARCHAR(255) NOT NULL DEFAULT '',
    "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "users"."id" IS '用户ID';

COMMENT ON COLUMN "users"."username" IS '用户名';

COMMENT ON COLUMN "users"."nickname" IS '用户昵称';

COMMENT ON COLUMN "users"."password" IS '用户密码';

COMMENT ON COLUMN "users"."email" IS '用户邮箱';

COMMENT ON COLUMN "users"."gender" IS '用户性别, 1:男, 2:女, 3: 未知';

COMMENT ON COLUMN "users"."avatar_url" IS '头像URL';

COMMENT ON COLUMN "users"."password_changed_at" IS '上次密码更新时间';

COMMENT ON COLUMN "users"."created_at" IS '创建时间';

COMMENT ON COLUMN "users"."updated_at" IS '更新时间';