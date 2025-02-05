package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// ParseMode represents formatting style
type ParseMode int8

const (
	// Plain it's no formatting
	Plain ParseMode = iota
	// Markdown formatting style
	Markdown
	// HTML formatting style
	HTML
)

// String returns string representation of parse mode.
// If not set returns empty string.
func (pm ParseMode) String() string {
	switch pm {
	case Markdown:
		return "markdown"
	case HTML:
		return "HTML"
	default:
		return ""
	}
}

// Peer define generic interface
type Peer interface {
	AddPeerToRequest(k string, r *Request)
}

func ParsePeer(v string) (Peer, error) {
	if strings.HasPrefix(v, "@") {
		return Username(v[1:]), nil
	} else {
		// NOTE: maybe check prefix (e.g. -100 is channels) and return and a more specific type?

		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		return ChatID(id), nil
	}
}

// UserID represents unique user identifier.
type UserID int

func (id UserID) AddPeerToRequest(k string, r *Request) { r.AddInt(k, int(id)) }

// Username represents user/supergroup/channel username.
type Username string

func (un Username) AddPeerToRequest(k string, r *Request) { r.AddString(k, string("@"+un)) }

// User represents Telegram user or bot.
type User struct {
	ID           UserID   `json:"id"`
	IsBot        bool     `json:"is_bot"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name,omitempty"`
	Username     Username `json:"username,omitempty"`
	LanguageCode string   `json:"language_code,omitempty"`
}

func (user User) AddPeerToRequest(k string, r *Request) { user.ID.AddPeerToRequest(k, r) }

// UserSlice it's just alias for slice of users.
type UserSlice []User

// FileID represents unique file identifier.
type FileID string

func (id FileID) AddFileToRequest(k string, r *Request) {
	r.AddString(k, string(id))
}

// File represents a file ready to be downloaded.
type File struct {
	ID   FileID `json:"file_id"`
	Size int    `json:"file_size,omitempty"`
	Path string `json:"file_path,omitempty"`

	client *Client
}

// NewReader creates a file content reader
func (file File) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return file.client.DownloadFile(ctx, file.Path)
}

// PhotoSize represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	// Unique identifier for this file
	FileID FileID `json:"file_id"`

	// Photo width
	Width int `json:"width"`

	// Photo height
	Height int `json:"height"`

	// Optional. File size
	FileSize int `json:"file_size,omitempty"`
}

// PhotoSizeSlice represents array of PhotoSize
type PhotoSizeSlice []PhotoSize

// UserProfilePhotos object represent a user's profile pictures.
type UserProfilePhotos struct {
	Total int              `json:"total_count"`
	Items []PhotoSizeSlice `json:"photos"`
}

// First returns first photo or nil if no items.
func (photos UserProfilePhotos) First() PhotoSizeSlice {
	if len(photos.Items) > 0 {
		return photos.Items[0]
	}

	return nil
}

// Last returns last photo or nil if no items.
func (photos UserProfilePhotos) Last() PhotoSizeSlice {
	if len(photos.Items) > 0 {
		return photos.Items[len(photos.Items)-1]
	}

	return nil
}

type ChatPhoto struct {
	// Unique file identifier of small (160x160) chat photo.
	// This file_id can be used only for photo download.
	SmallFileID FileID `json:"small_file_id"`

	// Unique file identifier of big (640x640) chat photo.
	// This file_id can be used only for photo download.
	BigFileID FileID `json:"big_file_id"`
}

// ChatID unique chat identifier.
type ChatID int64

func (id ChatID) AddPeerToRequest(k string, r *Request) { r.AddInt64(k, int64(id)) }

type ChatType string

const (
	PrivateChat    = ChatType("private")
	GroupChat      = ChatType("group")
	SupergroupChat = ChatType("supergroup")
	ChannelChat    = ChatType("channel")
)

// Chat represents a generic type for all chats.
type Chat struct {
	// Unique identifier for this chat.
	ID ChatID `json:"id"`

	// Type of chat
	Type ChatType `json:"type"`

	// Optional. Title, for supergroups, channels and group chats
	Title string `json:"title,omitempty"`

	// Optional. Username, for private chats, supergroups and channels if available
	Username Username `json:"username,omitempty"`

	// Optional. First name of the other party in a private chat
	FirstName string `json:"first_name,omitempty"`

	// Optional. Last name of the other party in a private chat
	LastName string `json:"last_name,omitempty"`

	// Optional. True if a group has ‘All Members Are Admins’ enabled.
	AllMembersAreAdministrators bool `json:"all_members_are_administrators,omitempty"`

	// Optional. Chat photo.
	// Returned only in GetChat.
	Photo *ChatPhoto `json:"photo,omitempty"`

	// Optional. Description, for supergroups and channel chats.
	// Returned only in GetChat.
	Description string `json:"description,omitempty"`

	// Optional. Chat invite link, for supergroups and channel chats.
	// Each administrator in a chat generates their own invite links,
	// so the bot must first generate the link using exportChatInviteLink.
	// Returned only in getChat.
	InviteLink string `json:"invite_link,omitempty"`

	// Optional. Pinned message, for groups, supergroups and channels.
	// Returned only in getChat.
	PinnedMessage json.RawMessage `json:"pinned_message,omitempty"`

	// Optional. For supergroups, name of group sticker set.
	// Returned only in getChat.
	StickerSetName string `json:"sticker_set_name,omitempty"`

	// Optional. True, if the bot can change the group sticker set.
	// Returned only in getChat.
	CanSetStickerSet bool `json:"can_set_sticker_set,omitempty"`
}

func (chat Chat) AddPeerToRequest(k string, r *Request) { chat.ID.AddPeerToRequest(k, r) }

type ChatMember struct {
	// Information about the user
	User User `json:"user"`

	// The member's status in the chat. Can be “creator”, “administrator”, “member”, “restricted”, “left” or “kicked”.
	Status string `json:"status"`

	// Optional. Restricted and kicked only.
	// Date when restrictions will be lifted for this user, unix time.
	UntilDate int64 `json:"until_date,omitempty"`

	// Optional. Administrators only.
	// True, if the bot is allowed to edit administrator privileges of that user.
	CanBeEdited bool `json:"can_be_edited,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can change the chat title, photo and other settings.
	CanChangeInfo bool `json:"can_change_info,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can post in the channel, channels only.
	CanPostMessages bool `json:"can_post_messages,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can edit messages of other users
	// and can pin messages, channels only.
	CanEditMessages bool `json:"can_edit_messages,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can delete messages of other users
	CanDeleteMessages bool `json:"can_delete_messages,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can invite new users to the chat.
	CanInviteUsers bool `json:"can_invite_users,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can restrict, ban or unban chat members.
	CanRestrictMembers bool `json:"can_restrict_members,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can pin messages, groups and supergroups only
	CanPinMessages bool `json:"can_pin_messages,omitempty"`

	// Optional. Administrators only.
	// True, if the administrator can add new administrators with a subset
	// of his own privileges or demote administrators that he has promoted,
	// directly or indirectly (promoted by administrators that were appointed by the user).
	CanPromoteMembers bool `json:"can_promote_members,omitempty"`

	// Optional. Restricted only. True, if the user is a member of the chat at the moment of the request
	IsMember bool `json:"is_member,omitempty"`

	// Optional. Restricted only. True, if the user can send text messages, contacts, locations and venues
	CanSendMessages bool `json:"can_send_messages,omitempty"`

	// Optional. Restricted only.
	// True, if the user can send audios, documents, photos, videos, video notes and voice notes,
	// implies can_send_messages
	CanSendMediaMessages bool `json:"can_send_media_messages,omitempty"`

	// Optional. Restricted only. True, if the user can send animations, games, stickers and use inline bots, i
	// implies can_send_media_messages
	CanSendOtherMessages bool `json:"can_send_other_messages,omitempty"`
}

// ChatMemberSlice define a array of chat members
type ChatMemberSlice []ChatMember

// CallbackQueryID represents unique CallbackQuery identifier.
type CallbackQueryID string

// CallbackQuery object represents an incoming callback query from a callback button in an inline keyboard.
// If the button that originated the query was attached to a message sent by the bot, the field message will be present.
// If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present.
// Exactly one of the fields data or game_short_name will be present.
type CallbackQuery struct {
	// Unique identifier for this query
	ID CallbackQueryID `json:"id"`

	// Sender
	From User `json:"from"`

	// Optional. Message with the callback button that originated the query.
	// Note that message content and message date will not be available if the message is too old.
	Message *Message `json:"message,omitempty"`

	// Optional. Identifier of the message sent via the bot in inline mode, that originated the query.
	InlineMessageID InlineMessageID `json:"inline_message_id,omitempty"`

	// Global identifier, uniquely corresponding to the chat to which the message with the callback button was sent.
	// Useful for high scores in games.
	ChatInstance string `json:"chat_instance,omitempty"`

	// Optional. Data associated with the callback button.
	// Be aware that a bad client can send arbitrary data in this field.
	Data string `json:"data,omitempty"`

	// Optional. Short name of a Game to be returned, serves as the unique identifier for the game
	GameShortName string `json:"game_short_name,omitempty"`
}

// WebhookError represent error that happened when trying to delivery update via webhook.
type WebhookError struct {
	// Description of error
	Message string

	// Date of error
	Date time.Time
}

func (err *WebhookError) Error() string {
	return fmt.Sprintf("%s at %s (%s ago)",
		err.Message,
		err.Date.Format(time.RFC850),
		time.Since(err.Date),
	)
}

// WebhookInfo contains information about the current status of a webhook.
type WebhookInfo struct {
	// WebhookInfo URL, may be empty if webhook is not set up
	URL string

	// True, if a custom certificate was provided for webhook certificate checks
	HasCustomCertificate bool

	// Most recent error happened when trying to delivery update via webhook.
	Error *WebhookError

	// Number of updates awaiting delivery
	PendingUpdateCount int

	// Optional. Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery
	MaxConnections int

	// Optional. A list of update types the bot is subscribed to. Defaults to all update types
	AllowedUpdates []UpdateType
}

// IsSet returns true if webhook is set up.
func (webhook WebhookInfo) IsSet() bool {
	return webhook.URL != ""
}

// HasError returns true if webhook has information about last error.
func (webhook WebhookInfo) HasError() bool {
	return webhook.Error != nil
}

func (webhook *WebhookInfo) UnmarshalJSON(data []byte) error {
	response := struct {
		URL                  string       `json:"url"`
		HasCustomCertificate bool         `json:"has_custom_certificate"`
		LastErrorMessage     string       `json:"last_error_message"`
		LastErrorDate        int64        `json:"last_error_date"`
		PendingUpdateCount   int          `json:"pending_update_count"`
		MaxConnections       int          `json:"max_connections"`
		AllowedUpdates       []UpdateType `json:"allowed_updates"`
	}{}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	var lastError *WebhookError

	if response.LastErrorMessage != "" {
		lastError = &WebhookError{
			Message: response.LastErrorMessage,
			Date:    time.Unix(response.LastErrorDate, 0),
		}
	}

	*webhook = WebhookInfo{
		URL:                  response.URL,
		HasCustomCertificate: response.HasCustomCertificate,
		PendingUpdateCount:   response.PendingUpdateCount,
		Error:                lastError,
		MaxConnections:       response.MaxConnections,
		AllowedUpdates:       response.AllowedUpdates,
	}

	return nil
}
