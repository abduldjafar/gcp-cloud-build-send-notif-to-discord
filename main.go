// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/sirupsen/logrus"
)

type Data struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Source struct {
		RepoSource struct {
			ProjectID string `json:"projectId"`
			RepoName  string `json:"repoName"`
			CommitSha string `json:"commitSha"`
		} `json:"repoSource"`
	} `json:"source"`
	CreateTime time.Time `json:"createTime"`
	StartTime  time.Time `json:"startTime"`
	Steps      []struct {
		Name string   `json:"name"`
		Args []string `json:"args"`
	} `json:"steps"`
	Timeout          string `json:"timeout"`
	ProjectID        string `json:"projectId"`
	LogsBucket       string `json:"logsBucket"`
	SourceProvenance struct {
		ResolvedRepoSource struct {
			ProjectID string `json:"projectId"`
			RepoName  string `json:"repoName"`
			CommitSha string `json:"commitSha"`
		} `json:"resolvedRepoSource"`
	} `json:"sourceProvenance"`
	BuildTriggerID string `json:"buildTriggerId"`
	Options        struct {
		SubstitutionOption   string `json:"substitutionOption"`
		Logging              string `json:"logging"`
		DynamicSubstitutions bool   `json:"dynamicSubstitutions"`
	} `json:"options"`
	LogURL        string `json:"logUrl"`
	Substitutions struct {
		RevisionID  string `json:"REVISION_ID"`
		CommitSha   string `json:"COMMIT_SHA"`
		ShortSha    string `json:"SHORT_SHA"`
		BranchName  string `json:"BRANCH_NAME"`
		RefName     string `json:"REF_NAME"`
		TriggerName string `json:"TRIGGER_NAME"`
		RepoName    string `json:"REPO_NAME"`
	} `json:"substitutions"`
	Tags     []string `json:"tags"`
	QueueTTL string   `json:"queueTtl"`
	Name     string   `json:"name"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func sendToDiscord(data []*api.EmbedField, status string) error {
	logger := logrus.New()
	colours := map[string]int{}
	colours["QUEUED"] = 15258703
	colours["WORKING"] = 1127128
	colours["DONE"] = 6606392
	colours["FAILED"] = 14177041
	colours["CANCELLED"] = 14177041
	colours["TIMEOUT"] = 14177041

	logger.SetLevel(logrus.DebugLevel)

	webhook, err := disgohook.NewWebhookClientByToken(nil, logger, "webhook token")
	if err != nil {
		logger.Errorf("failed to create webhook: %s", err)
		return err

	}
	var colour int = colours[status]
	var title string = "Build Status"

	_, err = webhook.SendMessage(api.NewWebhookMessageCreateBuilder().
		SetEmbeds(api.Embed{
			Color:  &colour,
			Title:  &title,
			Fields: data,
		}).
		Build(),
	)
	if err != nil {
		logger.Errorf("failed to send webhook message: %s", err)
		return err

	}

	return nil
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {

	str_data := string(m.Data)
	data := Data{}
	json.Unmarshal([]byte(str_data), &data)

	message := []*api.EmbedField{
		&api.EmbedField{
			Name:  "projectId",
			Value: data.Source.RepoSource.ProjectID,
		},
		&api.EmbedField{
			Name:  "repoName",
			Value: data.Source.RepoSource.RepoName,
		},
		&api.EmbedField{
			Name:  "branchName",
			Value: data.Substitutions.BranchName,
		},
		&api.EmbedField{
			Name:  "triggerName",
			Value: data.Substitutions.TriggerName,
		},
		&api.EmbedField{
			Name:  "status",
			Value: data.Status,
		},
		&api.EmbedField{
			Name:  "logs",
			Value: data.LogURL,
		},
	}

	if err := sendToDiscord(message, data.Status); err != nil {
		return err
	}

	return nil
}
