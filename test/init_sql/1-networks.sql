create extension if not EXISTS "uuid-ossp";

create table if not exists networks (
  id           uuid         PRIMARY KEY DEFAULT uuid_generate_v4(),
  name         VARCHAR(255) NOT NULL UNIQUE,
  private_key  VARCHAR(45)  NOT NULL UNIQUE,
  port         INTEGER      NOT NULL,
  ipv4         VARCHAR(20)  NOT NULL
)

