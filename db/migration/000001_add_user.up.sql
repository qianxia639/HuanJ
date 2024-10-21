-- 用户表
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(20) UNIQUE NOT NULL,
    "nickname" VARCHAR(60) UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "salt" VARCHAR NOT NULL,
    "email" VARCHAR(64) UNIQUE NOT NULL,
    "gender" SMALLINT NOT NULL DEFAULT 3,
    "avatar" VARCHAR NOT NULL DEFAULT 'default.jpg',
    "is_online" BOOLEAN NOT NULL DEFAULT FALSE,
    "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "users"."id" IS '用户ID';

COMMENT ON COLUMN "users"."username" IS '用户名';

COMMENT ON COLUMN "users"."nickname" IS '用户昵称';

COMMENT ON COLUMN "users"."password" IS '用户密码';

COMMENT ON COLUMN "users"."salt" IS '随机盐';

COMMENT ON COLUMN "users"."email" IS '用户邮箱';

COMMENT ON COLUMN "users"."gender" IS '用户性别, 1:男, 2:女, 3: 未知';

COMMENT ON COLUMN "users"."avatar" IS '用户头像';

COMMENT ON COLUMN "users"."is_online" IS '是否在线, F: 离线, T: 在线';

COMMENT ON COLUMN "users"."password_changed_at" IS '上次密码更新时间';

COMMENT ON COLUMN "users"."created_at" IS '创建时间';

COMMENT ON COLUMN "users"."updated_at" IS '更新时间';