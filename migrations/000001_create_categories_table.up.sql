CREATE TABLE IF NOT EXISTS categories(
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(4000),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at timestamp(0) with time zone NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL,
    deleted_at timestamp(0) with time zone NULL
);