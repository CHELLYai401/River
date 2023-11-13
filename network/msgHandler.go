package network

type MsgHandler struct {
	RouterMap map[int]IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		RouterMap: make(map[int]IRouter),
	}
}

func (mh *MsgHandler) AddRouter(msgId int, router IRouter) {
	if mh.RouterMap[msgId] != nil {
		return
	}
	mh.RouterMap[msgId] = router
}

func (mh *MsgHandler) DoMsgHandle(msgId int, requset Request) {
	if mh.RouterMap[msgId] == nil {
		return
	}
	mh.RouterMap[msgId].PreHandle(requset)
	mh.RouterMap[msgId].Handle(requset)
	mh.RouterMap[msgId].PostHandle(requset)
}
