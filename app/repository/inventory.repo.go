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
func (repo * InventoryRepository) GetStat(equipment model.Equipment) error {
	_, err := repo.db.Exec("INSERT INTO equipment(name,model,quantity, cost_price, date_received) VALUES(?, ?, ?, ?, ?)", equipment.Name, 
	equipment.ModelOrMake, equipment.Quantity, equipment.CostPrice, equipment.DateReceived)
  return err
}
func (repo * InventoryRepository) GetEquipments()([]model.Equipment, error ){
	equipments := make([]model.Equipment, 0)
	err := repo.db.Select(&equipments, "SELECT id, name, model, quantity, cost_price, date_received FROM equipment where deleted_at is null ORDER BY updated_at DESC	")
  	return equipments, err
}
func (repo * InventoryRepository) UpdateEquipment(equipment model.Equipment) error {
	_, err := repo.db.Exec("UPDATE equipment SET name = ?, model = ?, quantity = ?, cost_price = ?, date_received = ? where id = ?", equipment.Name, 
	equipment.ModelOrMake, equipment.Quantity, equipment.CostPrice, equipment.DateReceived, equipment.Id)
  return err
}
func (repo * InventoryRepository) DeleteEquipment(id int) error {
	_, err := repo.db.Exec("UPDATE equipment SET deleted_at = now() where id = ?", id)
  return err
}

func NewInventoryRepository() InventoryRepository{
	return InventoryRepository{
		db: db.GetConnection(),
	}
}