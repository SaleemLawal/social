CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INTEGER NOT NULL DEFAULT 1,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES ('User', 1, 'User role');
INSERT INTO roles (name, level, description) VALUES ('Moderator', 2, 'Moderator role');
INSERT INTO roles (name, level, description) VALUES ('Admin', 3, 'Admin role');