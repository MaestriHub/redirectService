package migrations_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"redirectServer/internal/database/migrations"
	"redirectServer/internal/database/models"
)

// Helper function to setup a test database
func setupTestDB() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func assertTableAndColumns(t *testing.T, db *gorm.DB, tableName string, expectedColumns map[string]string) {
	// Проверка существования таблицы
	if !db.Migrator().HasTable(tableName) {
		t.Errorf("Expected table %s to exist, but it does not", tableName)
		return
	}

	// Получаем информацию о столбцах таблицы
	rows, err := db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", tableName)).Rows()
	if err != nil {
		t.Fatalf("Failed to get table info for %s: %v", tableName, err)
	}
	defer rows.Close()

	// Создаем карту для хранения реальных типов данных столбцов
	actualColumns := make(map[string]string)

	// Сканируем все строки результата
	for rows.Next() {
		var cid int
		var name, colType string
		var notnull bool
		var dfltValue sql.NullString
		var pk int

		err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk)
		if err != nil {
			t.Fatalf("Failed to scan row for table %s: %v", tableName, err)
		}

		// Сохраняем имя столбца и его тип данных
		actualColumns[name] = colType
	}

	// Проверяем каждый ожидаемый столбец
	for columnName, columnType := range expectedColumns {
		// Проверяем существование столбца
		actualType, exists := actualColumns[columnName]
		if !exists {
			t.Errorf("Expected column %s in table %s to exist, but it does not", columnName, tableName)
			continue
		}

		// Проверяем тип данных столбца
		if actualType != columnType {
			t.Errorf("Expected column %s in table %s to have type %s, but got %s", columnName, tableName, columnType, actualType)
		}
	}
}

// Test Migration Up
func TestMigration1_Up(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	migration := migrations.NewMigration1(db)

	err = migration.Up()
	if err != nil {
		t.Fatalf("Migration Up failed: %v", err)
	}

	if !db.Migrator().HasTable(&models.DirectLink{}) {
		t.Errorf("Expected table DirectLink to exist, but it does not")
	}
	if !db.Migrator().HasTable(&models.Fingerprint{}) {
		t.Errorf("Expected table Fingerprint to exist, but it does not")
	}
}

// Test Migration Down
func TestMigration1_Down(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	migration := migrations.NewMigration1(db)

	err = migration.Up()
	if err != nil {
		t.Fatalf("Migration Up failed: %v", err)
	}

	err = migration.Down()
	if err != nil {
		t.Fatalf("Migration Down failed: %v", err)
	}

	if db.Migrator().HasTable(&models.DirectLink{}) {
		t.Errorf("Expected table DirectLink to be dropped, but it still exists")
	}
	if db.Migrator().HasTable(&models.Fingerprint{}) {
		t.Errorf("Expected table Fingerprint to be dropped, but it still exists")
	}
}

func TestMigration1_Up_NoMigrations(t *testing.T) {
	newDB, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup new test database: %v", err)
	}

	if newDB.Migrator().HasTable(&models.DirectLink{}) {
		t.Errorf("Expected table DirectLink to be dropped, but it still exists")
	}
	if newDB.Migrator().HasTable(&models.Fingerprint{}) {
		t.Errorf("Expected table Fingerprint to be dropped, but it still exists")
	}
}

// Test Error Handling in Migration Up
func TestMigration1_Up_ErrorHandling(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	migration := migrations.NewMigration1(db)

	err = migration.Up()

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

// Test Error Handling in Migration Down
func TestMigration1_Down_ErrorHandling(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	migration := migrations.NewMigration1(db)

	err = migration.Down()

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

// test created columns
func TestMigration1_Up_Types(t *testing.T) {
	// Arrange
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	migration := migrations.NewMigration1(db)

	// Act
	err = migration.Up()
	if err != nil {
		t.Fatalf("Migration Up failed: %v", err)
	}

	// Assert
	assertTableAndColumns(t, db, "direct_links", map[string]string{
		"id":         "uuid",
		"nano_id":    "TEXT",
		"clicks":     "INTEGER",
		"payload":    "jsonb",
		"event":      "TEXT",
		"created_at": "datetime",
		"updated_at": "datetime",
		"deleted_at": "datetime",
	})

	assertTableAndColumns(t, db, "fingerprints", map[string]string{
		"id":              "uuid",
		"ip":              "TEXT",
		"user_agent":      "TEXT",
		"language":        "TEXT",
		"languages":       "TEXT",
		"cores":           "INTEGER",
		"memory":          "INTEGER",
		"screen_width":    "INTEGER",
		"screen_height":   "INTEGER",
		"color_depth":     "INTEGER",
		"pixel_ratio":     "REAL",
		"viewport_width":  "INTEGER",
		"viewport_height": "INTEGER",
		"time_zone":       "TEXT",
		"link_id":         "TEXT",
		"created_at":      "datetime",
		"updated_at":      "datetime",
		"deleted_at":      "datetime",
	})
}
