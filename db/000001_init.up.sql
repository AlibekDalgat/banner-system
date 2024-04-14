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
    title TEXT,
    text TEXT,
    url TEXT,
    tag_ids INT[] NOT NULL,
    feature_id INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE history_banners (
    history_id SERIAL PRIMARY KEY,
    id INT NOT NULL,
    title TEXT,
    text TEXT,
    url TEXT,
    tag_ids INT[] NOT NULL,
    feature_id INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
