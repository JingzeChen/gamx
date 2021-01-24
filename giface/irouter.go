package giface

/*
	define Router module to handle client operation.
*/

type IRouter interface{
	//处理业务之前的操作
	PreHandle(req IRequest)
	//处理主业务
	Handle(req IRequest)
	//主业务处理完成后的操作
	PostHandle(req IRequest)
}
