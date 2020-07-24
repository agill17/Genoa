package gitSync

import (
	"context"
	"coveros.com/pkg/utils"
	"fmt"
	"github.com/google/go-github/github"
	"reflect"
)

func (wH WebhookHandler) syncHelmReleaseWithGithub(owner, repo, branch, releaseFile string, gitClient *github.Client, isRemovedFromGithub bool) {

	log.Info(fmt.Sprintf("Syncing %v from %v/%v into cluster", releaseFile, owner, repo))
	gitFileContents, errReadingFromGit := utils.GetFileContentsFromGitInString(owner, repo, branch, releaseFile, gitClient)
	if errReadingFromGit != nil {
		log.Error(errReadingFromGit, "Failed to get fileContents from github")
		return
	}

	hrFromGit, unMarshalErr := utils.UnMarshalStringDataToHelmRelease(gitFileContents)
	if unMarshalErr != nil {
		log.Info("Failed to unmarshal")
		return
	}

	if hrFromGit == nil {
		log.Info(fmt.Sprintf("%v is not a valid HelmRelease, therefore skipping", releaseFile))
		return
	}

	if hrFromGit.Spec.ValuesOverride.V == nil {
		hrFromGit.Spec.ValuesOverride.V = map[string]interface{}{}
	}

	if isRemovedFromGithub {
		if err := wH.Client.Delete(context.TODO(), hrFromGit); err != nil {
			log.Error(err, "Failed to delete %v which was removed from github: %v", hrFromGit.GetName(), releaseFile)
			return
		}
		log.Info(fmt.Sprintf("Delete %v HelmRelease from cluster initiated...", hrFromGit.GetName()))
		return
	}

	log.Info(fmt.Sprintf("Creating %v namespace if needed..", hrFromGit.GetNamespace()))
	if errCreatingNamespace := utils.CreateNamespace(hrFromGit.GetNamespace(), wH.Client); errCreatingNamespace != nil {
		log.Error(errCreatingNamespace, "Failed to create namespace")
		return
	}

	log.Info(fmt.Sprintf("Creating %v/%v HelmRelease", hrFromGit.GetNamespace(), hrFromGit.GetName()))
	hrFromCluster, errCreatingHR := utils.CreateHelmRelease(hrFromGit, wH.Client)
	if errCreatingHR != nil {
		log.Info(fmt.Sprintf("%v/%v failed to create helmRelease : %v", hrFromGit.GetNamespace(), hrFromGit.GetName(), errCreatingHR))
	}
	log.Info(fmt.Sprintf("Successfully created %v/%v HelmRelease", hrFromGit.GetNamespace(), hrFromGit.GetName()))

	specInSync := reflect.DeepEqual(hrFromCluster.Spec, hrFromGit.Spec)
	labelsInSync := reflect.DeepEqual(hrFromCluster.GetLabels(), hrFromGit.GetLabels())
	annotationsInSync := reflect.DeepEqual(hrFromCluster.GetAnnotations(), hrFromGit.GetAnnotations())
	if !specInSync || !labelsInSync || !annotationsInSync {
		hrFromCluster.SetAnnotations(hrFromGit.GetAnnotations())
		hrFromCluster.SetLabels(hrFromGit.GetLabels())
		hrFromCluster.Spec = hrFromGit.Spec
		if errUpdating := wH.Client.Update(context.TODO(), hrFromCluster); errUpdating != nil {
			log.Error(errUpdating, fmt.Sprintf("Failed to apply HelmRelease from %v/%v - %v", owner, repo, hrFromGit.GetName()))
			return
		}

		log.Info(fmt.Sprintf("Updated HelmRelease from %v/%v - %v/%v", owner, repo, hrFromGit.GetNamespace(), hrFromGit.GetName()))
	}

}
