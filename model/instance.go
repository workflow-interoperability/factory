package model

type Instance struct {
	Key         string `gorm:"primary_key"`
	State       int
	Name        string
	Subject     string
	Description string
	FactoryKey  string
	Observers   string // use `|` to split different observers
	ContextData string
	ResultData  string
	History     string
}

func GetInstance(key string) (Instance, error) {
	var ret Instance
	data := SqliteConn.Model(ret).Where("id = ?", key).First(&ret)
	return ret, data.Error
}

func InsertInstance(i Instance) error {
	return SqliteConn.Model(i).Create(i).Error
}

func ListInstances() ([]Instance, error) {
	var ret []Instance
	return ret, SqliteConn.Model(&Instance{}).Find(&ret).Error
}

func UpdateInstance(i Instance) error {
	return SqliteConn.Model(i).Where("key = ?", i.Key).Update(i).Error
}

func AddInstanceObserver(key string, observers string) error {
	return SqliteConn.Model(&Instance{}).Where("key = ?", key).UpdateColumn("observers = ?", observers).Error
}
