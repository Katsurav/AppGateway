package KIALogging

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type KIALogger struct {
	DB *gorm.DB
}

func (l *KIALogger) Connect(unixPath string) bool {
	var err error

	if len(unixPath) > 0 {
		dsnUnix := "kia:Trustno1!@unix(" + unixPath + ")/kia_demo?charset=utf8mb4&parseTime=True&loc=Local"
		l.DB, err = gorm.Open(mysql.Open(dsnUnix), &gorm.Config{})

	} else {
		dsnTCP := "kia:Trustno1!@tcp(35.184.207.242:3306)/kia_demo?charset=utf8mb4&parseTime=True&loc=Local"

		l.DB, err = gorm.Open(mysql.Open(dsnTCP), &gorm.Config{})
	}
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return false
	}

	// Migrate the schema
	l.DB.AutoMigrate(&KIAQuery{})
	l.DB.AutoMigrate(&KIAResponse{})
	l.DB.AutoMigrate(&KIACitation{})

	return true
	/*

		// Create
		db.Create(&Product{Code: "D42", Price: 100})

		// Read
		var product Product
		db.First(&product, 1)                 // find product with integer primary key
		db.First(&product, "code = ?", "D42") // find product with code D42

		// Update - update product's price to 200
		db.Model(&product).Update("Price", 200)
		// Update - update multiple fields
		db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
		db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

		// Delete - delete product
		db.Delete(&product, 1)

	*/
}

func (l *KIALogger) QueryCreate(id string, user string, query string, pre string, post string) {
	// Create
	k := KIAQuery{}
	k.Create(id, user, query, pre, post)
	l.DB.Create(&k)
}

func (l *KIALogger) ResponseCreate(id string, rank int, response string, aimodel string, compute int64) {
	// Create
	k := KIAResponse{}
	k.Create(id, rank, response, aimodel, compute)
	l.DB.Create(&k)
}

func (l *KIALogger) CitationCreate(id string, index int, response string) {
	// Create
	k := KIACitation{}
	k.Create(id, index, response)
	l.DB.Create(&k)
}
