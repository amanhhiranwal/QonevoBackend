package models

import "time"

type Product struct {
	ID int64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Subheading string `json:"subheading"`

	// =====================================
	// KEEP THESE FOR TABLE + LEGACY SUPPORT
	// =====================================

	Size string `json:"size"`

	Chipset string `json:"chipset"`

	Storage string `json:"storage"`

	Resolution string `json:"resolution"`

	// =====================================

	GoogleIntegration bool `json:"google_integration"`

	IsActive bool `json:"is_active"`

	Thumbnail *string `json:"thumbnail,omitempty"`

	Images []ProductImage `json:"images,omitempty"`

	Specifications []ProductSpecificationCategory `json:"specifications,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func (p Product) GetSpecValue(category string, key string) string {

	for _, specCategory := range p.Specifications {

		if specCategory.Category == category {

			for _, item := range specCategory.Items {

				if item.SpecKey == key {
					return item.SpecValue
				}
			}
		}
	}

	return ""
}
