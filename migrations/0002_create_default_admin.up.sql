DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        RAISE NOTICE 'Inserting admin user';
        INSERT INTO users (username, password, role, created_at)
        VALUES (
            'admin',
            crypt(current_setting('ADMIN_PASSWORD'), gen_salt('bf')),
            'admin',
            NOW()
        );
    ELSE
        RAISE NOTICE 'Admin user already exists';
    END IF;
END $$;