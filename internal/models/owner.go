package models


type OwnerFields struct {
 FirstName string `json:"firstName"`
 LastName  string `json:"lastName"`
 Address   string `json:"address"`
 City      string `json:"city"`
 Telephone string `json:"telephone"`
}


type Owner struct {
 OwnerFields
 ID   int32 `json:"id"`
 Pets []Pet `json:"pets"`
}