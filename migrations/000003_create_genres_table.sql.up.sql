CREATE TABLE IF NOT EXISTS genres(
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at timestamp(0) with time zone NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL,
    deleted_at timestamp(0) with time zone NULL
);

CREATE TABLE genres_categories(
    genre_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    constraint idx_genre_category unique (genre_id, category_id),
    constraint fk_genre_id foreign key (genre_id) references genres (id) on delete cascade,
    constraint fk_category_id foreign key (category_id) references categories (id) on delete cascade
);



   