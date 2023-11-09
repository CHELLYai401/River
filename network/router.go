package network

type IRouter interface {
	//在处理conn业务之前的钩子方法Hook
	PreHandle(requset *Request)
	//在处理conn业务时的主方法
	Handle(request *Request)
	//在处理conn业务之后的钩子方法Hook
	PostHandle(request *Request)
}

// 根据需要进行重写
type BaseRouter struct {
}

// 之所以BaseRouter的方法为空，是因为Router可以选择实现里面的方法，不需要全部实现
// 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(requset *Request) {}

// 在处理conn业务时的主方法
func (br *BaseRouter) Handle(request *Request) {}

// 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request *Request) {}
