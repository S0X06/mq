package utils

const (
	SUCCESS = 0
	FAILURE = -1
)

var msg map[int]string = map[int]string{
	SUCCESS: "成功",
	FAILURE: "失败",
}

func GetCode(c int) (code int, message string) {
	var ok bool
	code = c
	message, ok = msg[c]

	if !ok {
		message = msg[10001]
	}
	return
}
