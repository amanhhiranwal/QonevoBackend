CREATE TABLE IF NOT EXISTS product_images (

    id BIGSERIAL PRIMARY KEY,

    product_id BIGINT NOT NULL,

    image_url TEXT NOT NULL,

    is_primary BOOLEAN DEFAULT false,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_images_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);