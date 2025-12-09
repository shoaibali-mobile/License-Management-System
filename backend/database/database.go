package database

import (
	"license-mnm/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("license_mnm.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.SubscriptionPack{},
		&models.Subscription{},
	)
	if err != nil {
		return err
	}

	// Create indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_customers_user_id ON customers(user_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_subscriptions_customer_id ON subscriptions(customer_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_subscriptions_pack_id ON subscriptions(pack_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_subscription_packs_sku ON subscription_packs(sku)")

	return nil
}







