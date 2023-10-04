package KIALogging

import (
	"gorm.io/gorm"
	"time"
)

type KIAQuery struct {
	gorm.Model
	Id        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	User      string    `json:"user"`
	QueryTime time.Time `json:"querytime" gorm:"default:current_timestamp(3)"`
	Query     string    `json:"query"`
	Pre       string    `json:"pre"`
	Post      string    `json:"post"`
}

func (k *KIAQuery) Create(id string, user string, query string, pre string, post string) {
	k.Id = id
	k.User = user
	k.Query = query
	k.Pre = pre
	k.Post = post
}
