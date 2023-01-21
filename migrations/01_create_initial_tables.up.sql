DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;

CREATE TABLE users
(
    id      INT NOT NULL AUTO_INCREMENT,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(250),
    password VARCHAR(250),  
    email VARCHAR(250),
    PRIMARY KEY (id)
);

CREATE TABLE accounts
(
    id      INT NOT NULL AUTO_INCREMENT,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(250),  
    name VARCHAR(250),  
    balance FLOAT,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (id)
);


CREATE TABLE transactions
(
    id      INT NOT NULL AUTO_INCREMENT,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    amount FLOAT,
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    PRIMARY KEY (id)
);