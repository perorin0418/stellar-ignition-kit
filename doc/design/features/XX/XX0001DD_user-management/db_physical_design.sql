-- db_physical_design.sql
-- 物理設計テンプレート

CREATE TABLE users (
    user_id      VARCHAR(36) PRIMARY KEY,
    user_name    VARCHAR(100) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL
);
