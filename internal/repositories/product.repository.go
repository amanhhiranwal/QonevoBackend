package repositories

import (
	"database/sql"

	"qonevo-backend/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(
	db *sql.DB,
) *ProductRepository {

	return &ProductRepository{
		db: db,
	}
}

// =====================================
// Create Product
// =====================================

func (r *ProductRepository) CreateProduct(
	product *models.Product,
) (int64, error) {

	query := `
	INSERT INTO products (
		name,
		slug,
		subheading,
		size,
		chipset,
		storage,
		resolution,
		google_integration
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id
	`

	var productID int64

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Slug,
		product.Subheading,
		product.Size,
		product.Chipset,
		product.Storage,
		product.Resolution,
		product.GoogleIntegration,
	).Scan(&productID)

	if err != nil {
		return 0, err
	}

	return productID, nil
}

// =====================================
// Create Product Image
// =====================================

func (r *ProductRepository) CreateProductImage(
	image *models.ProductImage,
) error {

	query := `
	INSERT INTO product_images (
		product_id,
		image_url,
		is_primary
	)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		query,
		image.ProductID,
		image.ImageURL,
		image.IsPrimary,
	)

	return err
}

// =====================================
// Get Products
// =====================================

func (r *ProductRepository) GetProducts() (
	[]models.Product,
	error,
) {

	query := `
	SELECT
		p.id,
		p.name,
		p.slug,
		p.subheading,
		p.size,
		p.chipset,
		p.storage,
		p.resolution,
		p.google_integration,
		p.is_active,
		pi.image_url,
		p.created_at,
		p.updated_at
	FROM products p
	LEFT JOIN product_images pi
		ON pi.product_id = p.id
		AND pi.is_primary = true
	ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Subheading,
			&product.Size,
			&product.Chipset,
			&product.Storage,
			&product.Resolution,
			&product.GoogleIntegration,
			&product.IsActive,
			&product.Thumbnail,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(
			products,
			product,
		)
	}

	// check row iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// =====================================
// Get Product By ID
// =====================================

// func (r *ProductRepository) GetProductByID(
// 	id int64,
// ) (*models.Product, error) {

// 	query := `
// 	SELECT
// 		p.id,
// 		p.name,
// 		p.slug,
// 		p.subheading,
// 		p.size,
// 		p.chipset,
// 		p.storage,
// 		p.resolution,
// 		p.google_integration,
// 		p.is_active,
// 		pi.image_url,
// 		p.created_at,
// 		p.updated_at
// 	FROM products p
// 	LEFT JOIN product_images pi
// 		ON pi.product_id = p.id
// 		AND pi.is_primary = true
// 	WHERE p.id = $1
// 	`

// 	var product models.Product

// 	err := r.db.QueryRow(
// 		query,
// 		id,
// 	).Scan(
// 		&product.ID,
// 		&product.Name,
// 		&product.Slug,
// 		&product.Subheading,
// 		&product.Size,
// 		&product.Chipset,
// 		&product.Storage,
// 		&product.Resolution,
// 		&product.GoogleIntegration,
// 		&product.IsActive,
// 		&product.Thumbnail,
// 		&product.CreatedAt,
// 		&product.UpdatedAt,
// 	)

// 	if err != nil {

// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}

// 		return nil, err
// 	}

// 	return &product, nil
// }
// =====================================
// Get Product By ID
// =====================================

func (r *ProductRepository) GetProductByID(
	id int64,
) (*models.Product, error) {

	query := `
	SELECT
		p.id,
		p.name,
		p.slug,
		p.subheading,
		p.size,
		p.chipset,
		p.storage,
		p.resolution,
		p.google_integration,
		p.is_active,
		pi.image_url,
		p.created_at,
		p.updated_at
	FROM products p
	LEFT JOIN product_images pi
		ON pi.product_id = p.id
		AND pi.is_primary = true
	WHERE p.id = $1
	`

	var product models.Product

	err := r.db.QueryRow(
		query,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Subheading,
		&product.Size,
		&product.Chipset,
		&product.Storage,
		&product.Resolution,
		&product.GoogleIntegration,
		&product.IsActive,
		&product.Thumbnail,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	// =====================================
	// Load All Images
	// =====================================

	imageQuery := `
	SELECT
		id,
		product_id,
		image_url,
		is_primary
	FROM product_images
	WHERE product_id = $1
	ORDER BY id ASC
	`

	rows, err := r.db.Query(
		imageQuery,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var images []models.ProductImage

	for rows.Next() {

		var image models.ProductImage

		err := rows.Scan(
			&image.ID,
			&image.ProductID,
			&image.ImageURL,
			&image.IsPrimary,
		)

		if err != nil {
			return nil, err
		}

		images = append(
			images,
			image,
		)
	}

	product.Images = images

	return &product, nil
}

// =====================================
// Update Product
// =====================================

func (r *ProductRepository) DeleteProductImage(
	id int64,
) error {

	query := `
	DELETE FROM product_images
	WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		id,
	)

	return err
}

func (r *ProductRepository) UpdateProduct(
	product *models.Product,
) error {

	query := `
	UPDATE products
	SET
		name = $1,
		slug = $2,
		subheading = $3,
		size = $4,
		chipset = $5,
		storage = $6,
		resolution = $7,
		google_integration = $8,
		is_active = $9,
		updated_at = NOW()
	WHERE id = $10
	`

	_, err := r.db.Exec(
		query,
		product.Name,
		product.Slug,
		product.Subheading,
		product.Size,
		product.Chipset,
		product.Storage,
		product.Resolution,
		product.GoogleIntegration,
		product.IsActive,
		product.ID,
	)

	return err
}

// =====================================
// Delete Product
// =====================================

func (r *ProductRepository) DeleteProduct(
	id int64,
) error {

	query := `
	DELETE FROM products
	WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		id,
	)

	return err
}

// =====================================
// Count Products
// =====================================

func (r *ProductRepository) CountProducts() (
	int,
	error,
) {

	query := `
	SELECT COUNT(*)
	FROM products
	`

	var count int

	err := r.db.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// =====================================
// GET PRODUCT IMAGES
// =====================================

func (r *ProductRepository) GetProductImagesByProductID(
	productID int64,
) ([]models.ProductImage, error) {

	query := `
	SELECT
		id,
		product_id,
		image_url,
		is_primary,
		created_at
	FROM product_images
	WHERE product_id = $1
	ORDER BY is_primary DESC, created_at ASC
	`

	rows, err := r.db.Query(
		query,
		productID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var images []models.ProductImage

	for rows.Next() {

		var image models.ProductImage

		err := rows.Scan(
			&image.ID,
			&image.ProductID,
			&image.ImageURL,
			&image.IsPrimary,
			&image.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		images = append(
			images,
			image,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
