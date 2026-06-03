ALTER TABLE products ADD COLUMN IF NOT EXISTS cost_price NUMERIC(12,2) DEFAULT 0;

CREATE TABLE IF NOT EXISTS stock_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id),
    type        VARCHAR(20) NOT NULL CHECK (type IN ('purchase', 'sale', 'adjustment')),
    quantity    INT NOT NULL,
    cost_price  NUMERIC(12,2) DEFAULT 0,
    total_cost  NUMERIC(12,2) DEFAULT 0,
    note        VARCHAR(500),
    created_at  TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_stock_logs_product ON stock_logs(product_id);
CREATE INDEX idx_stock_logs_type ON stock_logs(type);
CREATE INDEX idx_stock_logs_created ON stock_logs(created_at);
