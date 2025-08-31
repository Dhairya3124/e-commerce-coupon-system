package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/model"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/service"
	"github.com/go-chi/chi"
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
func (h *CouponHandler) GetCouponByCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "couponID")
	res, err := h.service.GetCouponByCodeService(ctx, code)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(res)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

func (h *CouponHandler) GetAllCouponsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.service.GetAllCouponService(ctx)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(res)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CouponHandler) DeleteCouponByCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "couponID")

	err := h.service.DeleteCouponByCodeService(ctx, code)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
func (h *CouponHandler) UpdateCouponByCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "couponID")
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

	err := h.service.UpdateCouponByCodeService(ctx, code, coupon)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
