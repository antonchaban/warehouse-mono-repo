-- liquibase formatted sql

-- changeset anton:3
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       email VARCHAR(100) UNIQUE,
                       first_name VARCHAR(50),
                       last_name VARCHAR(50),
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
                       id BIGSERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE user_roles (
                            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                            PRIMARY KEY (user_id, role_id)
);

-- changeset anton:4
INSERT INTO roles (name) VALUES ('ROLE_ADMIN');
INSERT INTO roles (name) VALUES ('ROLE_STOREKEEPER');
INSERT INTO roles (name) VALUES ('ROLE_LOGISTICIAN');

-- Пароль: admin123
INSERT INTO users (username, password_hash, email, first_name, last_name)
VALUES ('admin', '$2a$10$wW.w4/1.1234567890123456789012345678901234567890', 'admin@company.com', 'Super', 'Admin');

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r WHERE u.username = 'admin' AND r.name = 'ROLE_ADMIN';