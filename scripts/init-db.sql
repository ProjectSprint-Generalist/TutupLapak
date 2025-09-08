-- Grant schema privileges to tutuplapak user
GRANT ALL ON SCHEMA public TO tutuplapak;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tutuplapak;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO tutuplapak;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tutuplapak;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tutuplapak;
