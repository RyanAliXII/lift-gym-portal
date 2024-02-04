package repository

import (
	"context"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
	objStore objstore.ObjectStorer
}
func (repo * InventoryRepository) NewEquipment(equipment model.Equipment) error {
	if(equipment.ImageFile != nil){
		file, err := equipment.ImageFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		filename := uuid.New().String()
		
		objname, err := repo.objStore.Upload(context.Background(), file, objstore.UploadConfig{
			Filename: filename,
			FolderName: "/items",
			AllowedFormats: []string{"png","jpg", "webp"},
		} )
		if err != nil {
			return err
		}
		equipment.Image = objname
	}

      _, err := repo.db.Exec("INSERT INTO equipment(name,model,quantity, cost_price, date_received, `condition`, quantity_threshold, condition_threshold, image) VALUES(?, ?, ?, ?, ?, ?, ?, ?,?)", equipment.Name, 
	  equipment.ModelOrMake, equipment.Quantity, equipment.CostPrice, equipment.DateReceived, equipment.Condition, equipment.QuantityThreshold, equipment.ConditionThreshold, equipment.Image)
	return err
}
func (repo * InventoryRepository) GetStat( ) (model.InventoryStat, error) {
	stat := model.InventoryStat{}
    err := repo.db.Get(&stat, "SELECT COALESCE(SUM(cost_price), 0) as total_cost from equipment where deleted_at is null")
  	return stat, err
}
func (repo * InventoryRepository) GetEquipments()([]model.Equipment, error ){
	equipments := make([]model.Equipment, 0)
	err := repo.db.Select(&equipments, "SELECT id, name, model, quantity, cost_price, date_received, `condition`, quantity_threshold, condition_threshold, COALESCE(image, '') as image FROM equipment where deleted_at is null ORDER BY updated_at DESC	")
  	return equipments, err
}
func (repo * InventoryRepository) UpdateEquipment(equipment model.Equipment) error {
	if(equipment.ImageFile != nil){
		file, err := equipment.ImageFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		filename := uuid.New().String()
		
		objname, err := repo.objStore.Upload(context.Background(), file, objstore.UploadConfig{
			Filename: filename,
			FolderName: "/items",
			AllowedFormats: []string{"png","jpg", "webp"},
		} )
		if err != nil {
			return err
		}
		equipment.Image = objname
	}

	_, err := repo.db.Exec("UPDATE equipment SET name = ?, model = ?, quantity = ?, cost_price = ?, date_received = ?, `condition` = ?, quantity_threshold = ?, condition_threshold = ?, image = ? where id = ?", equipment.Name, 
	equipment.ModelOrMake, equipment.Quantity, equipment.CostPrice, equipment.DateReceived,equipment.Condition, equipment.QuantityThreshold, equipment.ConditionThreshold,equipment.Image, equipment.Id)
  return err
}
func (repo * InventoryRepository) DeleteEquipment(id int) error {
	_, err := repo.db.Exec("UPDATE equipment SET deleted_at = now() where id = ?", id)
  return err
}

func NewInventoryRepository() InventoryRepository{
	objStore, _ := objstore.GetObjectStorage()
	return InventoryRepository{
		db: db.GetConnection(),
		objStore: objStore,
	}
}