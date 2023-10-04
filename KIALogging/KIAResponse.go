package KIALogging

import (
	"gorm.io/gorm"
)

type KIAResponse struct {
	gorm.Model
	Id       string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Rank     int    `json:"rank"`
	Response string `json:"response"`
	AIModel  string `json:"aimodel"`
	Compute  int64  `json:"compute"`
}

func (k *KIAResponse) Create(id string, rank int, response string, aimodel string, compute int64) {
	k.Id = id
	k.Rank = rank
	k.Response = response
	k.AIModel = aimodel
	k.Compute = compute
}
