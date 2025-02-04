CREATE TABLE IF NOT EXISTS "invitation_codes" (
    "id" SERIAL PRIMARY KEY,
    "code" VARCHAR(64) NOT NULL UNIQUE,
    "user_id" INT,
    "used_at" TIMESTAMPTZ,
    "status" INT NOT NULL DEFAULT -1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expired_at" TIMESTAMPTZ
);

COMMENT ON COLUMN "invitation_codes"."id" IS '编号ID';

COMMENT ON COLUMN "invitation_codes"."code" IS '邀请码';

COMMENT ON COLUMN "invitation_codes"."user_id" IS '使用者Id';

COMMENT ON COLUMN "invitation_codes"."used_at" IS '使用时间';

COMMENT ON COLUMN "invitation_codes"."status" IS '状态, -1: 未使用, 1: 已使用, -2: 已过期';

COMMENT ON COLUMN "invitation_codes"."created_at" IS '创建时间';

COMMENT ON COLUMN "invitation_codes"."expired_at" IS '过期时间';
