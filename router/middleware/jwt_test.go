package middleware

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GenerateToken(t *testing.T) {
	Convey("generate token : ", t, func() {
		var data map[string]interface{} = map[string]interface{}{
			"id":        1,
			"user_name": "s06",
		}
		token, err := GenerateToken(data)
		fmt.Println(token)
		So(err, ShouldBeNil)
	})
}
