# E-Commerce Coupon System

A system for managing and applying discount coupons in an e-commerce platform.

## Features

- Create, update, and delete coupons
- Apply coupons during checkout
- Validate coupon usage and expiration
- Support for percentage and fixed amount discounts
- Limit coupon usage per user and per order
- Admin dashboard for coupon management

## Installation

1. Clone the repository:
  ```bash
  git clone https://github.com/Dhairya3124/e-commerce-coupon-system.git
  ```
2. Initialize the project:
  ```bash
  cd e-commerce-coupon-system
  docker-compose up -d
  ```
3. Install dependencies:
  ```bash
  go mod tidy
  ```
4. Start the application:

  ```bash
  go run cmd/main.go
  ```

## Limitations

- Currently supports only basic coupon types (percentage and fixed amount).
- Only basic cache mechanism is implemented.
- No user authentication or authorization implemented.
- Limited error handling and logging.

