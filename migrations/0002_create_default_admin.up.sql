INSERT INTO users (username, password, role, created_at)
VALUES (
    'admin',
    crypt('secure_admin_password', gen_salt('bf')), -- Replace with a secure hashed password
    'admin',
    NOW()
)
ON CONFLICT (username) DO NOTHING; -- Prevents errors if the admin user already exists