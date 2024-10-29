-- migrations/init.sql

-- Создание роли (пользователя)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'your_username') THEN
        CREATE ROLE your_username WITH LOGIN PASSWORD 'your_password';
    END IF;
END $$;

-- Создание базы данных
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'your_database') THEN
        CREATE DATABASE your_database WITH OWNER your_username;
    END IF;
END $$;

-- Предоставление прав пользователю
GRANT ALL PRIVILEGES ON DATABASE your_database TO your_username;
