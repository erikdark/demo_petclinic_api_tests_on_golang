package models


type PetType struct {
 ID   int32  `json:"id"`
 Name string `json:"name"`
}


type PetFields struct {
 Name      string  `json:"name"`
 BirthDate string  `json:"birthDate"` 
 Type      PetType `json:"type"`
}

// Pet ...
type Pet struct {
 PetFields
 ID      int32   `json:"id"`
 OwnerID int32   `json:"ownerId"`
 Visits  []Visit `json:"visits"`
}


type Visit struct {
 ID          int32  `json:"id"`
 Date        string `json:"date"`
 Description string `json:"description"`
 PetID       int32  `json:"petId"`
}