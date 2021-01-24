package giface

type IMsgHandler interface {
	DoMsgHandler(req IRequest) error
	AddRouter(msgId uint32, router IRouter) error
	StartWorkPool()
	SendMsgToTaskQueue(req IRequest)
}
