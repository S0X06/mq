package cron

import (
	. "mq/service"
)

func callBack(ack int32, orderSn string) {

	if ack == REMOVEACK {
		//删除
		MapFn[OP_REMOVE](orderSn)

	} else if ack == PREACK {
		//失败
		MapFn[OP_SORT](orderSn)
	} else if ack == ACK {
		//成功
		MapFn[OP_UPDATE](orderSn)

	} else {
		MapFn[OP_UnLock](orderSn)
	}

	return
}
