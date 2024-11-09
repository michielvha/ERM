DO $$
DECLARE
    admin_password TEXT := 'PLACEHOLDER_PASSWORD';
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        RAISE NOTICE 'Inserting admin user';
        INSERT INTO users (username, password, role, created_at)
        VALUES (
            'admin',
            crypt(admin_password, gen_salt('bf')),
            'admin',
            NOW()
        );
    ELSE
        RAISE NOTICE 'Admin user already exists';
    END IF;
END $$;