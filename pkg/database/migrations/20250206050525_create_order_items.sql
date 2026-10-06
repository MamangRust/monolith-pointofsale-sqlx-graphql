-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order_items" (
    "order_item_id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL REFERENCES "orders" ("order_id") ON DELETE CASCADE,
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "quantity" INT NOT NULL,
    "price" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);

CREATE INDEX idx_order_items_product_id ON order_items (product_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_items_order_id ON order_items (order_id);

DROP INDEX IF EXISTS idx_order_items_product_id ON order_items (product_id);

DROP TABLE IF EXISTS "order_items";

-- +goose StatementEnd
