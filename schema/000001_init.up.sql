CREATE TABLE IF NOT EXISTS rassilka
(
    id              SERIAL PRIMARY KEY,
    start_time      TIMESTAMPTZ NOT NULL,
    message         TEXT NOT NULL,
    filter          TEXT[],
    end_time        TIMESTAMPTZ NOT NULL,
    supplemented    BOOL NOT NULL,
    recreated       BOOL NOT NULL
    CONSTRAINT valid_times CHECK (start_time < end_time)
);


CREATE TABLE IF NOT EXISTS clients
(
    id              SERIAL PRIMARY KEY,
    phone           VARCHAR(7) NOT NULL CHECK (phone ~ '^[7]\d{6}$'),
    code            VARCHAR(3) NOT NULL CHECK (code ~ '^\d{3}$'),
    tag             VARCHAR(50),
    timezone        VARCHAR(50) NOT NULL CHECK (timezone ~ '^[A-Za-z/0-9_]+$'),
    CONSTRAINT unique_code_phone_combo UNIQUE (code, phone)
);

CREATE SEQUENCE messages_id_seq;
