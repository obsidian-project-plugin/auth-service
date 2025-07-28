-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE SCHEMA IF NOT EXISTS dco;


CREATE TABLE dco.users (
                           id UUID PRIMARY KEY,
                           username VARCHAR(255) NOT NULL UNIQUE,
                           email VARCHAR(255) UNIQUE,
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX idx_users_email ON dco.users (email);


CREATE TABLE dco.clients (
                             id UUID PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             client_identifier VARCHAR(255) NOT NULL UNIQUE,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                             updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX idx_clients_client_identifier ON dco.clients (client_identifier);


CREATE TABLE dco.user_devices (
                                  id UUID PRIMARY KEY,
                                  user_id UUID REFERENCES dco.users(id) ON DELETE CASCADE,
                                  client_id UUID REFERENCES dco.clients(id) ON DELETE CASCADE,
                                  device_name VARCHAR(255),
                                  client_user_id VARCHAR(255),
                                  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                  last_used_at TIMESTAMP WITH TIME ZONE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
