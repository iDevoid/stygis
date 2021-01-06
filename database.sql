DROP TABLE IF exists users;
create table users (
    id BIGSERIAL PRIMARY KEY,
    username varchar(50) not null,
    email varchar(255) not null,
    hashed_email VARCHAR(128) not NULL,
    "password" varchar(128) not null,
    create_time timestamp not null DEFAULT CURRENT_TIMESTAMP,
    status smallint not null DEFAULT 0,
    update_time TIMESTAMP
);
ALTER TABLE users ADD CONSTRAINT user_username UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT user_email UNIQUE (hashed_email);
CREATE INDEX user_username_idx ON users (username, hashed_email,"password");

DROP TABLE IF exists "profiles";
create table "profiles" (
    id BIGINT PRIMARY KEY,
    username varchar(50) not null,
    full_name VARCHAR(255) not null DEFAULT '',
    profile_picture text not null DEFAULT '',
    cover_picture text not null DEFAULT '',
    bio text not null DEFAULT '',
    "card" text not null DEFAULT '',
    followers INTEGER not null DEFAULT 0,
    following INTEGER not null DEFAULT 0,
    register_time TIMESTAMP NULL,
    status smallint not null DEFAULT 0,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NULl
);
ALTER TABLE "profiles" ADD CONSTRAINT profile_id UNIQUE (id);
ALTER TABLE "profiles" ADD CONSTRAINT profile_username UNIQUE (username);
CREATE INDEX profile_username_idx ON "profiles" (username);
