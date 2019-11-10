package model

type Activity struct {
	Key            string `gorm:"primary_key"`
	State          int
	Name           string
	Description    string
	ValidStates    string // split by `|`
	InstanceKey    string
	RemoteInstance string
	StartedDate    int64
	DueDate        int64
	LastModified   int64
}

func GetActivity(key string) (Activity, error) {
	var activity Activity
	if err := SqliteConn.Model(activity).Where("key = ?", key).First(&activity).Error; err != nil {
		return Activity{}, err
	}
	return activity, nil
}

func UpdateActivity(key string, a Activity) error {
	return SqliteConn.Model(a).Where("key = ?", key).Update(a).Error
}
