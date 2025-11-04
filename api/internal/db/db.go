package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	carreraSoftwareModel "quepc/api/internal/carreraSoftware/model"
	carrerasModel "quepc/api/internal/carreras/model"
	discosModel "quepc/api/internal/componentes/discos/model"
	gabinetesModel "quepc/api/internal/componentes/gabinetes/model"
	gpusModel "quepc/api/internal/componentes/gpus/model"
	motherboardsModel "quepc/api/internal/componentes/motherboards/model"
	procesadoresModel "quepc/api/internal/componentes/procesadores/model"
	psusModel "quepc/api/internal/componentes/psus/model"
	ramsModel "quepc/api/internal/componentes/rams/model"
	softwaresModel "quepc/api/internal/softwares/model"
	userModel "quepc/api/internal/user/model"
)

var DB *gorm.DB

func InitDB() {
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

	if err := db.AutoMigrate(&softwaresModel.Software{},
		&carrerasModel.Carrera{},
		&carreraSoftwareModel.CarreraSoftware{},
		&discosModel.Disco{},
		&gabinetesModel.Gabinete{},
		&gpusModel.GPU{},
		&motherboardsModel.Motherboard{},
		&procesadoresModel.Procesador{},
		&psusModel.PSU{},
		&ramsModel.RAM{},
		&userModel.User{}); err != nil {
		log.Fatalf("Error al migrar tablas: %v", err)
	}
	DB = db
	log.Println("Conexión exitosa a PostgreSQL")
}
