package util

func SafeSend(ch chan interface{}, value []byte) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}

func SafeSendString(ch chan interface{},value string) (closed bool) {
	defer func() {
		if recover() != nil {
			// the return result can be altered
			// in a defer function call
			closed = true
		}
	}()
	ch <-value
	return false
}
