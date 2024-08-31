CREATE TABLE videos_video_media (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    checksum VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    encoded_path VARCHAR(500) NOT NULL,
    media_status VARCHAR(50) NOT NULL
);

CREATE TABLE videos_image_media (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    checksum VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL
);

CREATE TABLE videos (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    year_launched SMALLINT NOT NULL,
    opened BOOLEAN NOT NULL DEFAULT FALSE,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    rating VARCHAR(5),
    duration DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    video_id BIGINT NULL,
    trailer_id BIGINT NULL,
    banner_id BIGINT NULL,
    thumbnail_id BIGINT NULL,
    thumbnail_half_id BIGINT NULL,
    CONSTRAINT fk_v_video_id FOREIGN KEY (video_id) REFERENCES videos_video_media (id) ON DELETE CASCADE,
    CONSTRAINT fk_v_trailer_id FOREIGN KEY (trailer_id) REFERENCES videos_video_media (id) ON DELETE CASCADE,
    CONSTRAINT fk_v_banner_id FOREIGN KEY (banner_id) REFERENCES videos_image_media (id) ON DELETE CASCADE,
    CONSTRAINT fk_v_thumb_id FOREIGN KEY (thumbnail_id) REFERENCES videos_image_media (id) ON DELETE CASCADE,
    CONSTRAINT fk_v_thumb_half_id FOREIGN KEY (thumbnail_half_id) REFERENCES videos_image_media (id) ON DELETE CASCADE
);

CREATE TABLE videos_categories (
    video_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    CONSTRAINT idx_vcs_video_category UNIQUE (video_id, category_id),
    CONSTRAINT fk_vcs_video_id FOREIGN KEY (video_id) REFERENCES videos (id),
    CONSTRAINT fk_vcs_category_id FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE videos_genres (
    video_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,
    CONSTRAINT idx_vgs_video_genre UNIQUE (video_id, genre_id),
    CONSTRAINT fk_vgs_video_id FOREIGN KEY (video_id) REFERENCES videos (id),
    CONSTRAINT fk_vgs_genre_id FOREIGN KEY (genre_id) REFERENCES genres (id)
);

CREATE TABLE videos_cast_members (
    video_id BIGINT NOT NULL,
    cast_member_id BIGINT NOT NULL,
    CONSTRAINT idx_vcms_video_member UNIQUE (video_id, cast_member_id),
    CONSTRAINT fk_vcms_video_id FOREIGN KEY (video_id) REFERENCES videos (id),
    CONSTRAINT fk_vcms_genre_id FOREIGN KEY (cast_member_id) REFERENCES cast_members (id)
);