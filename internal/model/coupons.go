package model

import (
	"time"

	"github.com/lib/pq"
)

type UsageType string
type DiscountType string
type DiscountTarget string

const (
	UsageTypeCart    UsageType = "cart_wise"
	UsageTypeProduct UsageType = "product_wise"
	UsageTypeBxGy    UsageType = "Bxgy"
)
const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFlat       DiscountType = "flat"
)

type Coupon struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	Code                 string         `json:"code" gorm:"uniqueIndex"`
	DiscountType         DiscountType   `json:"discount_type"`
	DiscountValue        float64        `json:"discount_value"`
	MinOrderValue        float64        `json:"min_order_value"`
	MaxDiscount          float64        `json:"max_discount"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	UsageLimit           int            `json:"usage_limit"`
	UsageType            UsageType      `json:"usage_type"`
	UsageCount           int            `json:"usage_count" gorm:"default:0"`
	IsActive             bool           `json:"is_active" gorm:"default:true"`
	ApplicableCategories pq.StringArray `json:"applicable_items" gorm:"type:text;serializer:json"`
	ApplicableProductIDs pq.StringArray `json:"applicable_product_ids" gorm:"type:text;serializer:json"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}
type CreateCouponRequest struct {
	Code                 string         `json:"code"`
	DiscountType         DiscountType   `json:"discount_type"`
	DiscountValue        float64        `json:"discount_value"`
	MinOrderValue        float64        `json:"min_order_value"`
	MaxDiscount          float64        `json:"max_discount"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	UsageLimit           int            `json:"usage_limit"`
	IsActive             bool           `json:"is_active"`
	ApplicableCategories pq.StringArray `json:"applicable_items"`
	ApplicableProductIDs pq.StringArray `json:"applicable_product_ids"`
	UsageType            UsageType      `json:"usage_type"`
}

type CartItem struct {
	ID       string `json:"id"`
	Price    int    `json:"price"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
}

type Cart struct {
	Items []CartItem `json:"cart_items"`
	Total float64    `json:"total"`
}
