package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}
func (repo * InventoryRepository) NewEquipment(equipment model.Equipment) error {
      _, err := repo.db.Exec("INSERT INTO equipment(name,model,quantity, cost_price, date_received) VALUES(?, ?, ?, ?, ?)", equipment.Name, 
	  equipment.ModelOrMake, equipment.Quantity, equipment.CostPrice, equipment.DateReceived)
	return err
}
func (repo * InventoryRepository) GetEquipments()([]model.Equipment, error ){
	equipments := make([]model.Equipment, 0)
	err := repo.db.Select(&equipments, "SELECT id, name, model, quantity, cost_price, date_received FROM equipment where deleted_at is null")
  	return equipments, err
}

func NewInventoryRepository() InventoryRepository{
	return InventoryRepository{
		db: db.GetConnection(),
	}
}