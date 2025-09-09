-- Create action_type ENUM
CREATE TYPE action_type_enum AS ENUM (
    'retry_payment',
    'outreach',
    'linear_task',
    'email',
    'other'
);

-- Create action_status ENUM
CREATE TYPE action_status_enum AS ENUM (
    'pending',
    'approved',
    'modified',
    'denied'
);

-- Create action_result ENUM
CREATE TYPE action_result_enum AS ENUM (
    'success',
    'failure',
    'pending',
    'other'
);

-- Create actions table
CREATE TABLE actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    leak_id UUID NOT NULL,
    action_type action_type_enum NOT NULL,
    status action_status_enum NOT NULL,
    result action_result_enum NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (leak_id) REFERENCES leaks(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_actions_leak_id ON actions(leak_id);
CREATE INDEX idx_actions_action_type ON actions(action_type);
CREATE INDEX idx_actions_status ON actions(status);
CREATE INDEX idx_actions_result ON actions(result);

-- Create updated_at trigger for actions
CREATE TRIGGER update_actions_updated_at BEFORE UPDATE ON actions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();