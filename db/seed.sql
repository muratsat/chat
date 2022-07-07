CREATE DATABASE chat;

CREATE TABLE users (
    id SERIAL,
    username varchar,
    password_hash varchar,
    PRIMARY KEY(id)
);

CREATE TABLE rooms (
    id SERIAL,
    name varchar,
    PRIMARY KEY (id)
);

CREATE TABLE participants (
    user_id integer,
    room_id integer,
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (room_id) REFERENCES rooms
);

CREATE TABLE messages (
    id SERIAL,
    user_id integer,
    room_id integer,
    text text,
    date date,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (room_id) REFERENCES rooms
);

CREATE TABLE friends (
    user_id integer,
    friend_id integer,
    ir_request boolean,
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (friend_id) REFERENCES users
)
