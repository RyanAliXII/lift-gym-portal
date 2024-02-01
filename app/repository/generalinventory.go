package repository

import (
	"context"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)


type GeneralInventory struct {
	db *sqlx.DB
	objstore objstore.ObjectStorer
}
func NewGeneralInventory () GeneralInventory {
	objstore, _ := objstore.GetObjectStorage()

	return GeneralInventory{
		db: db.GetConnection(),
		objstore: objstore,
	}
}
func( repo * GeneralInventory)NewItem(item model.GeneralItem) error {
	dialect := goqu.Dialect("mysql")
	record := goqu.Record{
		"name": item.Name,
		"brand": item.Brand,
		"cost_price": item.CostPrice,
		"quantity": item.Quantity,
		"quantity_threshold": item.QuantityThreshold,
		"image": item.Image,
		"date_received": item.DateReceived,
		"unit_of_measure": item.UnitOfMeasure,
	}
	if(len(item.ExpirationDate) > 0){
		record["expiration_date"] = item.ExpirationDate
	}
	if(item.ImageFile != nil){
		file, err :=  item.ImageFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		filename := uuid.New()
		filePath, err := repo.objstore.Upload(context.Background(), file, objstore.UploadConfig{
			FolderName: "/items",
			Filename: filename.String(),
			AllowedFormats: []string{"jpg", "png", "webp"},
		})
		if err != nil {
			return err
		}
		record["image"] = filePath
	}
	ds := dialect.Insert(goqu.T("general_inventory")).Rows(record).Prepared(true)
	sql, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(sql, args...)
	return err
}
func( repo * GeneralInventory)Update(item model.GeneralItem) error {
	dialect := goqu.Dialect("mysql")
	record := goqu.Record{
		"name": item.Name,
		"brand": item.Brand,
		"cost_price": item.CostPrice,
		"quantity": item.Quantity,
		"quantity_threshold": item.QuantityThreshold,
		"image": item.Image,
		"date_received": item.DateReceived,
		"unit_of_measure": item.UnitOfMeasure,
	}
	if(len(item.ExpirationDate) > 0){
		record["expiration_date"] = item.ExpirationDate
	}
	if(item.ImageFile != nil){
		file, err :=  item.ImageFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		filename := uuid.New()
		filePath, err := repo.objstore.Upload(context.Background(), file, objstore.UploadConfig{
			FolderName: "/items",
			Filename: filename.String(),
			AllowedFormats: []string{"jpg", "png", "webp"},
		})
		if err != nil {
			return err
		}
		record["image"] = filePath
	}
	ds := dialect.Update(goqu.T("general_inventory")).Set(record).Prepared(true).Where(goqu.Ex{"id": item.Id })
	sql, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(sql, args...)
	return err
}
func (repo * GeneralInventory) GetItems() ([]model.GeneralItem, error) {
	items := make([]model.GeneralItem, 0)
	err := repo.db.Select(&items, `SELECT id,name, brand, quantity,
	 unit_of_measure,cost_price,date_received, quantity_threshold, 
	 COALESCE(expiration_date, '') 
	 as expiration_date, image from general_inventory 	 where deleted_at is null ORDER BY created_at desc`)

	return items, err
}
func (repo * GeneralInventory) DeleteItem(id int) error {
	_, err := repo.db.Exec("DELETE FROM general_inventory WHERE id = ? and deleted_at is null", id)
	return err
}