package tg

// Sticker object represents a sticker.
type Sticker struct {
	// Unique identifier for this file
	FileID FileID `json:"file_id"`

	// Sticker width
	Width int `json:"width"`

	// Sticker height
	Height int `json:"height"`

	// Optional. Sticker thumbnail in the .webp or .jpg format
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Optional. Emoji associated with the sticker
	Emoji string `json:"emoji,omitempty"`

	// Optional. Name of the sticker set to which the sticker belongs
	SetName string `json:"set_name,omitempty"`

	// Optional. For mask stickers, the position where the mask should be placed
	MastPosition *MaskPosition `json:"mast_position,omitempty"`
}

// MaskPosition object describes the position on faces where a mask should be placed by default.
type MaskPosition struct {
	// The part of the face relative to which the mask should be placed. One of “forehead”, “eyes”, “mouth”, or “chin”.
	Point string `json:"point"`

	// Shift by X-axis measured in widths of the mask scaled to the face size, from left to right. For example, choosing -1.0 will place mask just to the left of the default mask position.
	XShift float64 `json:"x_shift"`

	// Shift by Y-axis measured in heights of the mask scaled to the face size, from top to bottom. For example, 1.0 will place the mask just below the default mask position.
	YShift float64 `json:"y_shift"`

	// Mask scaling coefficient. For example, 2.0 means double size.
	Scale float64 `json:"scale"`
}
