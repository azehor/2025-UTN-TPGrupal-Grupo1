package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	carrerasModel "quepc/api/internal/carreras/model"
	softwaresModel "quepc/api/internal/softwares/model"
)

var DB *gorm.DB

func InitDB() {
	// Cargar archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, para usar las variables del entorno")
	}

	// Leer variables del entorno
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, name, port, sslmode)

	// Conectar con PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	if err := db.AutoMigrate(&softwaresModel.Software{}, &carrerasModel.Carrera{}); err != nil {
		log.Fatalf("Error al migrar tablas: %v", err)
	}
	DB = db
	log.Println("Conexión exitosa a PostgreSQL")
}
