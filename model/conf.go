package model

type Conf struct {
	AppId        string `bson:"app_id"`
	AppSecret    string `bson:"app_secret"`
	ServceName   string `bson:"servce_name"`
	QueueName    string `bson:"queue_name"`
	RoutingKey   string `bson:"routing_key"`
	ExchangeName string `bson:"exchange_name"`
	ExchangeType string `bson:"exchange_type"`
	Status       int    `bson:"status"` //0不发布 1发布
	UpdatedAt    int64  `bson:"updated_at"`
	CreatedAt    int64  `bson:"created_at"`
}
