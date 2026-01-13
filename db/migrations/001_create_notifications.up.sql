-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    
    -- Content
    type VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Delivery channels
    channels TEXT[] NOT NULL DEFAULT '{}',
    
    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    
    -- Read status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    
    -- Delivery tracking
    email_sent_at TIMESTAMP,
    push_sent_at TIMESTAMP,
    
    -- Retry mechanism
    retry_count INT DEFAULT 0,
    last_error TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '30 days')
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_user_notifications ON notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_unread_notifications ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_pending_notifications ON notifications(status) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_expires_at ON notifications(expires_at);

-- Create notification_preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    
    -- Global toggles
    email_enabled BOOLEAN DEFAULT TRUE,
    push_enabled BOOLEAN DEFAULT TRUE,
    in_app_enabled BOOLEAN DEFAULT TRUE,
    
    -- Category preferences
    order_notifications JSONB DEFAULT '{"email": true, "push": true, "in_app": true}'::jsonb,
    account_notifications JSONB DEFAULT '{"email": true, "push": false, "in_app": true}'::jsonb,
    product_notifications JSONB DEFAULT '{"email": false, "push": true, "in_app": true}'::jsonb,
    
    -- Quiet hours
    quiet_hours_enabled BOOLEAN DEFAULT FALSE,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);