package main

import (
	"crave/miner/cmd/configuration"
	"crave/miner/cmd/lib"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	//router.Use()
	var wg sync.WaitGroup
	wg.Add(1)
	//router.Use()
	router := gin.Default()
	go func() {
		defer wg.Done()
		startLib(router)
	}()
	wg.Wait()

}

func startLib(router *gin.Engine) {
	container := configuration.NewContainer(router)
	go lib.Start(container)

}
