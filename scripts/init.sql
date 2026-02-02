-- Flash Sale System Database Initialization
-- PostgreSQL 17.2

-- ============================================\dt
-- Product Domain Tables
-- ============================================

-- Products table (Aggregate Root)
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status SMALLINT NOT NULL DEFAULT 9 CHECK (status IN (1, 9)),
    available_stock INT NOT NULL DEFAULT 0 CHECK (available_stock >= 0),
    reserved_stock INT NOT NULL DEFAULT 0 CHECK (reserved_stock >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE products IS 'Product aggregate root';
COMMENT ON COLUMN products.status IS '1=active, 9=inactive';
COMMENT ON COLUMN products.available_stock IS 'Available stock for purchase';
COMMENT ON COLUMN products.reserved_stock IS 'Reserved stock for pending orders';

-- Product pricing table (Aggregate Root)
CREATE TABLE IF NOT EXISTS product_pricing (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('USD', 'TWD', 'JPY')),
    amount DECIMAL(19, 4) NOT NULL CHECK (amount >= 0),
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT uk_product_currency_period UNIQUE (product_id, currency, valid_from),
    CONSTRAINT valid_period CHECK (valid_until IS NULL OR valid_until > valid_from)
);

COMMENT ON TABLE product_pricing IS 'Product pricing with multi-currency and time-based periods';
COMMENT ON COLUMN product_pricing.valid_until IS 'NULL means valid indefinitely';

-- ============================================
-- Order Domain Tables
-- ============================================

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    total_price DECIMAL(19, 4) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- ============================================
-- Payment Domain Tables
-- ============================================

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    payment_method VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- ============================================
-- Indexes for Performance
-- ============================================

-- Product indexes
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_available_stock ON products(available_stock) WHERE available_stock > 0;

-- Product pricing indexes
CREATE INDEX idx_product_pricing_product_valid ON product_pricing(product_id, valid_from, valid_until);
-- Note: Cannot use CURRENT_TIMESTAMP in partial index (not IMMUTABLE)
-- Query will filter by time at runtime instead
CREATE INDEX idx_product_pricing_currency ON product_pricing(product_id, currency);

-- Order indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_product_id ON orders(product_id);
CREATE INDEX idx_orders_status ON orders(status);

-- Payment indexes
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);

-- ============================================
-- Sample Data
-- ============================================

-- Insert sample products
INSERT INTO products (sku, name, description, status, available_stock, reserved_stock) VALUES
    ('IPHONE-15-PRO', 'Flash Sale iPhone 15 Pro', 'Latest iPhone with A17 Pro chip', 1, 100, 0),
    ('MACBOOK-PRO-16', 'Flash Sale MacBook Pro 16"', 'M3 Max MacBook Pro', 1, 50, 0),
    ('AIRPODS-PRO-2', 'Flash Sale AirPods Pro 2', 'Active Noise Cancellation', 1, 200, 0);

-- Insert sample pricing (multi-currency)
INSERT INTO product_pricing (product_id, currency, amount, valid_from, valid_until) VALUES
    -- iPhone 15 Pro pricing
    (1, 'USD', 999.99, '2026-01-01 00:00:00', NULL),
    (1, 'TWD', 31000, '2026-01-01 00:00:00', NULL),
    (1, 'JPY', 145000, '2026-01-01 00:00:00', NULL),
    
    -- MacBook Pro pricing
    (2, 'USD', 2499.99, '2026-01-01 00:00:00', NULL),
    (2, 'TWD', 77000, '2026-01-01 00:00:00', NULL),
    (2, 'JPY', 363000, '2026-01-01 00:00:00', NULL),
    
    -- AirPods Pro pricing
    (3, 'USD', 249.99, '2026-01-01 00:00:00', NULL),
    (3, 'TWD', 7700, '2026-01-01 00:00:00', NULL),
    (3, 'JPY', 36300, '2026-01-01 00:00:00', NULL);

-- Create function to update timestamp automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for auto-updating timestamps
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
