-- 群组表
CREATE TABLE IF NOT EXISTS "groups" (
    "id" SERIAL PRIMARY KEY,
    "group_name" VARCHAR(64) NOT NULL UNIQUE,
    "creator_id" INT NOT NULL,
    "group_avatar_url" VARCHAR(512),
    "description" VARCHAR(255),
    "max_member" INT DEFAULT 500,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00',
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
    "waiting" BOOLEAN NOT NULL DEFAULT true,
    "joined_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    CONSTRAINT "group_members_group_id_fk" FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT "group_members_user_id_fk" FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

COMMENT ON COLUMN "group_members"."group_id" IS '群组ID';

COMMENT ON COLUMN "group_members"."user_id" IS '用户ID';

COMMENT ON COLUMN "group_members"."role" IS '成员角色, 1: 群主, 2: 管理员, 3: 普通成员';

COMMENT ON COLUMN "group_members"."waiting" IS '等待同意, f: 已同意, t: 未同意';

COMMENT ON COLUMN "group_members"."joined_at" IS '加入时间';