package alipay

func init() {
	DoAsyncAlipay()
}

type PayAction struct {
	Payload map[string]string
	Action  func(map[string]string) (Response, error)
	Return  func(Response, error)
}

var actionChan = make(chan *PayAction)

func PushPayActionToChan(action *PayAction) {
	actionChan <- action
}

func DoAsyncAlipay() {
	go func() {
		for {
			select {
			case action := <-actionChan:
				go func() {
					action.Return(action.Action(action.Payload))
				}()
			}
		}
	}()
}
