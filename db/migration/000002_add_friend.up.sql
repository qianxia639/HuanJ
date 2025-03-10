-- 好友关系表
CREATE TABLE IF NOT EXISTS "friendships" (
    "user_id" INT NOT NULL,
    "friend_id" INT NOT NULL,
    "note" VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT "friendships_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "friendships_friend_id_fk" FOREIGN KEY (friend_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "friendships"."user_id" IS '用户ID';

COMMENT ON COLUMN "friendships"."friend_id" IS '好友ID';

COMMENT ON COLUMN "friendships"."note" IS '好友备注';

COMMENT ON COLUMN "friendships"."created_at" IS '创建时间';

-- 好友请求表
CREATE TABLE IF NOT EXISTS "friend_requests" (
    "id" SERIAL PRIMARY KEY,
    "sender_id" INT NOT NULL,
    "receiver_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL DEFAULT '',
    "status" SMALLINT NOT NULL DEFAULT 1,
    "requested_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "expired_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '7 days',
    CONSTRAINT "friend_requests_sender_id_fk" FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "friend_requests_receiver_id_fk" FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "friend_requests"."id" IS '请求ID';

COMMENT ON COLUMN "friend_requests"."sender_id" IS '请求者ID';

COMMENT ON COLUMN "friend_requests"."receiver_id" IS '接收者ID';

COMMENT ON COLUMN "friend_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "friend_requests"."status" IS '请求状态, 1: 待处理, 2: 已同意, 3: 已拒绝, 4: 已过期';

COMMENT ON COLUMN "friend_requests"."requested_at" IS '请求时间';

COMMENT ON COLUMN "friend_requests"."updated_at" IS '变更时间';

COMMENT ON COLUMN "friend_requests"."expired_at" IS '申请过期时间';
