package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"cloud.google.com/go/orgpolicy/apiv2"
	"cloud.google.com/go/orgpolicy/apiv2/orgpolicypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	command_name = "has-constraint"
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

	doesHave, err := hasBooleanConstraint(context.Background(), *projectId, *constraintId)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	if !doesHave {
		fmt.Printf("%s does not have %s policy set\n", *projectId, *constraintId)
		os.Exit(1)
	}
	fmt.Printf("%s has %s policy set\n", *projectId, *constraintId)
}

func hasBooleanConstraint(ctx context.Context, projectId, constraintId string) (bool, error) {
	client, err := orgpolicy.NewClient(ctx)
	if err != nil {
		return false, err
	}
	s, err := client.GetEffectivePolicy(ctx, &orgpolicypb.GetEffectivePolicyRequest{
		Name: fmt.Sprintf("projects/%s/policies/%s", projectId, constraintId),
	})
	if err != nil {
		if isNotFoundError(err) {
			return false, nil
		}
		if isPermissionDenied(err) {
			return false, fmt.Errorf("The Organization Policy API must be enabled for the project %s", projectId)
		}
		return false, err
	}
	for _, policyRule := range s.Spec.GetRules() {
		if policyRule.GetEnforce() {
			return true, nil
		}
	}
	return false, nil
}

func isNotFoundError(err error) bool {
	return status.Code(err) == codes.NotFound
}

func isPermissionDenied(err error) bool {
	return status.Code(err) == codes.PermissionDenied
}
