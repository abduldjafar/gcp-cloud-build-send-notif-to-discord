// Package p contains a Pub/Sub Cloud Function.
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {

	str_data := string(m.Data)
	data := Data{}
	json.Unmarshal([]byte(str_data), &data)

	message := `
	projectId : ` + data.Source.RepoSource.ProjectID + `,
	repoName : ` + data.Source.RepoSource.RepoName + `,
	status : ` + data.Status + `
	`
	bot, err := tgbotapi.NewBotAPI("telegram api")
	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewMessage(int64(985052364), message)

	bot.Send(msg)
	return nil
}
