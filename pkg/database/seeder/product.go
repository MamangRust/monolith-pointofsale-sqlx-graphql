package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const seedGetCategories = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const seedCreateProduct = `INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, slug_product, image_product, barcode)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *`

type productSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	var merchants []*db.GetMerchantsRow
	if err := r.db.SelectContext(r.ctx, &merchants, seedGetMerchants, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get merchants:", zap.Any("error", err))
		return err
	}

	var categories []*db.GetCategoriesRow
	if err := r.db.SelectContext(r.ctx, &categories, seedGetCategories, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get categories:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(categories) == 0 {
		r.logger.Error("No merchants or categories found, skipping seeding")
		return nil
	}

	productNames := []string{
		"Smartphone", "Laptop", "Wireless Earbuds", "Gaming Mouse", "Mechanical Keyboard",
		"Smartwatch", "Power Bank", "Bluetooth Speaker", "External Hard Drive", "Monitor",
	}
	brands := []string{"Samsung", "Apple", "Sony", "Logitech", "Razer", "Xiaomi", "HP", "Dell", "Acer", "Asus"}
	images := []string{
		"image1.jpg", "image2.jpg", "image3.jpg", "image4.jpg", "image5.jpg",
		"image6.jpg", "image7.jpg", "image8.jpg", "image9.jpg", "image10.jpg",
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		category := categories[rand.Intn(len(categories))]
		name := productNames[rand.Intn(len(productNames))]
		brand := brands[rand.Intn(len(brands))]
		price := int32(rand.Intn(5000000) + 50000)
		countInStock := int32(rand.Intn(100) + 1)
		weight := int32(rand.Intn(5000) + 100)
		slug := fmt.Sprintf("%s-%d", name, rand.Intn(1000))
		image := images[rand.Intn(len(images))]
		barcode := fmt.Sprintf("BC-%d", rand.Intn(9999999))

		var product db.Product
		if err := r.db.GetContext(r.ctx, &product, seedCreateProduct,
			merchant.MerchantID,
			category.CategoryID,
			name,
			ptrString(fmt.Sprintf("Description for %s", name)),
			price,
			countInStock,
			ptrString(brand),
			ptrInt32(weight),
			ptrString(slug),
			ptrString(image),
			ptrString(barcode),
		); err != nil {
			r.logger.Error("Failed to create product:", zap.Error(err))
			return err
		}

	}

	r.logger.Info("Product seeding completed successfully.", zap.Int("count", 10))
	return nil
}
