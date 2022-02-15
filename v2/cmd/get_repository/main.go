// Example of reading a repository
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/magentaaps/dockerhub/v2"
)

var nameFlag = flag.String("name", "", "Name of the repository.")

func main() {
	flag.Parse()
	if *nameFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	client := dockerhub.NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	respository, err := client.GetRepository(context.Background(), *nameFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", respository)
}
