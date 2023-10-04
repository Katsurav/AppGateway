package KIALogging

import (
	"gorm.io/gorm"
	"strings"
)

type KIACitation struct {
	gorm.Model
	Id       string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Index    int    `json:"index" gorm:"primaryKey;autoIncrement:false"`
	Response string `json:"response"`
}

func (k *KIACitation) Create(id string, index int, response string) {
	k.Id = id
	k.Index = index
	k.Response = strings.Replace(response, "./documents/", "", -1)
}
