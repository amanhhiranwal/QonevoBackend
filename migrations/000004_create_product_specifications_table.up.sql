CREATE TABLE IF NOT EXISTS product_specifications (

    id BIGSERIAL PRIMARY KEY,

    product_id BIGINT NOT NULL,

    spec_key VARCHAR(255) NOT NULL,

    spec_value TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_specifications_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);