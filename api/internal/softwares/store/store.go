package store

import (
	"fmt"
	"quepc/api/internal/softwares/model"
)

type Store struct {
	//GORM db
}

func New() *Store {
	return &Store{}
}

func (s *Store) List() (model.Softwares, error) {
	softwares := hardcodedData()
	return softwares, nil
}

func (s *Store) Create(software *model.Software) (*model.Software, error) {
	return software, nil
}

func (s *Store) Read(id string) (*model.Software, error) {
	for _, sw := range hardcodedData() {
		if sw.ID == id {
			return sw, nil
		}
	}
	return nil, fmt.Errorf("software con id %s no encontrado", id)
}

func (s *Store) Update(software *model.Software) (int64, error) {
	return 0, nil
}

func (s *Store) Delete(id string) (int64, error) {
	return 0, nil
}

func hardcodedData() []*model.Software {
	c := make([]*model.Software, 0)
	c = append(c, &model.Software{
		ID:       "1",
		Nombre:   "AutoCAD",
		Empresa:  "Autodesk",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "2",
		Nombre:   "Photoshop",
		Empresa:  "Adobe",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "3",
		Nombre:   "Illustrator",
		Empresa:  "Adobe",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "4",
		Nombre:   "Premiere Pro",
		Empresa:  "Adobe",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "5",
		Nombre:   "After Effects",
		Empresa:  "Adobe",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "6",
		Nombre:   "Visual Studio",
		Empresa:  "Microsoft",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "7",
		Nombre:   "IntelliJ IDEA",
		Empresa:  "JetBrains",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "8",
		Nombre:   "Eclipse",
		Empresa:  "Eclipse Foundation",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "9",
		Nombre:   "NetBeans",
		Empresa:  "Apache",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "10",
		Nombre:   "Unity",
		Empresa:  "Unity Technologies",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "11",
		Nombre:   "Unreal Engine",
		Empresa:  "Epic Games",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "12",
		Nombre:   "Blender",
		Empresa:  "Blender Foundation",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "13",
		Nombre:   "Maya",
		Empresa:  "Autodesk",
		ImageURL: "",
	})
	c = append(c, &model.Software{
		ID:       "14",
		Nombre:   "3ds Max",
		Empresa:  "Autodesk",
		ImageURL: "",
	})
	return c
}
