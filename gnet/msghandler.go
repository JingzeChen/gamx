package gnet

import (
	"errors"
	"fmt"
	"gamx/giface"
	"gamx/utils"
)

type MsgHandle struct{
	//Map from requst message id to Router.
	Apis map[uint32]giface.IRouter
	//Number of workers in workpool.
	WorkPoolSize uint32
	//workpool Queue to handle request.
	TaskQueue []chan giface.IRequest
}

func (mh *MsgHandle)AddRouter(msgId uint32, router giface.IRouter) error{
	if _, ok := mh.Apis[msgId]; ok{
		fmt.Println("No need to Register Router because it already exits.")
		return nil
	}

	mh.Apis[msgId] = router
	fmt.Println("Succeed to Add Router with message id: ", msgId)
	return nil
}

func (mh *MsgHandle)DoMsgHandler(req giface.IRequest) error{
	handler, ok := mh.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("failed to get Router with MsgId: ", req.GetMsgId())
		return errors.New("failed to get Router with MsgId")
	}
	fmt.Println("is doing preHandler req id:", req.GetMsgId(), " data: ", req.GetData())
	//handler.PreHandle(req)
	handler.Handle(req)
	//handler.PostHandle(req)
	return nil
}

func (mh *MsgHandle)StartWorkPool(){
	fmt.Println("[Start Work Pool to handle request.]")

	for i:=0; i<int(mh.WorkPoolSize); i++{
		mh.TaskQueue[i] = make(chan giface.IRequest, utils.GlobalObject.MaxTaskQueueLen)
		go mh.StartOneWorker(uint32(i), mh.TaskQueue[i])
	}

}

func (mh *MsgHandle)StartOneWorker(workId uint32, taskQueue chan giface.IRequest){
	fmt.Println("Start WorkId: ", workId)
	for {
		select {
			case req := <-taskQueue:
				mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandle)SendMsgToTaskQueue(req giface.IRequest){
	workId := req.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println("Send request ConnID ", req.GetConnection().GetConnID(), " to worker id ", workId)

	mh.TaskQueue[workId] <- req
}

func NewMsgHandler() *MsgHandle {
	return &MsgHandle{
		Apis:make(map[uint32]giface.IRouter),
		WorkPoolSize:utils.GlobalObject.WorkPoolSize,
		TaskQueue: make([]chan giface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}
