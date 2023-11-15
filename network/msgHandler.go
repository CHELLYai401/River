package network

import "fmt"

type MsgHandler struct {
	RouterMap map[int]IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan Request
	//业务工作Worker池的消息队列
	WorkerPookSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		RouterMap:      make(map[int]IRouter),
		WorkerPookSize: 10,
		TaskQueue:      make([]chan Request, 10),
	}
}

func (mh *MsgHandler) AddRouter(msgId int, router IRouter) {
	if mh.RouterMap[msgId] != nil {
		return
	}
	mh.RouterMap[msgId] = router
}

func (mh *MsgHandler) DoMsgHandle(requset Request) {
	msgId := requset.Msg.GetMsgID()
	if mh.RouterMap[int(msgId)] == nil {
		return
	}
	mh.RouterMap[int(msgId)].PreHandle(requset)
	mh.RouterMap[int(msgId)].Handle(requset)
	mh.RouterMap[int(msgId)].PostHandle(requset)
}

func (mh *MsgHandler) StartWork(workId int, taskQueue chan Request) {
	fmt.Println("启动一条工作线程：", workId)
	for {
		select {
		case requset := <-taskQueue:
			mh.DoMsgHandle(requset)
		}
	}
}

func (mh *MsgHandler) StartWorkPool() {
	for i := 0; i < int(mh.WorkerPookSize); i++ {

		mh.TaskQueue[i] = make(chan Request)
		go mh.StartWork(i, mh.TaskQueue[i])

	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request Request) {

	workerId := request.Conn.GetConnID() % mh.WorkerPookSize
	fmt.Println("Add ConnId = ", request.Conn.GetConnID(),
		" request MsgId = ", request.Msg.GetMsgID(),
		" to WorkerId = ", workerId)
	mh.TaskQueue[workerId] <- request
}
