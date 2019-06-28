package tg

import (
	"context"
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
