-- 好友关系表
CREATE TABLE IF NOT EXISTS "friendships" (
    "user_id" INT NOT NULL,
    "friend_id" INT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT "friendships_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT "friendships_friend_id_fk" FOREIGN KEY (friend_id) REFERENCES users (id)
);

COMMENT ON COLUMN "friendships"."user_id" IS '用户ID';

COMMENT ON COLUMN "friendships"."friend_id" IS '好友的用户ID';

COMMENT ON COLUMN "friendships"."created_at" IS '创建时间';

-- 好友请求表
CREATE TABLE IF NOT EXISTS "friend_requests" (
    "id" SERIAL PRIMARY KEY,
    "from_user_id" INT NOT NULL,
    "to_user_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "requested_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "changed_at" TIMESTAMP NOT NULL DEFAULT '0001-01-01 00:00:00',
    CONSTRAINT "friend_requests_from_user_id_fk" FOREIGN KEY (from_user_id) REFERENCES users (id),
    CONSTRAINT "friend_requests_to_user_id_fk" FOREIGN KEY (to_user_id) REFERENCES users (id)
);

COMMENT ON COLUMN "friend_requests"."id" IS '请求ID';

COMMENT ON COLUMN "friend_requests"."from_user_id" IS '请求者ID';

COMMENT ON COLUMN "friend_requests"."to_user_id" IS '接收者ID';

COMMENT ON COLUMN "friend_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "friend_requests"."status" IS '请求状态, 1: 待处理, 2: 已同意, 3: 已拒绝, 4: 已忽略';

COMMENT ON COLUMN "friend_requests"."requested_at" IS '请求时间';

COMMENT ON COLUMN "friend_requests"."changed_at" IS '变更时间';
