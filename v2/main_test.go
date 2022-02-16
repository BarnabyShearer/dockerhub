package dockerhub

import (
    "testing"
	"fmt"
	"context"
	"os"
)

// DOCKER_USERNAME / PASSOWRD must be valid Magenta ApS credentials

func TestReadRepository(t *testing.T) {
    name := "magentaaps/lora"
	client := NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	repository, err := client.GetRepository(context.Background(), name)
	if err != nil {
        t.Fatalf(`Got error: %v`, err)
	}
    if repository.Name != "lora" {
        t.Fatalf(`Name wrong, got %s, expected %s`, repository.Name, name)
    }
}

func TestReadGroup(t *testing.T) {
    organization_name := "magentaaps"
    group_id := "owners"
	client := NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	group, err := client.GetGroup(context.Background(), organization_name, group_id)
	if err != nil {
        t.Fatalf(`Got error: %v`, err)
	}
    if group.Name != "owners" {
        t.Fatalf(`Name wrong, got %s, expected %s`, group.Name, group_id)
    }
    if group.Id != 51673 {
        t.Fatalf(`Id wrong, got %d, expected %d`, group.Id, 51673)
    }
}

func TestReadGroupFailure(t *testing.T) {
    organization_name := "magentaaps"
    group_id := "unknowngroupnamehere"
	client := NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	_, err := client.GetGroup(context.Background(), organization_name, group_id)
	if err == nil {
        t.Fatalf(`Did not get expected error`)
	}
    expected_err := "{\"detail\": \"Not found\"}"
    string_err := fmt.Sprint(err)
	if string_err != expected_err {
        t.Fatalf(`Wrong error, got %s, expected %s`, string_err, expected_err)
	}
}
