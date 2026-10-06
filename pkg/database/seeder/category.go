package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

const seedCreateCategory = `INSERT INTO categories (name, description, slug_category)
VALUES ($1, $2, $3)
  RETURNING *`

type categorySeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *categorySeeder {
	return &categorySeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *categorySeeder) Seed() error {
	categoryNames := []string{
		"Electronics", "Clothing", "Groceries", "Toys", "Home & Kitchen",
		"Books", "Beauty & Health", "Sports & Outdoors", "Automotive", "Furniture",
	}

	categoryDescriptions := []string{
		"Best electronics products", "Latest fashion trends", "Fresh groceries",
		"Fun toys for kids", "Essentials for home & kitchen",
		"Books for all ages", "Beauty and health products",
		"Outdoor sports equipment", "Automotive accessories", "Stylish furniture",
	}

	for i := 0; i < 10; i++ {
		name := categoryNames[i%len(categoryNames)]
		description := categoryDescriptions[i%len(categoryDescriptions)]
		slugCategory := fmt.Sprintf("%s-%d", name, i+1)

		var category db.Category
		if err := r.db.GetContext(r.ctx, &category, seedCreateCategory,
			name,
			ptrString(description),
			ptrString(slugCategory),
		); err != nil {
			r.logger.Error("Failed to create category:", zap.Error(err))
			return err
		}
		r.logger.Debug("Category created:", zap.String("slug", slugCategory))
	}

	r.logger.Info("Category seeding completed successfully.", zap.Int("count", 10))
	return nil
}
