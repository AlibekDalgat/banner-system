CREATE TABLE admins (
    login VARCHAR(255),
    password VARCHAR(255)
);
INSERT INTO admins VALUES ('admin', 'admin');

CREATE TABLE users (
   login VARCHAR(255),
   password VARCHAR(255)
);
INSERT INTO users VALUES ('user', 'user'), ('alibek', 'alibek'), ('dalghat', 'dalghat');

CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    url TEXT NOT NULL,
    tag_ids INT[] NOT NULL,
    feature_id INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);