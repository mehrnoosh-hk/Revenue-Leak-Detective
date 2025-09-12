-- Drop role if exists
DROP ROLE IF EXISTS service_account;

-- Drop all policies
DROP POLICY IF EXISTS tenant_isolation_users ON users;
DROP POLICY IF EXISTS tenant_isolation_customers ON customers;
DROP POLICY IF EXISTS tenant_isolation_payments ON payments;
DROP POLICY IF EXISTS tenant_isolation_leaks ON leaks;
DROP POLICY IF EXISTS tenant_isolation_events ON events;
DROP POLICY IF EXISTS tenant_isolation_integrations ON integrations;
DROP POLICY IF EXISTS tenant_isolation_actions ON actions;

-- Drop RLS on all tenant-scoped tables
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
ALTER TABLE customers DISABLE ROW LEVEL SECURITY;
ALTER TABLE payments DISABLE ROW LEVEL SECURITY;
ALTER TABLE leaks DISABLE ROW LEVEL SECURITY;
ALTER TABLE events DISABLE ROW LEVEL SECURITY;
ALTER TABLE integrations DISABLE ROW LEVEL SECURITY;
ALTER TABLE actions DISABLE ROW LEVEL SECURITY;

-- Drop function if exists
DROP FUNCTION IF EXISTS current_tenant_id();
DROP FUNCTION IF EXISTS is_service_account();