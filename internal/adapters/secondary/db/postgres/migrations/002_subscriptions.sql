-- +goose Up
CREATE TABLE subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    
    -- Subscription details
    plan VARCHAR(50) NOT NULL DEFAULT 'basic',        -- e.g., basic, pro, premium
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    auto_renew BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Payment/trial info
    last_payment_at TIMESTAMP,
    last_payment_amount NUMERIC(10,2),
    trial_ends_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_expires_after_start CHECK (expires_at IS NULL OR expires_at >= started_at)
);

-- Trigger to auto-update timestamps
CREATE TRIGGER update_subscription_timestamp
    BEFORE UPDATE ON subscription
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE subscription;
