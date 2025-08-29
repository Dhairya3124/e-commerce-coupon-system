package service

import (
	"context"
	"sync"

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
}

func NewCouponService(repo repository.Coupon, cache cache.Cache) CouponService {
	return &couponService{repo: repo, cache: cache}
}
func (s *couponService) CreateCouponService(ctx context.Context, coupon *model.Coupon) error {
	return s.repo.Create(ctx, coupon)

}
