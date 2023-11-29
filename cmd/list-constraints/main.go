package main

import (
	"context"
	"flag"
	"fmt"

	"cloud.google.com/go/orgpolicy/apiv2"
	"cloud.google.com/go/orgpolicy/apiv2/orgpolicypb"
	"google.golang.org/api/iterator"
)

func main() {
	projectId := flag.String("project-id", "", "The GCP project ID")

	// Parse the command-line arguments
	flag.Parse()

	// Check if the "project-id" flag is provided
	if *projectId == "" {
		fmt.Println("Usage: list-constraints --project-id=your_project_id")
		return
	}

	PrintConstraints(*projectId)
}

// The Organization Policy API must be enabled for the project in order for this method to work!
func PrintConstraints(projectId string) {
	client, err := orgpolicy.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	constraints := client.ListConstraints(context.TODO(), &orgpolicypb.ListConstraintsRequest{
		Parent: fmt.Sprintf("projects/%s", projectId),
	})
	for {
		constraint, err := constraints.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			panic(err)
		}
		fmt.Println(constraint)
	}
}
