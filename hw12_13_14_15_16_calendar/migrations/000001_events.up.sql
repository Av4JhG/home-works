create table if not exists events(
    id bigserial not null primary key,
    user_id bigint not null,
    title varchar(255) not null,
    date TIMESTAMP DEFAULT NULL,
    duration bigint not null,
    timeshift bigint not null,
    description varchar(255) not null
);