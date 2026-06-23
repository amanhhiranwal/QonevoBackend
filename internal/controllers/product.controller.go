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

	// products, err := c.service.GetProducts()

	products, err := c.service.GetProductsByType(
		"IFP",
	)
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
		20 << 20,
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

	googleIntegration :=
		r.FormValue("google_integration") == "on"

	productType := strings.TrimSpace(
		r.FormValue("product_type"),
	)

	var productTypePtr *string

	if productType != "" {
		productTypePtr = &productType
	}

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
		ProductType:       productTypePtr,
		Subheading:        subheading,
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
	// Save Specifications
	// =====================================

	specCategories := r.Form["spec_category[]"]
	specKeys := r.Form["spec_key[]"]
	specValues := r.Form["spec_value[]"]

	for i := range specKeys {

		if i >= len(specCategories) ||
			i >= len(specValues) {
			continue
		}

		key := strings.TrimSpace(
			specKeys[i],
		)

		value := strings.TrimSpace(
			specValues[i],
		)

		category := strings.TrimSpace(
			specCategories[i],
		)

		if key == "" || value == "" {
			continue
		}

		specification := &models.ProductSpecification{
			ProductID: productID,
			Category:  category,
			SpecKey:   key,
			SpecValue: value,
		}

		err := c.service.CreateProductSpecification(
			specification,
		)

		if err != nil {

			log.Printf(
				"failed to save specification: %v",
				err,
			)
		}
	}

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

	if productType == "CAMERA" {

		http.Redirect(
			w,
			r,
			"/cameras",
			http.StatusSeeOther,
		)

	} else {

		http.Redirect(
			w,
			r,
			"/products",
			http.StatusSeeOther,
		)
	}
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

	productType := strings.TrimSpace(
		r.FormValue("product_type"),
	)

	var productTypePtr *string

	if productType != "" {
		productTypePtr = &productType
	}

	product := &models.Product{
		ID: productID,

		Name: strings.TrimSpace(
			r.FormValue("name"),
		),

		ProductType: productTypePtr,

		Subheading: strings.TrimSpace(
			r.FormValue("subheading"),
		),

		GoogleIntegration: r.FormValue("google_integration") == "on",

		IsActive: true,
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

	if productType == "CAMERA" {

		http.Redirect(
			w,
			r,
			"/cameras",
			http.StatusSeeOther,
		)

	} else {

		http.Redirect(
			w,
			r,
			"/products",
			http.StatusSeeOther,
		)
	}
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

	// Load product before deleting so we know where to redirect
	product, err := c.service.GetProductByID(
		productID,
	)

	if err != nil {

		http.Error(
			w,
			"product not found",
			http.StatusNotFound,
		)

		return
	}

	// Actually delete the product
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

	// Redirect according to product type
	if product.ProductType != nil &&
		*product.ProductType == "CAMERA" {

		http.Redirect(
			w,
			r,
			"/cameras",
			http.StatusSeeOther,
		)

	} else {

		http.Redirect(
			w,
			r,
			"/products",
			http.StatusSeeOther,
		)
	}
}

// =====================================
// GET PRODUCTS API
// =====================================

// func (c *ProductController) GetProductsAPI(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	limitParam := r.URL.Query().Get(
// 		"limit",
// 	)

// 	products, err := c.service.GetProducts()

// 	if err != nil {

// 		log.Printf(
// 			"GET PRODUCTS API ERROR: %v",
// 			err,
// 		)

// 		http.Error(
// 			w,
// 			err.Error(),
// 			http.StatusInternalServerError,
// 		)

// 		return
// 	}

// 	if limitParam != "" {

// 		limit, err := strconv.Atoi(
// 			limitParam,
// 		)

// 		if err == nil &&
// 			limit > 0 &&
// 			limit < len(products) {

// 			products = products[:limit]
// 		}
// 	}

// 	w.Header().Set(
// 		"Content-Type",
// 		"application/json",
// 	)

// 	w.WriteHeader(
// 		http.StatusOK,
// 	)

// 	err = json.NewEncoder(w).Encode(
// 		products,
// 	)

// 	if err != nil {

// 		log.Printf(
// 			"JSON ENCODE ERROR: %v",
// 			err,
// 		)

// 		http.Error(
// 			w,
// 			"failed to encode products",
// 			http.StatusInternalServerError,
// 		)

// 		return
// 	}
// }

func (c *ProductController) GetProductsAPI(
	w http.ResponseWriter,
	r *http.Request,
) {

	limitParam := r.URL.Query().Get(
		"limit",
	)

	productType := strings.ToUpper(
		strings.TrimSpace(
			r.URL.Query().Get("type"),
		),
	)

	var (
		products []models.Product
		err      error
	)

	// Filter by type if provided
	if productType != "" {

		products, err = c.service.GetProductsByType(
			productType,
		)

	} else {

		products, err = c.service.GetProducts()
	}

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

	// Apply limit
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

func (c *ProductController) CreateCameraPage(
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
		"create-camera.html",
		map[string]any{
			"Title":      "Create Camera",
			"IsLoggedIn": true,
			"IsEdit":     false,
			"Specs":      map[string]string{}, // <-- ADD THIS
		},
	)
}

func (c *ProductController) CameraIndex(
	w http.ResponseWriter,
	r *http.Request,
) {

	products, err := c.service.GetProductsByType(
		"CAMERA",
	)

	log.Printf("found %d cameras", len(products))

	if err != nil {
		log.Printf("camera error: %v", err)

		http.Error(
			w,
			"failed to load cameras",
			http.StatusInternalServerError,
		)

		return
	}

	c.page.render(
		w,
		"cameras.html",
		map[string]any{
			"Title":      "Cameras",
			"IsLoggedIn": true,
			"Products":   products,
		},
	)
}

func (p *ProductController) EditCameraPage(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/camera/edit/",
	)

	productID, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		http.NotFound(
			w,
			r,
		)

		return
	}

	product, err := p.service.GetProductByID(
		productID,
	)

	if err != nil {

		http.NotFound(
			w,
			r,
		)

		return
	}

	specs, err := p.service.GetProductSpecsMap(
		productID,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	p.page.render(
		w,
		"create-camera.html",
		map[string]any{
			"Title":      "Edit Camera",
			"IsLoggedIn": true,
			"IsEdit":     true,
			"Product":    product,
			"Specs":      specs,
		},
	)
}

// =====================================
// Upload Additional Product Images
// =====================================

// func (c *ProductController) UploadImages(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	if r.Method != http.MethodPost {

// 		http.Error(
// 			w,
// 			"method not allowed",
// 			http.StatusMethodNotAllowed,
// 		)

// 		return
// 	}

// 	idStr := strings.TrimPrefix(
// 		r.URL.Path,
// 		"/products/images/upload/",
// 	)

// 	productID, err := strconv.ParseInt(
// 		idStr,
// 		10,
// 		64,
// 	)

// 	if err != nil {

// 		http.Error(
// 			w,
// 			"invalid product id",
// 			http.StatusBadRequest,
// 		)

// 		return
// 	}

// 	err = r.ParseMultipartForm(
// 		20 << 20,
// 	)

// 	if err != nil {

// 		http.Error(
// 			w,
// 			"failed to parse form",
// 			http.StatusBadRequest,
// 		)

// 		return
// 	}

// 	files := r.MultipartForm.File["images"]

// 	for _, fileHeader := range files {

// 		file, err := fileHeader.Open()

// 		if err != nil {
// 			continue
// 		}

// 		imageURL, err := c.s3Service.UploadFile(
// 			file,
// 			fileHeader,
// 			"products",
// 		)

// 		file.Close()

// 		if err != nil {
// 			log.Printf(
// 				"upload failed: %v",
// 				err,
// 			)

// 			continue
// 		}

// 		image := &models.ProductImage{
// 			ProductID: productID,
// 			ImageURL:  imageURL,
// 			IsPrimary: false,
// 		}

// 		err = c.service.CreateProductImage(
// 			image,
// 		)

// 		if err != nil {

// 			log.Printf(
// 				"failed to save image: %v",
// 				err,
// 			)
// 		}
// 	}

// 	http.Redirect(
// 		w,
// 		r,
// 		"/products/edit/"+idStr,
// 		http.StatusSeeOther,
// 	)
// }

// =====================================
// Upload Additional Product Images
// =====================================

func (c *ProductController) UploadImages(
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
		"/products/images/upload/",
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

	files := r.MultipartForm.File["images"]

	if len(files) == 0 {

		http.Redirect(
			w,
			r,
			"/products/edit/"+idStr,
			http.StatusSeeOther,
		)

		return
	}

	// =====================================
	// Check if product already has a primary image
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
			"failed to load product",
			http.StatusInternalServerError,
		)

		return
	}

	hasPrimary := false

	for _, img := range product.Images {

		if img.IsPrimary {

			hasPrimary = true
			break
		}
	}

	// =====================================
	// Upload images
	// =====================================

	for index, fileHeader := range files {

		file, err := fileHeader.Open()

		if err != nil {

			log.Printf(
				"failed to open file: %v",
				err,
			)

			continue
		}

		imageURL, err := c.s3Service.UploadFile(
			file,
			fileHeader,
			"products",
		)

		file.Close()

		if err != nil {

			log.Printf(
				"upload failed: %v",
				err,
			)

			continue
		}

		image := &models.ProductImage{
			ProductID: productID,
			ImageURL:  imageURL,

			// first uploaded image becomes primary
			// only if no primary image exists already
			IsPrimary: !hasPrimary && index == 0,
		}

		err = c.service.CreateProductImage(
			image,
		)

		if err != nil {

			log.Printf(
				"failed to save image: %v",
				err,
			)
		}
	}

	http.Redirect(
		w,
		r,
		"/products/edit/"+idStr,
		http.StatusSeeOther,
	)
}
