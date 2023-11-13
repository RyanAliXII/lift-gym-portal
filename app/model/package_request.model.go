package model

import validation "github.com/go-ozzo/ozzo-validation"

type PackageRequest struct {
	Id        int         `json:"id" db:"id"`
	ClientId  int         `json:"clientId" db:"client_id"`
	PackageId int         `json:"packageId" db:"package_id"`
	StatusId  int         `json:"statusId" db:"status_id"`
	Status    string      `json:"status" db:"status"`
	Client    ClientJSON  `json:"client" db:"client"`
	Remarks   string      `json:"remarks" db:"remarks"`
	Package   PackageJSON `json:"package" db:"package"`
	PackageSnapshot   PackageJSON `json:"packageSnapshot" db:"package_snapshot"`
	CreatedAt string 	  `json:"createdAt" db:"created_at"`
	Model
}

func(m PackageRequest)Validate() (map[string]string, error){
	return m.Model.Validate(&m, 
	validation.Field(&m.ClientId, validation.Required.Error("Client is required.")),
	validation.Field(&m.PackageId, validation.Required.Error("Package is required.")))
} 