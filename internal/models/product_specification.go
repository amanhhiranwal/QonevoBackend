package models

// type ProductSpecification struct {
// 	ID int64 `json:"id"`

// 	ProductID int64 `json:"product_id"`

// 	Category string `json:"category"`

// 	SpecKey string `json:"key"`

// 	SpecValue string `json:"value"`
// }

type ProductSpecification struct {
	ID int64 `json:"id"`

	ProductID int64 `json:"product_id"`

	Category string `json:"category"`

	SpecKey string `json:"spec_key"`

	SpecValue string `json:"spec_value"`
}

// type ProductSpecificationCategory struct {
// 	Category string `json:"category"`

// 	Items []ProductSpecification `json:"items"`
// }

type ProductSpecificationCategory struct {
	Category string `json:"category"`

	Items []ProductSpecification `json:"items"`
}
