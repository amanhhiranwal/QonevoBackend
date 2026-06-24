package repositories

import (
	"context"
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

	// query := `
	// INSERT INTO products (
	// 	name,
	// 	slug,
	// 	subheading,
	// 	google_integration
	// )
	// VALUES ($1, $2, $3, $4)
	// RETURNING id
	// `
	query := `
	INSERT INTO products (
	name,
	slug,
	subheading,
	google_integration,
	product_type
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	var productID int64

	// err := r.db.QueryRow(
	// 	query,
	// 	product.Name,
	// 	product.Slug,
	// 	product.Subheading,
	// 	product.GoogleIntegration,
	// ).Scan(&productID)

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Slug,
		product.Subheading,
		product.GoogleIntegration,
		product.ProductType,
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
// Create Product Specification
// =====================================

func (r *ProductRepository) CreateProductSpecification(
	specification *models.ProductSpecification,
) error {

	query := `
	INSERT INTO product_specifications (
		product_id,
		category,
		spec_key,
		spec_value
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		specification.ProductID,
		specification.Category,
		specification.SpecKey,
		specification.SpecValue,
	)

	return err
}

// =====================================
// Get Products
// =====================================

// func (r *ProductRepository) GetProducts() (
// 	[]models.Product,
// 	error,
// ) {

// 	query := `
// 	SELECT
// 		p.id,
// 		p.product_type,
// 		p.name,
// 		p.slug,
// 		p.subheading,
// 		p.google_integration,
// 		p.is_active,
// 		pi.image_url,
// 		p.created_at,
// 		p.updated_at
// 	FROM products p
// 	LEFT JOIN product_images pi
// 		ON pi.product_id = p.id
// 		AND pi.is_primary = true
// 	ORDER BY p.created_at DESC
// 	`

// 	rows, err := r.db.Query(query)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var products []models.Product

// 	for rows.Next() {

// 		var product models.Product

// 		// err := rows.Scan(
// 		// 	&product.ID,
// 		// 	&product.Name,
// 		// 	&product.Slug,
// 		// 	&product.Subheading,
// 		// 	&product.GoogleIntegration,
// 		// 	&product.IsActive,
// 		// 	&product.Thumbnail,
// 		// 	&product.CreatedAt,
// 		// 	&product.UpdatedAt,
// 		// )
// 		err := rows.Scan(
// 			&product.ID,
// 			&product.ProductType,
// 			&product.Name,
// 			&product.Slug,
// 			&product.Subheading,
// 			&product.GoogleIntegration,
// 			&product.IsActive,
// 			&product.Thumbnail,
// 			&product.CreatedAt,
// 			&product.UpdatedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		products = append(
// 			products,
// 			product,
// 		)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return products, nil
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
		p.product_type,
		p.name,
		p.slug,
		p.subheading,
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
		&product.ProductType,
		&product.Name,
		&product.Slug,
		&product.Subheading,
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
	// Load Images
	// =====================================

	images, err := r.GetProductImagesByProductID(
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	product.Images = images

	// =====================================
	// Load Specifications
	// =====================================

	specifications, err := r.GetProductSpecifications(
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	product.Specifications = specifications

	// log.Printf("%+v", product.Specifications)

	return &product, nil
}

// =====================================
// Get Product Specifications
// =====================================

func (r *ProductRepository) GetProductSpecifications(
	productID int64,
) ([]models.ProductSpecificationCategory, error) {

	query := `
	SELECT
		category,
		spec_key,
		spec_value
	FROM product_specifications
	WHERE product_id = $1
	ORDER BY category ASC, id DESC
	`
	// ORDER BY category ASC, id ASC

	rows, err := r.db.Query(
		query,
		productID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categoryMap := make(
		map[string][]models.ProductSpecification,
	)

	for rows.Next() {

		var category string
		var key string
		var value string

		err := rows.Scan(
			&category,
			&key,
			&value,
		)

		if err != nil {
			return nil, err
		}

		categoryMap[category] = append(

			categoryMap[category],

			models.ProductSpecification{
				Category:  category,
				SpecKey:   key,
				SpecValue: value,
			},
		)
	}

	var categories []models.ProductSpecificationCategory

	for category, items := range categoryMap {

		categories = append(

			categories,

			models.ProductSpecificationCategory{
				Category: category,
				Items:    items,
			},
		)
	}

	return categories, nil
}

func (r *ProductRepository) GetProductSpecsMap(
	productID int64,
) (map[string]string, error) {

	query := `
		SELECT
			spec_key,
			spec_value
		FROM product_specifications
		WHERE product_id = $1
	`

	rows, err := r.db.Query(
		query,
		productID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	specs := make(map[string]string)

	for rows.Next() {

		var key string
		var value string

		err := rows.Scan(
			&key,
			&value,
		)

		if err != nil {
			return nil, err
		}

		specs[key] = value
	}

	return specs, nil
}

// =====================================
// Delete Product Image
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

func (r *ProductRepository) GetSpecsMap(
	ctx context.Context,
	productID int64,
) (map[string]string, error) {

	query := `
    SELECT
        key,
        value
    FROM product_specs
    WHERE product_id = $1
    `

	rows, err := r.db.QueryContext(
		ctx,
		query,
		productID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	specs := make(map[string]string)

	for rows.Next() {

		var key string
		var value string

		err := rows.Scan(
			&key,
			&value,
		)

		if err != nil {
			return nil, err
		}

		specs[key] = value
	}

	return specs, nil
}

func (r *ProductRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Product, error) {

	query := `
    SELECT
        id,
        name,
        subheading,
        google_integration,
        is_active
    FROM products
    WHERE id = $1
    `

	var p models.Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Subheading,
		&p.GoogleIntegration,
		&p.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// =====================================
// Update Product
// =====================================

func (r *ProductRepository) UpdateProduct(
	product *models.Product,
) error {

	// query := `
	// UPDATE products
	// SET
	// 	name = $1,
	// 	slug = $2,
	// 	subheading = $3,
	// 	google_integration = $4,
	// 	is_active = $5,
	// 	updated_at = NOW()
	// WHERE id = $6
	// `
	query := `
	UPDATE products
	SET
	name = $1,
	subheading = $2,
	google_integration = $3,
	product_type = $4,
	is_active = $5,
	updated_at = NOW()
	WHERE id = $6
	`

	// _, err := r.db.Exec(
	// 	query,
	// 	product.Name,
	// 	product.Slug,
	// 	product.Subheading,
	// 	product.GoogleIntegration,
	// 	product.IsActive,
	// 	product.ID,
	// )

	_, err := r.db.Exec(
		query,
		product.Name,
		product.Subheading,
		product.GoogleIntegration,
		product.ProductType,
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
// Get Product Images
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

func (r *ProductRepository) GetProductSpecificationsByProductID(
	productID int64,
) ([]models.ProductSpecificationCategory, error) {

	query := `
	SELECT
		id,
		product_id,
		category,
		spec_key,
		spec_value
	FROM product_specifications
	WHERE product_id = $1
	ORDER BY category ASC, id ASC
	`

	rows, err := r.db.Query(
		query,
		productID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categoryMap := make(
		map[string][]models.ProductSpecification,
	)

	for rows.Next() {

		var spec models.ProductSpecification

		err := rows.Scan(
			&spec.ID,
			&spec.ProductID,
			&spec.Category,
			&spec.SpecKey,
			&spec.SpecValue,
		)

		if err != nil {
			return nil, err
		}

		categoryMap[spec.Category] =
			append(
				categoryMap[spec.Category],
				spec,
			)
	}

	var categories []models.ProductSpecificationCategory

	for category, items := range categoryMap {

		categories = append(
			categories,
			models.ProductSpecificationCategory{
				Category: category,
				Items:    items,
			},
		)
	}

	return categories, nil
}

func (r *ProductRepository) GetProducts() (
	[]models.Product,
	error,
) {

	query := `
	SELECT
    p.id,
    p.product_type,
    p.name,
    p.slug,
    p.subheading,
    p.google_integration,
    p.is_active,

    (
        SELECT image_url
        FROM product_images
        WHERE product_id = p.id
        ORDER BY is_primary DESC, created_at ASC
        LIMIT 1
    ) AS image_url,

    p.created_at,
    p.updated_at
FROM products p
ORDER BY p.created_at DESC
	`
	// query := `
	// SELECT
	// 	p.id,
	// 	p.product_type,
	// 	p.name,
	// 	p.slug,
	// 	p.subheading,
	// 	p.google_integration,
	// 	p.is_active,
	// 	pi.image_url,
	// 	p.created_at,
	// 	p.updated_at
	// FROM products p
	// LEFT JOIN product_images pi
	// 	ON pi.product_id = p.id
	// 	AND pi.is_primary = true
	// ORDER BY p.created_at DESC
	// `

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.scanProducts(rows)
}

func (r *ProductRepository) scanProducts(
	rows *sql.Rows,
) ([]models.Product, error) {

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.ProductType,
			&product.Name,
			&product.Slug,
			&product.Subheading,
			&product.GoogleIntegration,
			&product.IsActive,
			&product.Thumbnail,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductsByType(
	productType string,
) ([]models.Product, error) {

	query := `
	SELECT
    p.id,
    p.product_type,
    p.name,
    p.slug,
    p.subheading,
    p.google_integration,
    p.is_active,

    (
        SELECT image_url
        FROM product_images
        WHERE product_id = p.id
        ORDER BY is_primary DESC, created_at ASC
        LIMIT 1
    ) AS image_url,

    p.created_at,
    p.updated_at
FROM products p
WHERE p.product_type = $1
ORDER BY p.created_at DESC
	`
	// query := `
	// SELECT
	// 	p.id,
	// 	p.product_type,
	// 	p.name,
	// 	p.slug,
	// 	p.subheading,
	// 	p.google_integration,
	// 	p.is_active,
	// 	pi.image_url,
	// 	p.created_at,
	// 	p.updated_at
	// FROM products p
	// LEFT JOIN product_images pi
	// 	ON pi.product_id = p.id
	// 	AND pi.is_primary = true
	// WHERE p.product_type = $1
	// ORDER BY p.created_at DESC
	// `

	rows, err := r.db.Query(
		query,
		productType,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.scanProducts(rows)
}
