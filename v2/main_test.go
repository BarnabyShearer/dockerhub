package dockerhub

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// DOCKER_USERNAME / PASSOWRD must be valid credentials
// DOCKER_REPOSITORY and DOCKER_GROUP_ID must be valid

func TestReadRepository(t *testing.T) {
	name := os.Getenv("DOCKER_TEST_REPO")
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
	organization_name := strings.Split(os.Getenv("DOCKER_REPOSITORY"), "/")[0]
	group_id := "owners"
	client := NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	group, err := client.GetGroup(context.Background(), organization_name, group_id)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	if group.Name != "owners" {
		t.Fatalf(`Name wrong, got %s, expected %s`, group.Name, group_id)
	}
	if strconv.Itoa(group.Id) != os.Getenv("DOCKER_GROUP_ID") {
		t.Fatalf(`Id wrong, got %d, expected %d`, group.Id, 51673)
	}
}

func TestReadGroupFailure(t *testing.T) {
	organization_name := strings.Split(os.Getenv("DOCKER_TEST_REPO"), "/")[0]
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

func TestRepositoryPrivacy(t *testing.T) {
	organization_name := strings.Split(os.Getenv("DOCKER_TEST_REPO"), "/")[0]
	repository_name := strings.Split(os.Getenv("DOCKER_TEST_REPO"), "/")[1]
	client := NewClient(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	repo := Repository{
		Private:   true,
		Name:      repository_name,
		Namespace: organization_name,
	}
	created, err := client.CreateRepository(context.Background(), repo)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	if created.Private != true {
		t.Fatalf(`Repository not private, got %t, expected %t`, created.Private, true)
	}
	repo.Private = false
	err = client.UpdateRepository(context.Background(), fmt.Sprintf("%s/%s", created.Namespace, created.Name), repo)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	updated, err := client.GetRepository(context.Background(), fmt.Sprintf("%s/%s", created.Namespace, created.Name))
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	if updated.Private != false {
		t.Fatalf(`Repository not private, got %t, expected %t`, updated.Private, false)
	}
}
