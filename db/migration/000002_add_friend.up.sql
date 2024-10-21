-- 好友表
CREATE TABLE IF NOT EXISTS "friends" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL,
    "friend_id" INT NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, friend_id),  -- 确保一对好友关系唯一
    CONSTRAINT "friends_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT "friends_friend_id_fk" FOREIGN KEY (friend_id) REFERENCES users (id)
);

COMMENT ON COLUMN "friends"."id" IS '好友关系标识';

COMMENT ON COLUMN "friends"."user_id" IS '用户ID';

COMMENT ON COLUMN "friends"."friend_id" IS '好友的用户ID';

COMMENT ON COLUMN "friends"."status" IS '关系状态, 1: 待确认, 2: 已确认, 3: 已拒绝';

COMMENT ON COLUMN "friends"."created_at" IS '关系创建时间';

COMMENT ON COLUMN "friends"."created_at" IS '更新时间';