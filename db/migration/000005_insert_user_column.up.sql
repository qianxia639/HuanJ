ALTER TABLE "users"
ADD COLUMN "phone" VARCHAR(30) UNIQUE,
ADD COLUMN "birthday" DATE,
ADD COLUMN "bio" TEXT;

COMMENT ON COLUMN "users"."phone" IS '用户手机号';

COMMENT ON COLUMN "users"."birthday" IS '用户生日';

COMMENT ON COLUMN "users"."bio" IS '用户简介或个性签名';