package model

type Filter struct {
	OrderSn   string `bson:"order_sn"`
	Ack       int    `bson:"ack"`
	CreatedAt int64  `bson:"created_at"`
}

//解析单号
func ParseFilter(filters *[]Filter) []string {

	var filterItems []string
	for _, filter := range *filters {
		if filter.OrderSn != "" {
			filterItems = append(filterItems, filter.OrderSn)
		}
	}

	return filterItems
}
