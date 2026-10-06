-- +goose Up
-- +goose StatementBegin
CREATE TABLE "cashiers" (
    "cashier_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id"),
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "name" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_cashiers_merchant_id ON cashiers (merchant_id);

CREATE INDEX idx_cashiers_user_id ON cashiers (user_id);

CREATE INDEX idx_cashiers_created_at ON cashiers (created_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cashiers_merchant_id ON cashiers (merchant_id);

DROP INDEX IF EXISTS idx_cashiers_user_id ON cashiers (user_id);

DROP INDEX IF EXISTS idx_cashiers_created_at ON cashiers (created_at);

DROP TABLE IF EXISTS "cashiers";

-- +goose StatementEnd
