CREATE TABLE users (
    id bigserial not null primary key,
    name varchar unique not null,
    age integer not null,
    friends bigint[] not null
);
