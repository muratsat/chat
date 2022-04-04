CREATE TABLE user (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL UNIQUE ,
    password_hash VARCHAR(100) NOT NULL ,
    PRIMARY KEY (id)
);

CREATE TABLE room (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL ,
    PRIMARY KEY (id)
);

CREATE TABLE participants (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    room_id INT,
    PRIMARY KEY (id)
);

CREATE TABLE message (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    room_id INT,
    text VARCHAR(1000),
    date DATETIME,
    PRIMARY KEY (id)
);

CREATE TABLE friend (
    user_id INT,
    friend_id INT
);

CREATE TABLE authentication_tokens (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    auth_token VARCHAR(255),
    generated_at DATETIME,
    expires_at DATETIME,
    PRIMARY KEY (id)
);

CREATE TABLE friend_requests (
    user_id int,
    friend_id int
);
