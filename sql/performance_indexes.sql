-- High-QPS read path indexes. Run these once after schema creation.
-- Run once on a fresh schema. If an index already exists, skip that statement.

CREATE INDEX idx_order_user_status_created
    ON `order` (user_id, status, created_at DESC);

CREATE INDEX idx_order_user_created
    ON `order` (user_id, created_at DESC);

CREATE INDEX idx_order_available_created
    ON `order` (status, rider_id, can_reassign, created_at DESC);

CREATE INDEX idx_order_available_reward
    ON `order` (status, rider_id, can_reassign, reward_amount);

CREATE INDEX idx_address_user_type_default
    ON address (user_id, type, is_default);
