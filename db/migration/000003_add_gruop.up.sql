-- 群组请求表
CREATE TABLE IF NOT EXISTS "group_requests" (
    "id" SERIAL PRIMARY KEY,
    "sender_id" INT NOT NULL,
    "receiver_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL DEFAULT '',
    "status" SMALLINT NOT NULL DEFAULT 1,
    "requested_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "expired_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '7 days',
    CONSTRAINT "group_requests_sender_id_fk" FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "group_requests_receiver_id_fk" FOREIGN KEY (receiver_id) REFERENCES groups (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "group_requests"."id" IS '请求ID';

COMMENT ON COLUMN "group_requests"."sender_id" IS '请求者ID';

COMMENT ON COLUMN "group_requests"."receiver_id" IS '接收者ID';

COMMENT ON COLUMN "group_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "group_requests"."status" IS '请求状态, 1: 待处理, 2: 已同意, 3: 已拒绝, 4: 已忽略';

COMMENT ON COLUMN "group_requests"."requested_at" IS '请求时间';

COMMENT ON COLUMN "group_requests"."updated_at" IS '变更时间';

COMMENT ON COLUMN "group_requests"."expired_at" IS '申请过期时间';

-- 群组表
CREATE TABLE IF NOT EXISTS "groups" (
    "id" SERIAL PRIMARY KEY,
    "group_name" VARCHAR(64) NOT NULL UNIQUE,
    "creator_id" INT NOT NULL,
    "group_avatar_url" VARCHAR(512) NOT NULL DEFAULT '',
    "description" VARCHAR(255) NOT NULL DEFAULT '',
    "max_member" INT NOT NULL DEFAULT 500,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    CONSTRAINT "groups_creator_id_fk" FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "groups"."id" IS '群组ID';

COMMENT ON COLUMN "groups"."group_name" IS '群组名称';

COMMENT ON COLUMN "groups"."creator_id" IS '创建者ID';

COMMENT ON COLUMN "groups"."group_avatar_url" IS '群组头像URL';

COMMENT ON COLUMN "groups"."description" IS '群组描述';

COMMENT ON COLUMN "groups"."max_member" IS '群组最大成员数, 默认500';

COMMENT ON COLUMN "groups"."created_at" IS '群组创建时间';

COMMENT ON COLUMN "groups"."updated_at" IS '群组信息更新时间';

-- 群组成员表
CREATE TABLE IF NOT EXISTS "group_members" (
    "group_id" INT NOT NULL,
    "user_id" INT NOT NULL,
    "role"  SMALLINT NOT NULL DEFAULT 3,
    "joined_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    CONSTRAINT "group_members_group_id_fk" FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT "group_members_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "group_members"."group_id" IS '群组ID';

COMMENT ON COLUMN "group_members"."user_id" IS '用户ID';

COMMENT ON COLUMN "group_members"."role" IS '成员角色, 1: 群主, 2: 管理员, 3: 普通成员';

COMMENT ON COLUMN "group_members"."joined_at" IS '加入时间';