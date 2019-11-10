package model

type Factory struct {
	Key               string `gorm:"primary_key"`
	Name              string
	Subject           string
	Description       string
	ContextDataSchema string
	ResultDataSchema  string
	Expiration        int64
	Language          string
	Definition        string
}

func GetFactory(key string) (Factory, error) {
	var ret Factory
	data := SqliteConn.Model(ret).Where("key = ?", key).First(&ret)
	return ret, data.Error
}

func InsertFactory(f Factory) error {
	return SqliteConn.Model(f).Create(f).Error
}

func ListFactories() ([]Factory, error) {
	var ret []Factory
	return ret, SqliteConn.Model(&Factory{}).Find(&ret).Error
}

func UpdateFactory(key string, f Factory) error {
	return SqliteConn.Model(f).Where("key = ?", key).Update(f).Error
}
