package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/model"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/service"
)

type CouponHandler struct {
	service service.CouponService
}
type ValidationRequest struct {
	CouponCode string     `json:"code"`
	Cart       model.Cart `json:"cart"`
}
type ValidationResponse struct {
	IsValid bool `json:"is_valid"`
}
type GetApplicableCouponsRequest struct {
	Items []model.CartItem `json:"items"`
	Total float64          `json:"total"`
}

func NewCouponHandler(service service.CouponService) *CouponHandler {
	return &CouponHandler{
		service: service,
	}
}

func (h *CouponHandler) CreateCouponHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	coupon := &model.Coupon{
		Code:                 req.Code,
		DiscountType:         req.DiscountType,
		DiscountValue:        req.DiscountValue,
		MinOrderValue:        req.MinOrderValue,
		MaxDiscount:          req.MaxDiscount,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		UsageLimit:           req.UsageLimit,
		IsActive:             req.IsActive,
		ApplicableCategories: req.ApplicableCategories,
		ApplicableProductIDs: req.ApplicableProductIDs,
		UsageType:            req.UsageType,
	}

	err := h.service.CreateCouponService(ctx, coupon)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
