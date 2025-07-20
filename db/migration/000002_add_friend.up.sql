-- 好友关系表
CREATE TABLE IF NOT EXISTS "friendships" (
    "user_id" INT NOT NULL,
    "friend_id" INT NOT NULL,
    "remark" VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT "friendships_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "friendships_friend_id_fk" FOREIGN KEY (friend_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON TABLE "friendships" IS '好友关系表';

COMMENT ON COLUMN "friendships"."user_id" IS '用户ID';

COMMENT ON COLUMN "friendships"."friend_id" IS '好友ID';

COMMENT ON COLUMN "friendships"."remark" IS '好友备注';

COMMENT ON COLUMN "friendships"."created_at" IS '创建时间';

COMMENT ON COLUMN "friendships"."updated_at" IS '更新时间';

-- 好友申请表
CREATE TABLE IF NOT EXISTS "friend_requests" (
    "id" SERIAL PRIMARY KEY,
    "from_user_id" INT NOT NULL,
    "to_user_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL DEFAULT '',
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "friend_requests_from_user_id_fk" FOREIGN KEY (from_user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "friend_requests_to_user_id_fk" FOREIGN KEY (to_user_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON TABLE "friend_requests" IS '好友申请表';

COMMENT ON COLUMN "friend_requests"."id" IS '请求ID';

COMMENT ON COLUMN "friend_requests"."from_user_id" IS '申请者ID';

COMMENT ON COLUMN "friend_requests"."to_user_id" IS '接收者ID';

COMMENT ON COLUMN "friend_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "friend_requests"."status" IS '请求状态, 1: 待处理, 2: 已同意, 3: 已拒绝';

COMMENT ON COLUMN "friend_requests"."created_at" IS '请求时间';

COMMENT ON COLUMN "friend_requests"."updated_at" IS '变更时间';