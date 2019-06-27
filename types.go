package tg

import (
	"context"
	"encoding/json"
	"io"
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

// Peer define generic interface
type Peer interface {
	AddPeerToRequest(k string, r *Request)
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

// FileID represents unique file identifier.
type FileID string

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
