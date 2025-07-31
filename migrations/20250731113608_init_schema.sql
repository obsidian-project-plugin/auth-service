-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE SCHEMA obsidian;


CREATE TABLE obsidian.users (
                                id UUID PRIMARY KEY,
                                username VARCHAR(255) NOT NULL UNIQUE,
                                email VARCHAR(255) UNIQUE,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX  idx_users_email ON obsidian.users (email);


CREATE TABLE obsidian.clients (
                                  id UUID PRIMARY KEY,
                                  name VARCHAR(255) NOT NULL,
                                  client_identifier VARCHAR(255) NOT NULL UNIQUE,
                                  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX  idx_clients_client_identifier ON obsidian.clients (client_identifier);


CREATE TABLE obsidian.user_devices (
                                       id UUID PRIMARY KEY,
                                       user_id UUID REFERENCES obsidian.users(id) ON DELETE CASCADE,
                                       client_id UUID REFERENCES obsidian.clients(id) ON DELETE CASCADE,
                                       device_name VARCHAR(255),
                                       client_user_id VARCHAR(255),
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       last_used_at TIMESTAMP WITH TIME ZONE
);


CREATE OR REPLACE FUNCTION obsidian.update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_trigger
    BEFORE UPDATE ON obsidian.users
    FOR EACH ROW
    EXECUTE PROCEDURE obsidian.update_updated_at();

CREATE TRIGGER clients_update_trigger
    BEFORE UPDATE ON obsidian.clients
    FOR EACH ROW
    EXECUTE PROCEDURE obsidian.update_updated_at();
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
