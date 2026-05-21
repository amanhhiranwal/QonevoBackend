package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"qonevo-backend/internal/models"
	"qonevo-backend/internal/services"
)

type ProductController struct {
	service   *services.ProductService
	page      *PageController
	s3Service *services.S3Service
}

func NewProductController(
	service *services.ProductService,
	page *PageController,
	s3Service *services.S3Service,
) *ProductController {

	return &ProductController{
		service:   service,
		page:      page,
		s3Service: s3Service,
	}
}

// =====================================
// Products Listing Page
// =====================================

func (c *ProductController) Index(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	products, err := c.service.GetProducts()

	if err != nil {

		log.Printf(
			"failed to fetch products: %v",
			err,
		)

		http.Error(
			w,
			"failed to load products",
			http.StatusInternalServerError,
		)

		return
	}

	c.page.render(
		w,
		"products.html",
		map[string]any{
			"Title":      "Products",
			"IsLoggedIn": true,
			"Products":   products,
		},
	)
}

// =====================================
// Render Create Product Page
// =====================================

func (c *ProductController) CreatePage(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	c.page.render(
		w,
		"create-product.html",
		map[string]any{
			"Title":      "Create Product",
			"IsLoggedIn": true,
		},
	)
}

// =====================================
// Store Product
// =====================================

func (c *ProductController) Store(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================
	// Parse Multipart Form
	// =====================================

	err := r.ParseMultipartForm(
		20 << 20, // 20MB
	)

	if err != nil {

		log.Printf(
			"multipart parse error: %v",
			err,
		)

		http.Error(
			w,
			"failed to parse form",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// Sanitize Inputs
	// =====================================

	name := strings.TrimSpace(
		r.FormValue("name"),
	)

	subheading := strings.TrimSpace(
		r.FormValue("subheading"),
	)

	size := strings.TrimSpace(
		r.FormValue("size"),
	)

	chipset := strings.TrimSpace(
		r.FormValue("chipset"),
	)

	storage := strings.TrimSpace(
		r.FormValue("storage"),
	)

	resolution := strings.TrimSpace(
		r.FormValue("resolution"),
	)

	googleIntegration :=
		r.FormValue("google_integration") == "on"

	// =====================================
	// Validation
	// =====================================

	if name == "" {

		http.Error(
			w,
			"product name is required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// Build Product Model
	// =====================================

	product := &models.Product{
		Name:              name,
		Subheading:        subheading,
		Size:              size,
		Chipset:           chipset,
		Storage:           storage,
		Resolution:        resolution,
		GoogleIntegration: googleIntegration,
		IsActive:          true,
	}

	// =====================================
	// Save Product
	// =====================================

	productID, err := c.service.CreateProduct(
		product,
	)

	if err != nil {

		log.Printf(
			"failed to create product: %v",
			err,
		)

		http.Error(
			w,
			"failed to create product",
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf(
		"product created successfully: %d",
		productID,
	)

	// =====================================
	// Handle Multiple Image Uploads
	// =====================================

	form := r.MultipartForm

	if form != nil {

		files := form.File["images"]

		log.Printf(
			"received %d images",
			len(files),
		)

		for index, fileHeader := range files {

			log.Printf(
				"processing image: %s",
				fileHeader.Filename,
			)

			file, err := fileHeader.Open()

			if err != nil {

				log.Printf(
					"failed to open image: %v",
					err,
				)

				continue
			}

			// =================================
			// Upload To S3
			// =================================

			imageURL, err := c.s3Service.UploadFile(
				file,
				fileHeader,
				"products",
			)

			file.Close()

			if err != nil {

				log.Printf(
					"S3 upload failed: %v",
					err,
				)

				continue
			}

			log.Printf(
				"uploaded image url: %s",
				imageURL,
			)

			// =================================
			// Save Image In Database
			// =================================

			// image := &models.ProductImage{
			// 	ProductID: productID,
			// 	ImageURL:  imageURL,
			// 	SortOrder: index + 1,
			// }

			image := &models.ProductImage{
				ProductID: productID,
				ImageURL:  imageURL,
				IsPrimary: index == 0,
			}

			err = c.service.CreateProductImage(
				image,
			)

			if err != nil {

				log.Printf(
					"failed to save product image: %v",
					err,
				)

				continue
			}
		}
	}

	// =====================================
	// Redirect
	// =====================================

	http.Redirect(
		w,
		r,
		"/products",
		http.StatusSeeOther,
	)
}

// =====================================
// Edit Product Page
// =====================================

func (c *ProductController) EditPage(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================
	// Extract Product ID
	// =====================================

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/products/edit/",
	)

	productID, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		http.Error(
			w,
			"invalid product id",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// Fetch Product
	// =====================================

	product, err := c.service.GetProductByID(
		productID,
	)

	if err != nil {

		log.Printf(
			"failed to fetch product: %v",
			err,
		)

		http.Error(
			w,
			"product not found",
			http.StatusNotFound,
		)

		return
	}

	c.page.render(
		w,
		"edit-product.html",
		map[string]any{
			"Title":      "Edit Product",
			"IsLoggedIn": true,
			"Product":    product,
		},
	)
}

// =====================================
// Delete Product Image
// =====================================

func (c *ProductController) DeleteImage(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================
	// Extract Image ID
	// =====================================

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/products/image/delete/",
	)

	imageID, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		http.Error(
			w,
			"invalid image id",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// Delete Image
	// =====================================

	err = c.service.DeleteProductImage(
		imageID,
	)

	if err != nil {

		log.Printf(
			"failed to delete image: %v",
			err,
		)

		http.Error(
			w,
			"failed to delete image",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================
	// Redirect Back
	// =====================================

	http.Redirect(
		w,
		r,
		r.Referer(),
		http.StatusSeeOther,
	)
}

// =====================================
// Update Product
// =====================================

func (c *ProductController) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/products/update/",
	)

	productID, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		http.Error(
			w,
			"invalid product id",
			http.StatusBadRequest,
		)

		return
	}

	err = r.ParseMultipartForm(
		20 << 20,
	)

	if err != nil {

		http.Error(
			w,
			"failed to parse form",
			http.StatusBadRequest,
		)

		return
	}

	product := &models.Product{
		ID:                productID,
		Name:              strings.TrimSpace(r.FormValue("name")),
		Subheading:        strings.TrimSpace(r.FormValue("subheading")),
		Size:              strings.TrimSpace(r.FormValue("size")),
		Chipset:           strings.TrimSpace(r.FormValue("chipset")),
		Storage:           strings.TrimSpace(r.FormValue("storage")),
		Resolution:        strings.TrimSpace(r.FormValue("resolution")),
		GoogleIntegration: r.FormValue("google_integration") == "on",
		IsActive:          true,
	}

	err = c.service.UpdateProduct(
		product,
	)

	if err != nil {

		log.Printf(
			"failed to update product: %v",
			err,
		)

		http.Error(
			w,
			"failed to update product",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/products",
		http.StatusSeeOther,
	)
}

// =====================================
// Delete Product
// =====================================

func (c *ProductController) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/products/delete/",
	)

	productID, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		http.Error(
			w,
			"invalid product id",
			http.StatusBadRequest,
		)

		return
	}

	err = c.service.DeleteProduct(
		productID,
	)

	if err != nil {

		log.Printf(
			"failed to delete product: %v",
			err,
		)

		http.Error(
			w,
			"failed to delete product",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/products",
		http.StatusSeeOther,
	)
}

// =====================================
// GET PRODUCTS API
// =====================================

func (c *ProductController) GetProductsAPI(
	w http.ResponseWriter,
	r *http.Request,
) {

	// =====================================
	// Query Params
	// =====================================

	limitParam := r.URL.Query().Get(
		"limit",
	)

	// =====================================
	// Fetch Products
	// =====================================

	products, err := c.service.GetProducts()

	if err != nil {

		log.Printf(
			"GET PRODUCTS API ERROR: %v",
			err,
		)

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================
	// Apply Limit
	// =====================================

	if limitParam != "" {

		limit, err := strconv.Atoi(
			limitParam,
		)

		if err == nil &&
			limit > 0 &&
			limit < len(products) {

			products = products[:limit]
		}
	}

	// =====================================
	// Response
	// =====================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	err = json.NewEncoder(w).Encode(
		products,
	)

	if err != nil {

		log.Printf(
			"JSON ENCODE ERROR: %v",
			err,
		)

		http.Error(
			w,
			"failed to encode products",
			http.StatusInternalServerError,
		)

		return
	}
}
