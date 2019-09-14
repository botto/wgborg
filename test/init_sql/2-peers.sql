CREATE extension if NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS peers (
  id          uuid         PRIMARY KEY DEFAULT uuid_generate_v4(),
  public_key  VARCHAR(45)  NOT NULL UNIQUE,
  peer_name   VARCHAR(255) NOT NULL,
  psk         VARCHAR(45)  NOT NULL,
  cidr        VARCHAR(20)  NOT NULL,
  network     uuid         REFERENCES networks(id)
);

