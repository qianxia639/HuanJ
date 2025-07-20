-- 删除索引
DROP INDEX IF EXISTS idx_unique_pending_request;

DROP TABLE IF EXISTS friendships;

DROP TABLE IF EXISTS friend_requests;

-- 删除枚举类型
DROP TYPE IF EXISTS friendship_status;