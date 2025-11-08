package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
)

func ValidarUUID(uuid string) bool {
	var validador = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return validador.MatchString(uuid)
}

func GuardarImagen(file multipart.File, fileheader *multipart.FileHeader, endpoint string) (string, error) {
	defer file.Close()
	imgDirPath := filepath.Join(".", "static/images", endpoint)
	err := os.MkdirAll(imgDirPath, os.ModePerm)
	if err != nil {
		return "", errors.New("Error al crear directorio de imagen")
	}
	imgPath := filepath.Join(imgDirPath, fileheader.Filename)
	imagen, err := os.Create(imgPath)
	if err != nil {
		return "", errors.New("Error al crear la imagen")
	}
	defer imagen.Close()
	if _, err := io.Copy(imagen, file); err != nil {
		return "", errors.New("Error al copiar la imagen enviada a local")
	}
	return filepath.Clean(imgPath), nil
}

func EliminarImagen(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func CompletarImageURL(host string, imageURL string) string {
	return "http://" + filepath.Join(host, imageURL)
}
