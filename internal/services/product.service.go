package services

import (
	"regexp"
	"strings"

	"qonevo-backend/internal/models"
	"qonevo-backend/internal/repositories"
)

type ProductService struct {
	repo      *repositories.ProductRepository
	s3Service *S3Service
}

func NewProductService(
	repo *repositories.ProductRepository,
	s3Service *S3Service,
) *ProductService {

	return &ProductService{
		repo:      repo,
		s3Service: s3Service,
	}
}

// =====================================
// Create Product
// =====================================

func (s *ProductService) CreateProduct(
	product *models.Product,
) (int64, error) {

	product.Slug = generateSlug(
		product.Name,
	)

	// default active
	if !product.IsActive {
		product.IsActive = true
	}

	return s.repo.CreateProduct(
		product,
	)
}

// =====================================
// Create Product Specification
// =====================================

func (s *ProductService) CreateProductSpecification(
	specification *models.ProductSpecification,
) error {

	return s.repo.CreateProductSpecification(
		specification,
	)
}

// =====================================
// Create Product Image
// =====================================

func (s *ProductService) CreateProductImage(
	image *models.ProductImage,
) error {

	return s.repo.CreateProductImage(
		image,
	)
}

// =====================================
// Get All Products
// =====================================

// func (s *ProductService) GetProducts() (
// 	[]models.Product,
// 	error,
// ) {

// 	products, err := s.repo.GetProducts()

// 	if err != nil {
// 		return nil, err
// 	}

// 	// =====================================
// 	// Load images for each product
// 	// =====================================

// 	for i := range products {

// 		images, err := s.repo.GetProductImagesByProductID(
// 			products[i].ID,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		products[i].Images = images

// 		// =====================================
// 		// Set Thumbnail
// 		// =====================================

// 		for _, image := range images {

// 			if image.IsPrimary {

// 				products[i].Thumbnail =
// 					&image.ImageURL

// 				break
// 			}
// 		}

// 		// fallback thumbnail
// 		if products[i].Thumbnail == nil &&
// 			len(images) > 0 {

// 			products[i].Thumbnail =
// 				&images[0].ImageURL
// 		}
// 	}

// 	return products, nil
// }

func (s *ProductService) GetProducts() (
	[]models.Product,
	error,
) {

	products, err := s.repo.GetProducts()

	if err != nil {
		return nil, err
	}

	for i := range products {

		// =========================
		// LOAD IMAGES
		// =========================

		images, err := s.repo.GetProductImagesByProductID(
			products[i].ID,
		)

		if err != nil {
			return nil, err
		}

		products[i].Images = images

		// =========================
		// THUMBNAIL
		// =========================

		for _, image := range images {

			if image.IsPrimary {

				products[i].Thumbnail =
					&image.ImageURL

				break
			}
		}

		if products[i].Thumbnail == nil &&
			len(images) > 0 {

			products[i].Thumbnail =
				&images[0].ImageURL
		}

		// =========================
		// LOAD SPECIFICATIONS
		// =========================

		specifications, err :=
			s.repo.GetProductSpecificationsByProductID(
				products[i].ID,
			)

		if err != nil {
			return nil, err
		}

		products[i].Specifications =
			specifications
	}

	return products, nil
}

// =====================================
// Get Product By ID
// =====================================

func (s *ProductService) GetProductByID(
	id int64,
) (*models.Product, error) {

	product, err := s.repo.GetProductByID(
		id,
	)

	if err != nil || product == nil {
		return product, err
	}

	// =====================================
	// Load Product Images
	// =====================================

	images, err := s.repo.GetProductImagesByProductID(
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	product.Images = images

	// =====================================
	// Set Thumbnail
	// =====================================

	for _, image := range images {

		if image.IsPrimary {

			product.Thumbnail =
				&image.ImageURL

			break
		}
	}

	// fallback thumbnail
	if product.Thumbnail == nil &&
		len(images) > 0 {

		product.Thumbnail =
			&images[0].ImageURL
	}

	// =====================================
	// Load Specifications
	// =====================================

	specifications, err := s.repo.GetProductSpecifications(
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	product.Specifications = specifications

	return product, nil
}

// =====================================
// Update Product
// =====================================

func (s *ProductService) UpdateProduct(
	product *models.Product,
) error {

	product.Slug = generateSlug(
		product.Name,
	)

	return s.repo.UpdateProduct(
		product,
	)
}

// =====================================
// Delete Product
// =====================================

func (s *ProductService) DeleteProduct(
	id int64,
) error {

	return s.repo.DeleteProduct(
		id,
	)
}

// =====================================
// Delete Product Image
// =====================================

func (s *ProductService) DeleteProductImage(
	id int64,
) error {

	return s.repo.DeleteProductImage(
		id,
	)
}

// =====================================
// Count Products
// =====================================

func (s *ProductService) CountProducts() (
	int,
	error,
) {

	return s.repo.CountProducts()
}

// =====================================
// Generate Slug
// =====================================

func generateSlug(
	name string,
) string {

	slug := strings.ToLower(name)

	slug = strings.TrimSpace(slug)

	slug = strings.ReplaceAll(
		slug,
		" ",
		"-",
	)

	reg := regexp.MustCompile(
		`[^a-z0-9\-]`,
	)

	slug = reg.ReplaceAllString(
		slug,
		"",
	)

	reg2 := regexp.MustCompile(
		`-+`,
	)

	slug = reg2.ReplaceAllString(
		slug,
		"-",
	)

	slug = strings.Trim(
		slug,
		"-",
	)

	return slug
}
