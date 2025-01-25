-- 好友关系表
CREATE TABLE IF NOT EXISTS "friendships" (
    "user_id" INT NOT NULL,
    "friend_id" INT NOT NULL,
    "created_at" TIMESTAMP NULL,
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
    "sender_id" INT NOT NULL,
    "receiver_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "requested_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "changed_at" TIMESTAMP NULL,
    CONSTRAINT "friend_requests_sender_id_fk" FOREIGN KEY (sender_id) REFERENCES users (id),
    CONSTRAINT "friend_requests_receiver_id_fk" FOREIGN KEY (receiver_id) REFERENCES users (id)
);

COMMENT ON COLUMN "friend_requests"."id" IS '请求ID';

COMMENT ON COLUMN "friend_requests"."sender_id" IS '请求者ID';

COMMENT ON COLUMN "friend_requests"."receiver_id" IS '接收者ID';

COMMENT ON COLUMN "friend_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "friend_requests"."status" IS '状态, 1: 待处理, 2: 已添加, 3: 已过期';

COMMENT ON COLUMN "friend_requests"."requested_at" IS '请求时间';

COMMENT ON COLUMN "friend_requests"."changed_at" IS '变更时间';
