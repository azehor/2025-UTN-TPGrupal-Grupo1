package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidarUUID(uuid string) bool {
	var validador = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return validador.MatchString(uuid)
}

// GuardarImagen ahora guarda todos los archivos directamente bajo ./static
// y devuelve SOLO el nombre de archivo para almacenar en la DB.
// endpoint se conserva para posible categorización futura pero por ahora
// solo se usa para crear subdirectorios opcionales si se desea.
func GuardarImagen(file multipart.File, fileheader *multipart.FileHeader, endpoint string) (string, error) {
	// Ahora la estructura física es: ./static/images/<endpoint>/<filename>
	// endpoint esperado: "carreras", "juegos" o "softwares"
	defer file.Close()
	baseDir := filepath.Join(".", "static", "images", endpoint)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return "", errors.New("error al crear directorio de imagen")
	}
	filename := fileheader.Filename
	imgPath := filepath.Join(baseDir, filename)
	imagen, err := os.Create(imgPath)
	if err != nil {
		return "", errors.New("error al crear la imagen")
	}
	defer imagen.Close()
	if _, err := io.Copy(imagen, file); err != nil {
		return "", errors.New("error al copiar la imagen enviada a local")
	}
	// Guardamos en la DB el path relativo dentro de static: images/<endpoint>/<filename>
	return filepath.ToSlash(filepath.Join("images", endpoint, filename)), nil
}

func EliminarImagen(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func CompletarImageURL(host string, imageURL string) string {
	// Almacenar en DB ejemplos: "images/carreras/carrera-arquitectura.png"
	// Debe exponerse como /static/images/carreras/carrera-arquitectura.png
	base := os.Getenv("API_ASSETS_BASE_URL")
	cleaned := filepath.ToSlash(imageURL)
	cleaned = strings.TrimPrefix(cleaned, "./")
	cleaned = strings.TrimPrefix(cleaned, ".\\")
	// Si vino solo filename (caso legacy), asumir juegos como default? Mejor poner en raiz images.
	if !strings.Contains(cleaned, "/") {
		cleaned = filepath.ToSlash(filepath.Join("images", cleaned))
	}
	// Asegurar prefijo /static/
	if !strings.HasPrefix(cleaned, "static/") && !strings.HasPrefix(cleaned, "/static/") {
		cleaned = filepath.ToSlash(filepath.Join("static", cleaned))
	}
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	if base != "" {
		base = strings.TrimRight(base, "/")
		return base + cleaned
	}
	return "http://" + host + cleaned
}
