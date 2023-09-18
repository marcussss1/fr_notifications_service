CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS client (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    phone_number VARCHAR(11),
    mobile_operator_code INTEGER,
    tag VARCHAR(100),
    time_zone TIMESTAMP WITH TIME ZONE DEFAULT (now() at time zone 'UTC')
);

CREATE TABLE IF NOT EXISTS mailing (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at time zone 'UTC'),
    text VARCHAR(1000),
    filter VARCHAR(100),
    finished_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at time zone 'UTC')
);

CREATE TABLE IF NOT EXISTS pending_mailing (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    mailing_id UUID,
    status INTEGER
);

CREATE TABLE IF NOT EXISTS message (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at time zone 'UTC'),
    status INTEGER,
    mailing_id UUID,
    receiver_client_id UUID
);
