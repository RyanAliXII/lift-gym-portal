package model

import (
	"fmt"
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Equipment struct {
	Id 			 int `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	ModelOrMake  string  `json:"model" db:"model"`
	Quantity     int     `json:"quantity" db:"quantity"`
	CostPrice    float64 `json:"costPrice" db:"cost_price"`
	Condition 		int `json:"condition" db:"condition"`	
	QuantityThreshold int `json:"quantityThreshold" db:"quantity_threshold"`
	ConditionThreshold int `json:"conditionThreshold" db:"condition_threshold"`
	DateReceived string  `json:"dateReceived" db:"date_received"`
	Model
}

type InventoryStat struct {
 	TotalCost float64 `json:"totalCost" db:"total_cost"`
}

func (m Equipment) Validate()(error, map[string]string) {
	return m.Model.ValidationRules(&m, 
	validation.Field(&m.Name, validation.Required.Error("Equipment name is required.",), validation.Length(1, 100).Error("Equipment name should be atleast 1 to 100 characters.")),
    validation.Field(&m.ModelOrMake, validation.Required.Error("Model/make is required.",), validation.Length(1, 100).Error("Model/make should be atleast 1 to 100 characters.")),
	validation.Field(&m.Quantity, validation.Required.Error("Quantity is required."), validation.Min(1).Error("Quantity shoud be atleast 1")),
	validation.Field(&m.CostPrice, validation.Required.Error("Cost price is required."), validation.Min(float64(1)).Error("Cost price should be atleast 1")),
	validation.Field(&m.Condition, validation.Required.Error("Condition is required."), validation.Min(1).Error("Condition must not be less than 1"), validation.Max(100).Error("Condition must no exceed 100")),
	validation.Field(&m.QuantityThreshold, validation.Min(0).Error("Quantity threshold must not be less than 0")),
	validation.Field(&m.ConditionThreshold, validation.Min(0).Error("Condition threshold must not be less than 0")),
	validation.Field(&m.DateReceived, validation.Required.Error("Date received is required."), validation.By(func(value interface{}) error {
		format := "2006-01-02"
		strDate ,_ := value.(string)
		_, err := time.Parse(format, strDate)
		if err != nil {
			return fmt.Errorf("Date of birth is required.")
		}
		return nil
	})))
	
}

type GeneralItem struct {
	Id 			 int `json:"id" form:"id" db:"id"`
	Name         string  `json:"name" form:"name" db:"name"`
	Brand  string  `json:"brand" form:"brand" db:"brand"`
	Quantity     int     `json:"quantity" form:"quantity" db:"quantity"`
	CostPrice    float64 `json:"costPrice" form:"costPrice" db:"cost_price"`	
	QuantityThreshold int `json:"quantityThreshold" form:"quantityThreshold" db:"quantity_threshold"`
	UnitOfMeasure string `json:"unitOfMeasure" form:"unitOfMeasure" db:"unit_of_measure"`
	DateReceived string  `json:"dateReceived" form:"dateReceived" db:"date_received"`
	ExpirationDate string `json:"expirationDate" form:"expirationDate" db:"expiration_date"`
	Image string `json:"image" db:"image"`
	ImageFile *multipart.FileHeader `json:"imageFile"`
	Model
}

func (m * GeneralItem) Validate()(map[string]string, error) {
	return m.Model.Validate(m, 
	validation.Field(&m.Name, validation.Required.Error("Item name is required.",), validation.Length(1, 100).Error("Item name should be atleast 1 to 100 characters.")),
    validation.Field(&m.Brand, validation.Required.Error("Brand is required.",), validation.Length(1, 100).Error("Brand should be atleast 1 to 100 characters.")),
	validation.Field(&m.Quantity, validation.Required.Error("Quantity is required."), validation.Min(1).Error("Quantity shoud be atleast 1")),
	validation.Field(&m.CostPrice, validation.Required.Error("Cost price is required."), validation.Min(float64(1)).Error("Cost price should be atleast 1")),
	validation.Field(&m.UnitOfMeasure, validation.Required.Error("Unit of measure is required."), validation.In("pack", "pieces", "pounds", "kilograms", "grams", "litres")),
	validation.Field(&m.QuantityThreshold,validation.Min(0).Error("Quantity threshold must not be less than 0")),
	validation.Field(&m.ExpirationDate,  validation.By(func(value interface{}) error {
		if(len(m.ExpirationDate) == 0){
			return nil
		}
		format := "2006-01-02"
		strDate ,_ := value.(string)
		_, err := time.Parse(format, strDate)
		if err != nil {
			return fmt.Errorf("Invalid expiration date")
		}
		return nil
	})),
	validation.Field(&m.ImageFile, validation.By(func(value interface{}) error {
		if(m.ImageFile == nil) {
			return nil
		}
		contentType := m.ImageFile.Header["Content-Type"][0]
		if(contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp"){
			return fmt.Errorf("image should be .png, .jpg or .webp")
		}
		return nil
	})),
	validation.Field(&m.DateReceived, validation.Required.Error("Date received is required."), validation.By(func(value interface{}) error {
		format := "2006-01-02"
		strDate ,_ := value.(string)
		_, err := time.Parse(format, strDate)
		if err != nil {
			return fmt.Errorf("Date of birth is required.")
		}
		return nil
	})))

}