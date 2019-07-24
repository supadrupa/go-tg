package tg

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	token                 string
	transport             Transport
	defaultParseMode      ParseMode
	defaultWebPagePreview bool
}

// ClientOption represents client option.
type ClientOption func(client *Client)

// WithTransport sets client transport.
func WithTransport(t Transport) ClientOption {
	return func(c *Client) {
		c.transport = t
	}
}

// WithParseMode sets client default parse mode.
func WithParseMode(pm ParseMode) ClientOption {
	return func(c *Client) {
		c.defaultParseMode = pm
	}
}

// WithWebPagePreview sets client default web page preview options
func WithWebPagePreview(enable bool) ClientOption {
	return func(c *Client) {
		c.defaultWebPagePreview = enable
	}
}

// NewClient creates a Telegram Bot API client.
func NewClient(token string, options ...ClientOption) *Client {
	client := &Client{
		token:     token,
		transport: NewHTTPTransport(),

		defaultParseMode:      Plain,
		defaultWebPagePreview: true,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// Invoke request.
func (client *Client) Invoke(
	ctx context.Context,
	req *Request,
	dst interface{},
) error {
	req = req.WithToken(client.token)

	res, err := client.transport.Execute(ctx, req)
	if err != nil {
		return err
	}

	// TODO: handle bot api errors here
	if !res.OK {
		return errors.New(res.Description)
	}

	if dst != nil {
		return res.UnmarshalResult(&dst)
	}

	return nil
}

// GetMe returns bot profile.
//
// Source: https://core.telegram.org/bots/api#getme
func (client *Client) GetMe(
	ctx context.Context,
) (user *User, err error) {
	err = client.Invoke(ctx,
		NewRequest("getMe"),
		&user,
	)

	return
}

// GetFile returns file info and path to download.
//
// Source: https://core.telegram.org/bots/api#getfile
func (client *Client) GetFile(
	ctx context.Context,
	id FileID,
) (file *File, err error) {
	req := NewRequest("getFile").
		AddString("file_id", string(id))

	err = client.Invoke(ctx,
		req,
		&file,
	)

	if file != nil {
		file.client = client
	}

	return
}

// DownloadFile downloads file by path.
func (client *Client) DownloadFile(
	ctx context.Context,
	path string,
) (io.ReadCloser, error) {
	return client.transport.Download(ctx, client.token, path)
}

// ProfilePhotosOptions contains options for method GetUserProfilePhotos.
type ProfilePhotosOptions struct {
	// Sequential number of the first photo to be returned.
	// By default, all photos are returned.
	Offset int

	// Limits the number of photos to be retrieved.
	// Values between 1—100 are accepted.
	// Defaults to 100.
	Limit int
}

func (opts *ProfilePhotosOptions) AddToRequest(r *Request) {
	if opts != nil {
		r.AddOptInt("offset", opts.Offset).
			AddOptInt("limit", opts.Limit)
	}
}

// GetUserProfilePhotos use this method to get a list of profile pictures for a user.
//
// Source: https://core.telegram.org/bots/api#getuserprofilephotos
func (client *Client) GetUserProfilePhotos(
	ctx context.Context,
	userID UserID,
	opts *ProfilePhotosOptions,
) (photos *UserProfilePhotos, err error) {
	req := NewRequest("getUserProfilePhotos").
		AddInt("user_id", int(userID)).
		AddPart(opts)

	err = client.Invoke(ctx,
		req,
		&photos,
	)

	return
}

// GetChat get up to date information about the chat.
//
// Source: https://core.telegram.org/bots/api#getchat
func (client *Client) GetChat(
	ctx context.Context,
	peer Peer,
) (chat *Chat, err error) {
	err = client.Invoke(ctx,
		NewRequest("getChat").AddChatID(peer),
		&chat,
	)

	return
}

// SetChatTitle use this method to change the title of a chat.
// Titles can't be changed for private chats.
// The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
//
// New chat title should be between 1-255 characters.
//
// Source: https://core.telegram.org/bots/api#setchattitle
func (client *Client) SetChatTitle(
	ctx context.Context,
	peer Peer,
	title string,
) error {
	return client.Invoke(ctx,
		NewRequest("setChatTitle").
			AddChatID(peer).
			AddString("title", title),
		nil,
	)
}

// SetChatDescription use this method to change the description of a supergroup or a channel.
// The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
//
// New chat description should be between 0-255 characters.
//
// Source: https://core.telegram.org/bots/api#setchatdescription
func (client *Client) SetChatDescription(
	ctx context.Context,
	peer Peer,
	description string,
) error {
	return client.Invoke(ctx,
		NewRequest("setChatDescription").
			AddChatID(peer).
			AddString("description", description),
		nil,
	)
}

// GetChatMembersCount returns numbers of members in chat.
//
// Source: https://core.telegram.org/bots/api#getchatmemberscount
func (client *Client) GetChatMembersCount(
	ctx context.Context,
	peer Peer,
) (count int, err error) {
	err = client.Invoke(ctx,
		NewRequest("getChatMembersCount").AddChatID(peer),
		&count,
	)

	return
}

// GetChatAdministrators use this method to get a list of administrators in a chat.
// On success, returns an Array of ChatMember objects that contains information
// about all chat administrators except other bots.
// If the chat is a group or a supergroup and no administrators were appointed,
// only the creator will be returned.
//
// Source: https://core.telegram.org/bots/api#getchatadministrators
func (client *Client) GetChatAdministrators(
	ctx context.Context,
	peer Peer,
) (admins ChatMemberSlice, err error) {
	err = client.Invoke(ctx,
		NewRequest("getChatAdministrators").AddChatID(peer),
		&admins,
	)

	return
}

// KickOptions contains optional options to kick.
type KickOptions struct {
	// Date when the user will be unbanned, unix time.
	// If user is banned for more than 366 days or less than 30 seconds
	// from the current time they are considered to be banned forever
	Until time.Time
}

func (opts *KickOptions) AddToRequest(r *Request) {
	if opts != nil {
		r.AddOptTime("until_date", opts.Until)
	}
}

// KickChatMember use this method to kick a user from a group, a supergroup or a channel.
// In the case of supergroups and channels, the user will not be able to return to the group on their own using invite links, etc.,
// unless unbanned first.
// The bot must be an administrator in the chat for this to work and must have the appropriate admin rights.
//
// Source: https://core.telegram.org/bots/api#kickchatmember
func (client *Client) KickChatMember(
	ctx context.Context,
	peer Peer,
	userID UserID,
	opts *KickOptions,
) error {
	return client.Invoke(ctx,
		NewRequest("kickChatMember").
			AddChatID(peer).
			AddInt("user_id", int(userID)).
			AddPart(opts),
		nil,
	)
}

// UnbanChatMember use this method to unban a previously kicked user in a supergroup or channel.
// The user will not return to the group or channel automatically, but will be able to join via link, etc.
// The bot must be an administrator for this to work.
//
// Source: https://core.telegram.org/bots/api#unbanchatmember
func (client *Client) UnbanChatMember(
	ctx context.Context,
	peer Peer,
	userID UserID,
) error {
	return client.Invoke(ctx,
		NewRequest("unbanChatMember").
			AddChatID(peer).
			AddInt("user_id", int(userID)),
		nil,
	)
}

// RestrictOptions contains optional options for restrict chat member method.
type RestrictOptions struct {
	// Date when restrictions will be lifted for the user, unix time.
	// If user is restricted for more than 366 days or less than 30 seconds from the current time,
	// they are considered to be restricted forever
	Until time.Time

	// Pass True, if the user can send text messages, contacts, locations and venues
	CanSendMessages bool

	// Pass True, if the user can send audios, documents, photos, videos, video notes and voice notes, implies can_send_messages
	CanSendMediaMessages bool

	// Pass True, if the user can send animations, games, stickers and use inline bots, implies can_send_media_messages
	CanSendOtherMessages bool

	// Pass True, if the user may add web page previews to their messages, implies can_send_media_messages
	CanSendWebPagePreviews bool
}

func (opts *RestrictOptions) AddToRequest(r *Request) {
	if opts != nil {
		r.AddOptTime("until_date", opts.Until).
			AddOptBool("can_send_messages", opts.CanSendMessages).
			AddOptBool("can_send_media_messages", opts.CanSendMediaMessages).
			AddOptBool("can_send_other_messages", opts.CanSendOtherMessages).
			AddOptBool("can_send_web_page_previews", opts.CanSendWebPagePreviews)
	}
}

// RestrictChatMember use this method to restrict a user in a supergroup.
// The bot must be an administrator in the supergroup for this to work and must have the appropriate admin rights.
// Pass True for all boolean parameters to lift restrictions from a user.
//
// Source: https://core.telegram.org/bots/api#restrictchatmember
func (client *Client) RestrictChatMember(
	ctx context.Context,
	peer Peer,
	userID UserID,
	opts *RestrictOptions,
) error {
	return client.Invoke(ctx,
		NewRequest("restrictChatMember").
			AddChatID(peer).
			AddInt("user_id", int(userID)).
			AddPart(opts),
		nil,
	)
}

type OutgoingMessage interface {
	BuildSendRequest() (*Request, error)
}

func (client *Client) Send(ctx context.Context, msg OutgoingMessage, dst interface{}) error {
	req, err := msg.BuildSendRequest()
	if err != nil {
		return err
	}

	return client.Invoke(ctx,
		req,
		&dst,
	)
}

type UpdatesOptions struct {
	// Identifier of the first update to be returned.
	// Must be greater by one than the highest among the identifiers of previously received updates.
	// By default, updates starting with the earliest unconfirmed update are returned.
	// An update is considered confirmed as soon as getUpdates is called with an offset higher than its update_id.
	// The negative offset can be specified to retrieve updates starting from -offset update from the end of the updates queue.
	// All previous updates will forgotten.
	Offset UpdateID

	// Limits the number of updates to be retrieved.
	// Values between 1—100 are accepted.
	// Defaults to 100.
	Limit int

	// Timeout for long polling.
	// Defaults to 0, i.e. usual short polling.
	// Should be positive, short polling should be used for testing purposes only.
	Timeout time.Duration

	// List the types of updates you want your bot to receive.
	// Specify an empty list to receive all updates regardless of type (default).
	// If not specified, the previous setting will be used.
	AllowedUpdates []UpdateType
}

func (opts *UpdatesOptions) addToRequestAllowedUpdates(r *Request) error {
	if opts.AllowedUpdates != nil {
		val, err := json.Marshal(opts.AllowedUpdates)
		if err != nil {
			return errors.Wrap(err, "marshal allowed_updates")
		}
		r.AddString("allowed_updates", string(val))
	}
	return nil
}

func (opts *UpdatesOptions) addToRequest(r *Request) error {
	if opts != nil {
		r.AddOptInt("offset", int(opts.Offset)).
			AddOptInt("limit", int(opts.Limit)).
			AddOptInt("timeout", int(opts.Timeout.Seconds()))

		return opts.addToRequestAllowedUpdates(r)
	}

	return nil
}

// GetUpdates returns incoming updates using long polling.
func (client *Client) GetUpdates(ctx context.Context, opts *UpdatesOptions) (updates UpdateSlice, err error) {
	r := NewRequest("getUpdates")

	if err := opts.addToRequest(r); err != nil {
		return nil, err
	}

	err = client.Invoke(ctx, r, &updates)

	return
}

// WebhookOptions contains optional params for Client.SetWebhook.
type WebhookOptions struct {
	// Optional. Upload your public key certificate so that the root certificate in use can be checked.
	Certificate *InputFile

	// Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100.
	// Defaults to 40.
	// Use lower values to limit the load on your bot‘s server, and higher values to increase your bot’s throughput.
	MaxConnections int

	// List the types of updates you want your bot to receive.
	// Specify an empty list to receive all updates regardless of type (default).
	// If not specified, the previous setting will be used.
	AllowedUpdates []UpdateType
}

func (opts *WebhookOptions) addToRequestAllowedUpdates(r *Request) error {
	if opts.AllowedUpdates != nil {
		val, err := json.Marshal(opts.AllowedUpdates)
		if err != nil {
			return errors.Wrap(err, "marshal allowed_updates")
		}
		r.AddString("allowed_updates", string(val))
	}
	return nil
}

func (opts *WebhookOptions) addToRequest(r *Request) error {
	if opts != nil {
		if opts.Certificate != nil {
			r.AddFile("certificate", *opts.Certificate)
		}

		r.AddOptInt("max_connections", opts.MaxConnections)

		return opts.addToRequestAllowedUpdates(r)
	}

	return nil
}

// SetWebhook use this method to specify a url and receive incoming updates via an outgoing webhook.
// Whenever there is an update for the bot, we will send an HTTPS POST request to the specified url, containing a JSON-serialized Update.
// In case of an unsuccessful request, we will give up after a reasonable amount of attempts. Returns True on success.
func (client *Client) SetWebhook(ctx context.Context, url string, opts *WebhookOptions) error {
	r := NewRequest("setWebhook")

	r.AddString("url", url)

	if err := opts.addToRequest(r); err != nil {
		return err
	}

	return client.Invoke(ctx, r, nil)
}

// GetWebhookInfo returns current webhook status.
//
// Source: https://core.telegram.org/bots/api#getwebhookinfo
func (client *Client) GetWebhookInfo(
	ctx context.Context,
) (info *WebhookInfo, err error) {
	// TODO: maybe define requests without argument globally and use it instead create new?

	err = client.Invoke(ctx,
		NewRequest("getWebhookInfo"),
		&info,
	)

	return
}

// DeleteWebhook remove webhook.
//
// Source: https://core.telegram.org/bots/api#deletewebhook
func (client *Client) DeleteWebhook(ctx context.Context) error {
	return client.Invoke(ctx,
		NewRequest("deleteWebhook"),
		nil,
	)
}
