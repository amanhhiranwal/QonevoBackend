package models

import "time"

type Product struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	Subheading        string     `json:"subheading"`
	Size              string     `json:"size"`
	Chipset           string     `json:"chipset"`
	Storage           string     `json:"storage"`
	Resolution        string     `json:"resolution"`
	GoogleIntegration bool       `json:"google_integration"`
	IsActive          bool       `json:"is_active"`

	// =====================================
	// Images
	// =====================================

	Thumbnail *string `json:"thumbnail"`

	Images []ProductImage `json:"images"`

	// =====================================
	// Relations
	// =====================================

	Features       []ProductFeature       `json:"features"`
	Specifications []ProductSpecification `json:"specifications"`

	// =====================================
	// Timestamps
	// =====================================

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}