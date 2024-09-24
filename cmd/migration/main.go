package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
)

const configPath = "internal/database/migrations"

func main() {
	if err := configs.Setup(); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Usage: create migration <name> or migrate <direction>")
		return
	}

	command := os.Args[1]
	switch command {
	case "create":
		if len(os.Args) < 4 || os.Args[2] != "migration" {
			fmt.Println("Usage: create migration <name>")
			return
		}
		createMigration(os.Args[3])
	case "migrate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: migrate <direction>")
			return
		}
		migrateDB(os.Args[2])
	case "migrate:fresh":
		dropDB()
		migrateDB("up")
	default:
		fmt.Println("Unknown command. Usage: create migration <name> or migrate <direction>")
	}
}

func dropDB() {
	dbURL := configs.Config.DatabaseURL

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	// Drop the existing database
	_, err = db.Exec("DROP DATABASE IF EXISTS `" + configs.Config.DatabaseName + "`;")
	if err != nil {
		fmt.Println("Error dropping the database:", err)
		return
	}

	// Recreate the database
	_, err = db.Exec("CREATE DATABASE `" + configs.Config.DatabaseName + "`;")
	if err != nil {
		fmt.Println("Error creating the database:", err)
		return
	}

	fmt.Println("Database dropped and recreated successfully")
}



func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	upMigrationPath := fmt.Sprintf("%s/%s_%s.up.sql", configPath, timestamp, name)
	downMigrationPath := fmt.Sprintf("%s/%s_%s.down.sql", configPath, timestamp, name)

	if err := createFile(upMigrationPath); err != nil {
		fmt.Println("Error creating up migration file:", err)
		return
	}
	fmt.Println("Up migration file created successfully at:", upMigrationPath)

	if err := createFile(downMigrationPath); err != nil {
		fmt.Println("Error creating down migration file:", err)
		return
	}
	fmt.Println("Down migration file created successfully at:", downMigrationPath)
}

func createFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func migrateDB(direction string) {
	dbURL := configs.Config.DatabaseURL

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Println("Error creating MySQL driver instance:", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+configPath, "mysql", driver)
	if err != nil {
		fmt.Println("Error creating migration instance:", err)
		return
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Error migrating up:", err)
			return
		}
		fmt.Println("Migration up complete")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Error migrating down:", err)
			return
		}
		fmt.Println("Migration down complete")
	default:
		fmt.Println("Unknown migration direction. Use 'up' or 'down'")
	}
}
