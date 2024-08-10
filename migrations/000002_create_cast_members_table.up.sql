CREATE TABLE cast_members (
  id bigserial PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(32) NOT NULL,
  created_at timestamp(0) with time zone NOT NULL,
  updated_at timestamp(0) with time zone NOT NULL
);