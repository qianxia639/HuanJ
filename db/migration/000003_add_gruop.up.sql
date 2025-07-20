-- 群组表
CREATE TABLE IF NOT EXISTS "groups" (
    "id" SERIAL PRIMARY KEY,
    "group_name" VARCHAR(64) NOT NULL UNIQUE,
    "owner_id" INT NOT NULL,
    "avatar_url" VARCHAR(512) NOT NULL DEFAULT '',
    "description" VARCHAR(255) NOT NULL DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    CONSTRAINT "groups_owner_id_fk" FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON TABLE "groups" IS '群组表';

COMMENT ON COLUMN "groups"."id" IS '群组ID';

COMMENT ON COLUMN "groups"."group_name" IS '群组名称';

COMMENT ON COLUMN "groups"."owner_id" IS '群主ID';

COMMENT ON COLUMN "groups"."avatar_url" IS '群组头像URL';

COMMENT ON COLUMN "groups"."description" IS '群组描述';

COMMENT ON COLUMN "groups"."created_at" IS '群组创建时间';

COMMENT ON COLUMN "groups"."updated_at" IS '群组信息更新时间';

-- 群组申请表
CREATE TABLE IF NOT EXISTS "group_requests" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL,
    "group_id" INT NOT NULL,
    "request_desc" VARCHAR(100) NOT NULL DEFAULT '',
    "status" SMALLINT NOT NULL DEFAULT 1,
    "processor_id" INT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    CONSTRAINT "group_requests_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT "group_requests_group_id_fk" FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT "group_requests_processor_id_fk" FOREIGN KEY (processor_id) REFERENCES user (id) ON DELETE CASCADE
);

COMMENT ON TABLE "group_requests" IS '群组申请表';

COMMENT ON COLUMN "group_requests"."id" IS '请求ID';

COMMENT ON COLUMN "group_requests"."user_id" IS '申请者ID';

COMMENT ON COLUMN "group_requests"."group_id" IS '群组ID';

COMMENT ON COLUMN "group_requests"."request_desc" IS '请求信息';

COMMENT ON COLUMN "group_requests"."status" IS '请求状态, 1: 待处理, 2: 已同意, 3: 已拒绝';

COMMENT ON COLUMN "group_requests"."processor_id" IS '处理者ID';

COMMENT ON COLUMN "group_requests"."created_at" IS '创建时间';

COMMENT ON COLUMN "group_requests"."updated_at" IS '变更时间';

-- 群组成员表
CREATE TABLE IF NOT EXISTS "group_members" (
    "group_id" INT NOT NULL,
    "user_id" INT NOT NULL,
    "role"  SMALLINT NOT NULL DEFAULT 3,
    "mute_until" TIMESTAMPTZ DEFAULT NULL,
    "nickname" VARCHAR(20) NOT NULL,
    "joined_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    CONSTRAINT "group_members_group_id_fk" FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT "group_members_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON TABLE "group_members" IS '群组成员表';

COMMENT ON COLUMN "group_members"."group_id" IS '群组ID';

COMMENT ON COLUMN "group_members"."user_id" IS '用户ID';

COMMENT ON COLUMN "group_members"."role" IS '成员角色, 1: 群主, 2: 管理员, 3: 普通成员';

COMMENT ON COLUMN "group_members"."mute_until" IS '禁言截止时间(NULL表示未被禁言)';

COMMENT ON COLUMN "group_members"."nickname" IS '群内昵称';

COMMENT ON COLUMN "group_members"."joined_at" IS '加入时间';