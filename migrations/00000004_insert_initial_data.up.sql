INSERT INTO users (login, password_hash)
VALUES ('alice', encode(digest('secret', 'md5'), 'hex'))
    ON CONFLICT DO NOTHING;