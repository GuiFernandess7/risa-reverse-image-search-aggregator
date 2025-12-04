package database

import (
	"fmt"

	validatorpkg "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CrudGeneric[T any] struct {
	DB *gorm.DB
}

func (c *CrudGeneric[T]) ValidateModel(payload *T) error {
	var validator = validatorpkg.New()
	return validator.Struct(payload)
}

func (c *CrudGeneric[T]) Create(payload *T) error {
	if err := c.ValidateModel(payload); err != nil {
		fmt.Println("Invalid payload:", err)
		return err
	}
	return c.DB.Create(payload).Error
}

func (c *CrudGeneric[T]) Read(field string, value any) ([]T, error) {
	var models []T

	if err := c.DB.Where(fmt.Sprintf("%s = ?", field), value).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (c *CrudGeneric[T]) ReadAll() ([]T, error) {
	var items []T
	if err := c.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (c *CrudGeneric[T]) FindBy(field string, value any) (*T, error) {
	var model T
	if err := c.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (c *CrudGeneric[T]) Update(id any, updated *T) error {
	return c.DB.Model(new(T)).Where("id = ?", id).Updates(updated).Error
}

func (c *CrudGeneric[T]) Delete(id any) error {
	return c.DB.Delete(new(T), id).Error
}
