package main

import (
	"fmt"
	"github.com/remikj/pages-ms/src/contoller"
	"github.com/remikj/pages-ms/src/repository"
	"github.com/remikj/pages-ms/src/server"
	"github.com/remikj/pages-ms/src/service"
)

func main() {
	fmt.Println("Starting application")
	pageRepositoryAsync, err := repository.InitPageRepositoryAsyncFromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	serverImpl, err := server.NewServerFromEnv(
		contoller.NewPageController(
			service.NewPageService(pageRepositoryAsync),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = serverImpl.Run()
	if err != nil {
		fmt.Println(err)
	}
}
