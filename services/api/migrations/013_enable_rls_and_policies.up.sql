-- Enable RLS on all tenant-scoped tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE leaks ENABLE ROW LEVEL SECURITY;
ALTER TABLE events ENABLE ROW LEVEL SECURITY;
ALTER TABLE integrations ENABLE ROW LEVEL SECURITY;
ALTER TABLE actions ENABLE ROW LEVEL SECURITY;

-- Create function to get current tenant ID from session
CREATE OR REPLACE FUNCTION current_tenant_id() RETURNS UUID AS $$
BEGIN
    RETURN COALESCE(
        current_setting('app.current_tenant_id', true)::UUID,
        '00000000-0000-0000-0000-000000000000'::UUID
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create function to check if current user is service account
CREATE OR REPLACE FUNCTION is_service_account() RETURNS BOOLEAN AS $$
BEGIN
    RETURN COALESCE(
        current_setting('app.is_service_account', true)::BOOLEAN,
        false
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- users table policies
CREATE POLICY tenant_isolation_users ON users
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Customers table policies  
CREATE POLICY tenant_isolation_customers ON customers
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Payments table policies
CREATE POLICY tenant_isolation_payments ON payments
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Leaks table policies
CREATE POLICY tenant_isolation_leaks ON leaks
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Events table policies
CREATE POLICY tenant_isolation_events ON events
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Integrations table policies
CREATE POLICY tenant_isolation_integrations ON integrations
    FOR ALL
    TO PUBLIC
    USING (tenant_id = current_tenant_id() OR is_service_account())
    WITH CHECK (tenant_id = current_tenant_id() OR is_service_account());

-- Actions table policies (via leak relationship)
CREATE POLICY tenant_isolation_actions ON actions
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM leaks 
            WHERE leaks.id = actions.leak_id 
            AND leaks.tenant_id = current_tenant_id()
        ) OR is_service_account()
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM leaks 
            WHERE leaks.id = actions.leak_id 
            AND leaks.tenant_id = current_tenant_id()
        ) OR is_service_account()
    );

-- Create service account role
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'service_account') THEN
        CREATE ROLE service_account;
    END IF;
END
$$;
GRANT USAGE ON SCHEMA public TO service_account;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO service_account;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO service_account;