CREATE TABLE IF NOT EXISTS product_features (

    id BIGSERIAL PRIMARY KEY,

    product_id BIGINT NOT NULL,

    feature_name TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_features_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);