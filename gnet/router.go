package gnet

import "gamx/giface"

//定义一个基类，规范接口，客户使用时继承基类再自定义接口具体实现
//所有Router都继承BaseRouter
type BaseRouter struct{}

//处理业务前的Hook
func (br *BaseRouter)PreHandle(req giface.IRequest){}

//处理业务的主方法Hook
func (br *BaseRouter)Handler(req giface.IRequest){}

//处理业务之后的Hook
func (br *BaseRouter)PostHandle(req giface.IRequest){}
