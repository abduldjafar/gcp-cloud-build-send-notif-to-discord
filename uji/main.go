package main

import (
	"github.com/DisgoOrg/disgohook/api"

	"github.com/DisgoOrg/disgohook"
	"github.com/sirupsen/logrus"
)

func sendToDiscord(data []*api.EmbedField, status string) {
	logger := logrus.New()
	colours := map[string]int{}
	colours["QUEUED"] = 15258703
	colours["WORKING"] = 1127128
	colours["DONE"] = 6606392
	colours["FAILED"] = 14177041
	colours["CANCELLED"] = 14177041
	colours["TIMEOUT"] = 14177041

	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting example...")

	webhook, err := disgohook.NewWebhookClientByToken(nil, logger, "869350087736819773/vWIqYETISw8pGWFz9TAvLd2w0lLxU4tXlcWFFQzKXju35prz1QOWAV0YGslXLQolrstq")
	if err != nil {
		logger.Errorf("failed to create webhook: %s", err)

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

	}
}

func main() {

	message := []*api.EmbedField{
		&api.EmbedField{
			Name:  "projectId",
			Value: "quantum-bonus-318613",
		},
		&api.EmbedField{
			Name:  "repoName",
			Value: "naruto-xxxx-com",
		},
		&api.EmbedField{
			Name:  "status",
			Value: "DONE",
		},
	}
	sendToDiscord(message, "FAILED")
}
