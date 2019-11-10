package model

type ServiceRegistry struct {
	Key         string
	Name        string
	Description string
	Version     string
	Status      string
}

func UpdateServiceRegistry(key string, s ServiceRegistry) error {
	return SqliteConn.Model(s).Where("key = ?", key).Update(s).Error
}

func GetServiceRegistry(key string) (ServiceRegistry, error) {
	var s ServiceRegistry
	err := SqliteConn.Model(s).Where("key = ?", key).First(&s).Error
	return s, err
}
func InsertServiceRegistry(s ServiceRegistry) error {
	return SqliteConn.Model(s).Create(s).Error
}
