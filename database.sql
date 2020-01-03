DROP TABLE IF EXISTS account;
create table account (
    id BIGSERIAL PRIMARY KEY,
    email varchar (255) not null,
    hash_password varchar (64) not null,
    username varchar (50) not null,
    created_at timestamp not null,
    last_login timestamp null,
    status smallint not null
);
ALTER TABLE account ADD CONSTRAINT account_username UNIQUE (username);
ALTER TABLE account ADD CONSTRAINT account_email UNIQUE (email);
CREATE INDEX account_username_idx ON account (username,hash_password);
CREATE INDEX account_email_idx ON account (email,hash_password);