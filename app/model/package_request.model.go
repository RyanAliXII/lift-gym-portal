package model

type PackageRequest struct {
	Id        int        `json:"id" db:"id"`
	ClientId  int        `json:"clientId" db:"client_id"`
	PackageId int        `json:"packageId" db:"package_id"`
	StatusId  int        `json:"statusId" db:"status_id"`
	Status    string     `json:"status" db:"status"`
	Client    ClientJSON `json:"client" db:"client"`
	Remarks   string     `json:"remarks" db:"remarks"`
	Model
}