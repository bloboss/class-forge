-- Test database initialization script
-- This script runs when the test database container starts

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC';

-- Create test user with necessary permissions
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'fgc_test') THEN
        CREATE ROLE fgc_test WITH LOGIN PASSWORD 'fgc_test_password';
    END IF;
END
$$;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE forgejo_classroom_test TO fgc_test;
GRANT ALL ON SCHEMA public TO fgc_test;

-- Test data cleanup function (called before each test)
CREATE OR REPLACE FUNCTION reset_test_data() RETURNS void AS $$
BEGIN
    -- Truncate all tables in dependency order
    TRUNCATE TABLE IF EXISTS team_members RESTART IDENTITY CASCADE;
    TRUNCATE TABLE IF EXISTS teams RESTART IDENTITY CASCADE;
    TRUNCATE TABLE IF EXISTS submissions RESTART IDENTITY CASCADE;
    TRUNCATE TABLE IF EXISTS assignments RESTART IDENTITY CASCADE;
    TRUNCATE TABLE IF EXISTS roster_entries RESTART IDENTITY CASCADE;
    TRUNCATE TABLE IF EXISTS classrooms RESTART IDENTITY CASCADE;

    -- Reset sequences
    SELECT setval(pg_get_serial_sequence('classrooms', 'id'), 1, false);
    SELECT setval(pg_get_serial_sequence('roster_entries', 'id'), 1, false);
    SELECT setval(pg_get_serial_sequence('assignments', 'id'), 1, false);
    SELECT setval(pg_get_serial_sequence('submissions', 'id'), 1, false);
    SELECT setval(pg_get_serial_sequence('teams', 'id'), 1, false);
    SELECT setval(pg_get_serial_sequence('team_members', 'id'), 1, false);
END;
$$ LANGUAGE plpgsql;

-- Insert test data function (for integration tests)
CREATE OR REPLACE FUNCTION insert_test_data() RETURNS void AS $$
BEGIN
    -- Insert test classroom
    INSERT INTO classrooms (id, name, slug, description, organization_name, organization_id, instructor_id, instructor_login, public, archived, created_at, updated_at)
    VALUES (1, 'Test Classroom', 'test-classroom', 'A test classroom for integration tests', 'test-org', 1, 1, 'instructor', true, false, NOW(), NOW());

    -- Insert test roster entries
    INSERT INTO roster_entries (id, classroom_id, student_name, student_email, student_id, forgejo_username, forgejo_user_id, role, linked_at, created_at, updated_at)
    VALUES
        (1, 1, 'John Doe', 'john@example.com', 'john123', 'johndoe', 101, 'student', NOW(), NOW(), NOW()),
        (2, 1, 'Jane Smith', 'jane@example.com', 'jane456', 'janesmith', 102, 'student', NOW(), NOW(), NOW()),
        (3, 1, 'Bob Wilson', 'bob@example.com', 'bob789', NULL, NULL, 'student', NULL, NOW(), NOW());

    -- Insert test assignment
    INSERT INTO assignments (id, classroom_id, name, slug, description, template_repository, template_repository_id, deadline, max_team_size, auto_accept, public, created_at, updated_at)
    VALUES (1, 1, 'Test Assignment', 'test-assignment', 'A test assignment', 'https://example.com/template', 201, NOW() + INTERVAL '7 days', 2, false, true, NOW(), NOW());

    -- Reset sequences to continue from inserted data
    SELECT setval(pg_get_serial_sequence('classrooms', 'id'), 1, true);
    SELECT setval(pg_get_serial_sequence('roster_entries', 'id'), 3, true);
    SELECT setval(pg_get_serial_sequence('assignments', 'id'), 1, true);
END;
$$ LANGUAGE plpgsql;