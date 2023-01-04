package lib

type Router struct {
	Routes map[string]func(req Request, res Response)
}

func GetRouter() Router {
	return Router{Routes: make(map[string]func(req Request, res Response))}
}

func (r Router) Execute(req Request, res Response) bool {
	route := r.Routes[req.Path]
	if route == nil {
		return false
	}
	route(req, res)
	return true
}
