// Copyright (c) 2016 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

package autocreation

import (
	"strconv"

	"github.com/mattermost/mattermost-load-test/loadtestconfig"
	"github.com/mattermost/platform/model"
)

type ChannelsCreationResult struct {
	Channels []*model.Channel
	Errors   []error
}

func CreateChannels(client *model.Client, config *loadtestconfig.ChannelsConfiguration) *ChannelsCreationResult {
	totalChannels := config.NumChannelsPerTeam * len(config.TeamIds)

	channelResults := &ChannelsCreationResult{
		Channels: make([]*model.Channel, totalChannels),
		Errors:   make([]error, totalChannels),
	}

	for _, teamId := range config.TeamIds {
		client.SetTeamId(teamId)

		ThreadSplit(config.NumChannelsPerTeam, config.CreateThreads, func(channelNumber int) {
			channel := &model.Channel{
				Name:        config.ChannelNamePrefix + strconv.Itoa(channelNumber),
				DisplayName: config.ChannelDisplayName + strconv.Itoa(channelNumber),
				Type:        model.CHANNEL_OPEN,
				TeamId:      teamId,
			}

			if config.UseRandomId {
				channel.Name = channel.Name + model.NewId()
			}

			result, err := client.CreateChannel(channel)
			if err != nil {
				channelResults.Errors[channelNumber] = err
			} else {
				channelResults.Channels[channelNumber] = result.Data.(*model.Channel)
			}
		})
	}

	client.SetTeamId("")

	return channelResults
}
