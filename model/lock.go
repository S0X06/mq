package model

type Lock struct {
	Key       string `bson:"key"`
	Value     int    `bson:"value"`
	UpdatedAt int64  `bson:"updated_at"`
}
