// +build integration

package tg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(fmt.Sprintf("env '%s' is not provided, but required for tests", k))
	}
	return v
}

func getEnvInt(k string) int {
	v := getEnv(k)

	result, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return result
}

var (
	config = struct {
		// Bot Token
		Token string
		// Bot Username
		Username Username

		// Bot ID
		ID UserID

		// Any file_id available for bot
		FileID FileID

		// Any user id who has conversation with bot
		ExampleUserID UserID

		// Any channel where bot is admin.
		ExampleChannelID ChatID
	}{
		Token:            getEnv("TEST_BOT_TOKEN"),
		Username:         Username(getEnv("TEST_BOT_USERNAME")),
		ID:               UserID(getEnvInt("TEST_BOT_ID")),
		FileID:           FileID(getEnv("TEST_BOT_FILE_ID")),
		ExampleUserID:    UserID(getEnvInt("TEST_BOT_EXAMPLE_USER_ID")),
		ExampleChannelID: ChatID(getEnvInt("TEST_EXAMPLE_CHANNEL_ID")),
	}

	integrationClient = NewClient(config.Token)
)

func TestClient_GetMe_Integration(t *testing.T) {
	bot, err := integrationClient.GetMe(context.Background())
	require.NoError(t, err)
	assert.Equal(t, config.ID, bot.ID)
	assert.Equal(t, config.Username, bot.Username)
}

func TestClient_Integration_GetFile(t *testing.T) {
	file, err := integrationClient.GetFile(
		context.Background(),
		config.FileID,
	)

	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotZero(t, file.Size)
	require.NotZero(t, file.Path)

	reader, err := file.NewReader(context.Background())
	require.NoError(t, err)
	require.NotNil(t, reader)

	body, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, file.Size, len(body))
}

func TestClient_Integration_GetUserProfilePhotos(t *testing.T) {
	profilePhotos, err := integrationClient.GetUserProfilePhotos(
		context.Background(),
		config.ExampleUserID,
		nil,
	)

	require.NoError(t, err)
	require.NotNil(t, profilePhotos)

	if profilePhotos.Total > 0 {
		assert.NotEmpty(t, profilePhotos.Items)
	} else {
		t.Log("look like user does not have photos...")
	}
}

func TestClient_GetChat_Integration(t *testing.T) {
	t.Run("Channel", func(t *testing.T) {
		chat, err := integrationClient.GetChat(
			context.Background(),
			config.ExampleChannelID,
		)

		require.NoError(t, err)
		require.NotNil(t, chat)

		assert.NotEmpty(t, chat.Title)
		assert.Equal(t, ChannelChat, chat.Type)
	})

	t.Run("Private", func(t *testing.T) {
		chat, err := integrationClient.GetChat(
			context.Background(),
			config.ExampleUserID,
		)

		require.NoError(t, err)
		require.NotNil(t, chat)

		assert.NotEmpty(t, chat.FirstName)
		assert.Equal(t, PrivateChat, chat.Type)
	})

	t.Run("Error", func(t *testing.T) {
		chat, err := integrationClient.GetChat(
			context.Background(),
			UserID(1),
		)

		require.Error(t, err)
		require.Nil(t, chat)
	})
}

func TestClient_SetChatTitle_Integration(t *testing.T) {
	err := integrationClient.SetChatTitle(
		context.Background(),
		config.ExampleChannelID,
		fmt.Sprintf("mr-linch/go-tg integration tests [%d]", time.Now().Unix()),
	)

	assert.NoError(t, err)
}

func TestClient_SetChatDescription_Integration(t *testing.T) {
	err := integrationClient.SetChatDescription(
		context.Background(),
		config.ExampleChannelID,
		fmt.Sprintf("this channel is used for integration tests of github.com/mr-linch/go-tg\n\n last run: [%d]", time.Now().Unix()),
	)

	assert.NoError(t, err)
}

func TestClient_GetChatMembersCount_Integration(t *testing.T) {
	count, err := integrationClient.GetChatMembersCount(
		context.Background(),
		config.ExampleChannelID,
	)

	assert.NoError(t, err)
	assert.NotZero(t, count)
}

func TestClient_GetChatAdministrators_Integration(t *testing.T) {
	admins, err := integrationClient.GetChatAdministrators(
		context.Background(),
		config.ExampleChannelID,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, admins)

}

func TestClient_Send_TextMessage_Integration(t *testing.T) {
	msg := NewTextMessage(config.ExampleChannelID, "*Text*: `TestClient_Send_TextMessage_Integration`").
		WithParseMode(Markdown).
		WithNotification(false).
		WithWebPagePreview(false).
		WithReplyMarkup(
			NewInlineKeyboardMarkup(
				NewInlineKeyboardRow(
					NewInlineKeyboardButtonURL(
						"Sources",
						"github.com/mr-linch/go-tg",
					),
				),
			),
		)

	sended := Message{}
	ctx := context.Background()

	err := integrationClient.Send(ctx, msg, &sended)
	require.NoError(t, err)
	require.NotZero(t, sended.ID)

	t.Run("ForwardMessage", func(t *testing.T) {
		err := integrationClient.Send(ctx, NewForwardMessage(
			config.ExampleUserID,
			sended,
		), nil)

		assert.NoError(t, err)
	})
}

func TestClient_Send_PhotoMessage_Integration(t *testing.T) {

	photo, err := NewInputFileLocal("testdata/go-work.png")
	require.NoError(t, err, "no test data!")
	defer photo.Close()

	msg := NewPhotoMessage(config.ExampleChannelID, photo).
		WithCaption("*Text*: `TestClient_Send_PhotoMessage_Integration`").
		WithParseMode(Markdown).
		WithNotification(false).
		WithReplyMarkup(
			NewInlineKeyboardMarkup(
				NewInlineKeyboardRow(
					NewInlineKeyboardButtonURL(
						"Sources",
						"github.com/mr-linch/go-tg",
					),
				),
			),
		)

	err = integrationClient.Send(
		context.Background(),
		msg,
		nil,
	)

	assert.NoError(t, err)
}

func TestClient_Send_AudioMessage_Integration(t *testing.T) {
	// open audio file
	audioFile, err := NewInputFileLocal("testdata/audio.mp3")
	require.NoError(t, err, "no test data!")
	defer audioFile.Close()

	// open thumb file
	thumbFile, err := NewInputFileLocal("testdata/gopher.jpg")
	require.NoError(t, err, "no test data!")
	defer thumbFile.Close()

	msg := NewAudioMessage(config.ExampleChannelID, audioFile).
		WithCaption("*Text*: `TestClient_Send_PhotoMessage_Integration`").
		WithTitle("go-tg").
		WithPerformer("mr-linch").
		WithThumb(thumbFile).
		WithDuration(time.Second * 80).
		WithParseMode(Markdown).
		WithNotification(false).
		WithReplyMarkup(
			NewInlineKeyboardMarkup(
				NewInlineKeyboardRow(
					NewInlineKeyboardButtonURL(
						"Sources",
						"github.com/mr-linch/go-tg",
					),
				),
			),
		)

	err = integrationClient.Send(
		context.Background(),
		msg,
		nil,
	)

	assert.NoError(t, err)
}

func TestClient_Webhook_Integration(t *testing.T) {
	const (
		url            = "https://httpbin.org/status/200"
		maxConnections = 1
	)

	allowedUpdates := []UpdateType{UpdateEditedMessage}

	ctx := context.Background()

	t.Run("Install", func(t *testing.T) {
		err := integrationClient.SetWebhook(ctx, url, &WebhookOptions{
			MaxConnections: maxConnections,
			AllowedUpdates: allowedUpdates,
		})

		assert.NoError(t, err)
	})

	t.Run("Retrieve", func(t *testing.T) {
		info, err := integrationClient.GetWebhookInfo(ctx)

		if assert.NoError(t, err) && assert.NotNil(t, info) {
			assert.Equal(t, url, info.URL)
			assert.Equal(t, maxConnections, info.MaxConnections)
			assert.Equal(t, allowedUpdates, info.AllowedUpdates)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := integrationClient.DeleteWebhook(ctx)
		assert.NoError(t, err)
	})

	t.Run("RetrieveAfterDelete", func(t *testing.T) {
		info, err := integrationClient.GetWebhookInfo(ctx)

		if assert.NoError(t, err) && assert.NotNil(t, info) {
			assert.Empty(t, info.URL)
		}
	})

	info, err := integrationClient.GetWebhookInfo(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, info)
}
