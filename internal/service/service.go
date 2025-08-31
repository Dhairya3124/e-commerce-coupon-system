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
	GetAllCouponService(ctx context.Context) ([]model.Coupon, error)
	GetCouponByCodeService(ctx context.Context, code string) (*model.Coupon, error)
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
