//revive:disable:package-comments,exported
package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		repositoryName := "ctr-latex"
		repository, err := github.NewRepository(ctx, "newRepositoryCtrLatex", &github.RepositoryArgs{
			DeleteBranchOnMerge: pulumi.Bool(true),
			Description:         pulumi.String("This is a repository for containerized LaTeX environments."),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			Name:                pulumi.String(repositoryName),
			Topics: pulumi.StringArray{
				pulumi.String("dagger"),
				pulumi.String("docker"),
				pulumi.String("github"),
				pulumi.String("gitlab"),
				pulumi.String("go"),
				pulumi.String("golang"),
				pulumi.String("latex"),
				pulumi.String("pulumi"),
			},
			Visibility: pulumi.String("public"),
			// VulnerabilityAlerts: pulumi.Bool(true),
		}, pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewBranchProtection(ctx, "branchProtection", &github.BranchProtectionArgs{
			RepositoryId:          repository.NodeId,
			Pattern:               pulumi.String("main"),
			RequiredLinearHistory: pulumi.Bool(true),
		}, pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(ctx, "newIssueLabelGhActions", &github.IssueLabelArgs{
			Color:       pulumi.String("E66E01"),
			Description: pulumi.String("This issue is related to github-actions dependencies"),
			Name:        pulumi.String("github-actions dependencies"),
			Repository:  repository.Name,
		}, pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(ctx, "newIssueLabelGoModules", &github.IssueLabelArgs{
			Color:       pulumi.String("9BE688"),
			Description: pulumi.String("This issue is related to go modules dependencies"),
			Name:        pulumi.String("go-modules dependencies"),
			Repository:  repository.Name,
		}, pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(ctx, "newIssueLabelDockerDep", &github.IssueLabelArgs{
			Color:       pulumi.String("0A4DFE"),
			Description: pulumi.String("This issue is related to docker dependencies"),
			Name:        pulumi.String("docker dependencies"),
			Repository:  repository.Name,
		})
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretGLR", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("GITLAB_REPOSITORY"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretGLT", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("GITLAB_TOKEN"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretGLO", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("GITLAB_OWNER"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretDU", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("DOCKERHUB_USERNAME"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretDR", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("DOCKERHUB_REPOSITORY"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, "newActionsSecretDT", &github.ActionsSecretArgs{
			Repository: repository.Name,
			SecretName: pulumi.String("DOCKERHUB_TOKEN"),
		}, pulumi.Parent(repository), pulumi.Protect(false))
		if err != nil {
			return err
		}

		ctx.Export("repository", repository.Name)
		ctx.Export("repositoryUrl", repository.HtmlUrl)

		return nil
	})
}
