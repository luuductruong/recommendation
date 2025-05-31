-- create table product
create table if not exists product
(
    product_id  integer not null,
    name        text,
    price       numeric(15, 2),
    category_id text,
    constraint product_pk
    primary key (product_id)
    );

-- create user history
create table if not exists user_view_history
(
    id         text    not null,
    product_id integer not null,
    user_id    text    not null,
    view_at    timestamp with time zone,
    constraint product_view_pk
    primary key (id)
    );

-- create index
create index if not exists idx_user_view_history_user_product_time
    on user_view_history (user_id asc, product_id asc, view_at desc);

-- create category view history
create table category_view_history
(
    id           text not null
        constraint category_view_history_pk
            primary key,
    category_id  text not null
        constraint category_view_history_uk
            unique,
    total_view   integer,
    last_view_at timestamp with time zone
);
