package models

type ProductFeature struct {
	ID int64 `json:"id"`

	ProductID int64 `json:"product_id"`

	FeatureName string `json:"feature_name"`
}
