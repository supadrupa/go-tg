package tg

import (
	"context"
	"io"
	"time"
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

type clientMixin interface {
	WithClient(client *Client)
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

	// handle bot api errors here

	if dst == nil {
		return nil
	}

	// unmarshal response
	if err := res.UnmarshalResult(&dst); err != nil {
		return err
	}

	return nil
}

// GetMe returns bot profile.
func (client *Client) GetMe(
	ctx context.Context,
) (user *User, err error) {
	err = client.Invoke(ctx,
		NewRequest("getMe"),
		&user,
	)

	return
}

func (client *Client) GetFile(
	ctx context.Context,
	id FileID,
) (*File, error) {
	file := &File{client: client}

	req := NewRequest("getFile").AddString("file_id", string(id))

	if err := client.Invoke(ctx,
		req,
		&file,
	); err != nil {
		return nil, err
	}

	return file, nil
}

// DownloadFile downloads file by path.
func (client *Client) DownloadFile(
	ctx context.Context,
	path string,
) (io.ReadCloser, error) {
	return client.transport.Download(ctx, client.token, path)
}

type ProfilePhotosOptions struct {
	// Sequential number of the first photo to be returned.
	// By default, all photos are returned.
	Offset int

	// Limits the number of photos to be retrieved.
	// Values between 1—100 are accepted.
	// Defaults to 100.
	Limit int
}

// GetUserProfilePhotos use this method to get a list of profile pictures for a user.
func (client *Client) GetUserProfilePhotos(
	ctx context.Context,
	userID UserID,
	opts *ProfilePhotosOptions,
) (photos *UserProfilePhotos, err error) {

	req := NewRequest("getUserProfilePhotos").
		AddInt("user_id", int(userID))

	if opts != nil {
		req.AddOptInt("offset", opts.Offset)
		req.AddOptInt("limit", opts.Limit)
	}

	err = client.Invoke(ctx,
		req,
		&photos,
	)

	return
}

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

// Use this method to change the title of a chat.
// Titles can't be changed for private chats.
// The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
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

// Use this method to change the description of a supergroup or a channel.
// The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
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

// Use this method to get a list of administrators in a chat.
// On success, returns an Array of ChatMember objects that contains information
// about all chat administrators except other bots.
// If the chat is a group or a supergroup and no administrators were appointed,
// only the creator will be returned.
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
	UntilDate time.Time
}

func (opts *KickOptions) AddToRequest(r *Request) {
	if opts != nil {
		r.AddOptTime("until_date", opts.UntilDate)
	}
}

// KickChatMember use this method to kick a user from a group, a supergroup or a channel.
// In the case of supergroups and channels, the user will not be able to return to the group on their own using invite links, etc.,
// unless unbanned first.
//  The bot must be an administrator in the chat for this to work and must have the appropriate admin rights.
// Returns True on success.
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
// The bot must be an administrator for this to work. Returns True on success.
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
	UntilDate time.Time

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
	if opts == nil {
		return
	}

	r.AddOptTime("until_date", opts.UntilDate).
		AddOptBool("can_send_messages", opts.CanSendMessages).
		AddOptBool("can_send_media_messages", opts.CanSendMediaMessages).
		AddOptBool("can_send_other_messages", opts.CanSendOtherMessages).
		AddOptBool("can_send_web_page_previews", opts.CanSendWebPagePreviews)

}

// RestrictChatMember use this method to restrict a user in a supergroup.
// The bot must be an administrator in the supergroup for this to work and must have the appropriate admin rights.
// Pass True for all boolean parameters to lift restrictions from a user.
// Returns True on success.
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
	SendRequest() (*Request, error)
}

func (client *Client) Send(
	ctx context.Context,
	msg OutgoingMessage,
	result interface{},
) error {
	return nil
}
