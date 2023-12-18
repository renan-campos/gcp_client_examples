package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"cloud.google.com/go/orgpolicy/apiv2"
	"cloud.google.com/go/orgpolicy/apiv2/orgpolicypb"
	"google.golang.org/api/iterator"
)

const (
	command_name = "find-constraint"
	needed_flags = "--project-id=your_project_id, --constraint-id=constraint_id_to_find"
)

func main() {
	projectId := flag.String("project-id", "", "The GCP project ID")
	constraintId := flag.String("constraint-id", "", "The GCP constraint to find")

	// Parse the command-line arguments
	flag.Parse()

	// Check if the "project-id" flag is provided
	if *projectId == "" || *constraintId == "" {
		fmt.Printf("Usage: %s %s\n", command_name, needed_flags)
		return
	}

	hasConstraint, err := HasConstraint(context.Background(), *projectId, *constraintId)
	if err != nil {
		// Todo: Error handling
		panic(err)
	}
	if !hasConstraint {
		fmt.Printf("%s not found\n", *constraintId)
		os.Exit(1)
	}

	fmt.Printf("%s found\n", *constraintId)
}

func HasConstraint(ctx context.Context, projectId, constraintId string) (bool, error) {
	client, err := orgpolicy.NewClient(ctx)
	if err != nil {
		return false, err
	}
	constraints := client.ListConstraints(ctx, &orgpolicypb.ListConstraintsRequest{
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
		if constraint.GetName() == fmt.Sprintf("projects/%s/constraints/%s", projectId, constraintId) {
			return true, nil
		}
	}
	return false, nil
}
