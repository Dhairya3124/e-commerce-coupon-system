package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/cache"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/model"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/repository"
)

type couponService struct {
	mu    sync.RWMutex
	repo  repository.Coupon
	cache cache.Cache
}
type CouponService interface {
	CreateCouponService(ctx context.Context, coupon *model.Coupon) error
	GetAllCouponService(ctx context.Context) ([]model.Coupon, error)
	GetCouponByCodeService(ctx context.Context, code string) (*model.Coupon, error)
	DeleteCouponByCodeService(ctx context.Context, code string) error
	UpdateCouponByCodeService(ctx context.Context, code string, coupon *model.Coupon) error
	GetApplicableCouponsService(ctx context.Context, cart *model.Cart) ([]model.Coupon, error)
	ApplyCouponService(ctx context.Context, code string, cart *model.Cart) (*model.Cart, error)
}

func NewCouponService(repo repository.Coupon, cache cache.Cache) CouponService {
	return &couponService{repo: repo, cache: cache}
}
func (s *couponService) CreateCouponService(ctx context.Context, coupon *model.Coupon) error {
	return s.repo.Create(ctx, coupon)

}
func (s *couponService) GetAllCouponService(ctx context.Context) ([]model.Coupon, error) {
	return s.repo.GetAllCoupons(ctx)

}
func (s *couponService) GetCouponByCodeService(ctx context.Context, code string) (*model.Coupon, error) {
	return s.repo.GetCouponByCode(ctx, code)
}
func (s *couponService) DeleteCouponByCodeService(ctx context.Context, code string) error {
	return s.repo.DeleteCoupon(ctx, code)
}
func (s *couponService) UpdateCouponByCodeService(ctx context.Context, code string, coupon *model.Coupon) error {
	return s.repo.UpdateCoupon(ctx, code, coupon)
}
func (s *couponService) GetApplicableCouponsService(ctx context.Context, cart *model.Cart) ([]model.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	coupons, err := s.repo.GetAllCoupons(ctx)
	if err != nil {
		return nil, err
	}

	applicableCoupons := make([]model.Coupon, 0)
	now := time.Now()

	for _, coupon := range coupons {
		if !coupon.IsActive ||
			now.Before(coupon.StartDate) ||
			now.After(coupon.EndDate) ||
			cart.Total < coupon.MinOrderValue ||
			coupon.UsageCount >= coupon.UsageLimit {
			continue
		}

		switch coupon.UsageType {
		case model.UsageTypeCart:
			applicableCoupons = append(applicableCoupons, coupon)

		case model.UsageTypeProduct:
			hasApplicableProduct := false

			for _, item := range cart.Items {
				for _, pid := range coupon.ApplicableProductIDs {
					if item.ID == pid {
						hasApplicableProduct = true
						break
					}
				}
				if hasApplicableProduct {
					break
				}
			}

			if !hasApplicableProduct {
				for _, item := range cart.Items {
					for _, cat := range coupon.ApplicableCategories {
						if item.Category == cat {
							hasApplicableProduct = true
							break
						}
					}
					if hasApplicableProduct {
						break
					}
				}
			}

			if hasApplicableProduct {
				applicableCoupons = append(applicableCoupons, coupon)
			}

		case model.UsageTypeBxGy:
			if len(coupon.ApplicableProductIDs) >= 2 {
				xID := coupon.ApplicableProductIDs[0]
				yID := coupon.ApplicableProductIDs[1]
				xCount, yCount := 0, 0

				for _, item := range cart.Items {
					if item.ID == xID {
						xCount += item.Quantity
					}
					if item.ID == yID {
						yCount += item.Quantity
					}
				}

				if xCount > 0 && yCount > 0 {
					applicableCoupons = append(applicableCoupons, coupon)
				}
			}
		}
	}

	return applicableCoupons, nil
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

func (s *couponService) ApplyCouponService(ctx context.Context, code string, cart *model.Cart) (*model.Cart, error) {
	coupon, err := s.repo.GetCouponByCode(ctx, code)
	if err != nil || coupon == nil {
		return nil, fmt.Errorf("coupon not found")
	}

	now := time.Now()
	if !coupon.IsActive ||
		now.Before(coupon.StartDate) ||
		now.After(coupon.EndDate) ||
		cart.Total < coupon.MinOrderValue ||
		coupon.UsageCount >= coupon.UsageLimit {
		return nil, fmt.Errorf("coupon not applicable")
	}

	discountedItems := make([]model.CartItem, len(cart.Items))
	copy(discountedItems, cart.Items)
	totalDiscount := 0.0

	switch coupon.UsageType {
	case model.UsageTypeCart:
		discount := 0.0
		if coupon.DiscountType == model.DiscountPercentage {
			discount = cart.Total * coupon.DiscountValue / 100
		} else if coupon.DiscountType == model.DiscountFlat {
			discount = coupon.DiscountValue
		}
		if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
			discount = coupon.MaxDiscount
		}

		for i, item := range cart.Items {
			ratio := float64(item.Price*item.Quantity) / cart.Total
			itemDiscount := discount * ratio
			discountedItems[i] = model.CartItem{
				ID:       item.ID,
				Price:    item.Price - int(math.Round(itemDiscount/float64(item.Quantity))),
				Category: item.Category,
				Quantity: item.Quantity,
			}
			totalDiscount += itemDiscount
		}

	case model.UsageTypeProduct:
		for i, item := range cart.Items {
			apply := false

			for _, pid := range coupon.ApplicableProductIDs {
				if item.ID == pid {
					apply = true
					break
				}
			}

			if !apply {
				for _, cat := range coupon.ApplicableCategories {
					if item.Category == cat {
						apply = true
						break
					}
				}
			}

			if apply {
				discount := 0.0
				if coupon.DiscountType == model.DiscountPercentage {
					discount = float64(item.Price) * coupon.DiscountValue / 100
				} else if coupon.DiscountType == model.DiscountFlat {
					discount = coupon.DiscountValue
				}
				if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
					discount = coupon.MaxDiscount
				}

				discountedItems[i] = model.CartItem{
					ID:    item.ID,
					Price: item.Price - int(math.Round(discount)),
				}
				totalDiscount += discount
			}
		}

	case model.UsageTypeBxGy:
		if len(coupon.ApplicableProductIDs) >= 2 {
			xID := coupon.ApplicableProductIDs[0]
			yID := coupon.ApplicableProductIDs[1]

			xCount, yCount := 0, 0
			for _, item := range cart.Items {
				if item.ID == xID {
					xCount += item.Quantity
				}
				if item.ID == yID {
					yCount += item.Quantity
				}
			}

			eligibleYItems := min(xCount, yCount)
			appliedDiscounts := 0

			for i, item := range cart.Items {
				if item.ID == yID && appliedDiscounts < eligibleYItems {
					itemsToDiscount := min(item.Quantity, eligibleYItems-appliedDiscounts)

					discount := 0.0
					if coupon.DiscountType == model.DiscountPercentage {
						discount = float64(item.Price) * coupon.DiscountValue / 100
					} else if coupon.DiscountType == model.DiscountFlat {
						discount = coupon.DiscountValue
					}
					if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
						discount = coupon.MaxDiscount
					}

					itemDiscount := discount * float64(itemsToDiscount)
					discountedItems[i] = model.CartItem{
						ID:       item.ID,
						Category: item.Category,
						Quantity: item.Quantity,
						Price:    item.Price - int(math.Round(itemDiscount)),
					}
					totalDiscount += itemDiscount
					appliedDiscounts += itemsToDiscount
				}
			}
		}

	default:
	}

	for i := range discountedItems {
		if discountedItems[i].Price < 0 {
			totalDiscount -= float64(-discountedItems[i].Price)
			discountedItems[i].Price = 0
		}
	}

	updatedCart := &model.Cart{
		Items: discountedItems,
		Total: cart.Total - totalDiscount,
	}

	return updatedCart, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
