package models

import "gorm.io/gorm"

func (pv *ProductVariant) AfterSave(tx *gorm.DB) error {
	return updateProductAggregates(tx, pv.ProductID)
}
func (pv *ProductVariant) AfterDelete(tx *gorm.DB) error {
	return updateProductAggregates(tx, pv.ProductID)
}
func (pv *ProductVariant) BeforeSave(tx *gorm.DB) (err error) {
	pv.IsActive = pv.Stock > 0
	return nil
}
func (pv *ProductVariant) BeforeUpdate(tx *gorm.DB) (err error) {
	pv.IsActive = pv.Stock > 0
	return nil
}
func (pv *ProductVariant) AfterUpdate(tx *gorm.DB) error {
	return updateProductAggregates(tx, pv.ProductID)
}

func updateProductAggregates(tx *gorm.DB, productID uint) error {
	// Ищем активные варианты
	var variants []ProductVariant
	if err := tx.Where("product_id = ? AND is_active = ?", productID, true).Find(&variants).Error; err != nil {
		return err
	}

	// Если нет активных вариантов, то деактивируем продукт и обнуляем его цену и количество
	if len(variants) == 0 {
		return tx.Model(&Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
			"is_active": false,
			"stock":     0,
			"price":     0,
		}).Error
	}

	// суммируем количество всех вариантов и формируем мин стоимость
	var totalStock int
	minPrice := variants[0].Price
	for _, v := range variants {
		if v.Price < minPrice {
			minPrice = v.Price
		}
		totalStock += v.Stock
	}

	// Обновляем продукт с минимальной ценой и общим количеством и активируем его
	return tx.Model(&Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"is_active": true,
		"stock":     totalStock,
		"price":     minPrice,
	}).Error
}

type ProductVariant struct {
	gorm.Model
	Name      string  `json:"name" gorm:"index:idx_product_id_name,unique"`
	Stock     int     `json:"stock" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	IsActive  bool    `json:"is_active" gorm:"default:false"`
	ProductID uint    `json:"product_id" gorm:"not null;index:idx_product_id_name,unique"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
