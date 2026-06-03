ALTER TABLE orders DROP COLUMN IF EXISTS paid_at;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_status;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_url;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_transaction_id;
