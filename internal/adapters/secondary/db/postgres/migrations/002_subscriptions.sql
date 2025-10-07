-- +goose Up
CREATE TABLE subscription (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Subscription details
    plan VARCHAR(50) NOT NULL,        -- e.g., basic, premium
    billing_period VARCHAR(20) NOT NULL, -- e.g, monthly, yearly
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    auto_renew BOOLEAN NOT NULL,
    cancelled_at TIMESTAMP,
    
    -- Payment info
    last_payment_at TIMESTAMP,
    last_payment_amount NUMERIC(10,2),
    last_payment_currency VARCHAR(3), -- USD, EUR, GBP, etc.

    -- Trial 
    trial_ends_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    -- Constraints
    CONSTRAINT chk_expires_after_start CHECK (expires_at IS NULL OR expires_at >= started_at)
);

-- Trigger to auto-update timestamps
CREATE TRIGGER update_subscription_timestamp
    BEFORE UPDATE ON subscription
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE subscription;
