package main

import (
	l "main.go/lib"
)

func helloWorld(req l.Request, res l.Response) {
	res.Redirect("/test")
}

func main() {
	s := l.CreateServer(5000)

	s.AddRoute("/", helloWorld)

	s.Run()
}
