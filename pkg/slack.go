package slack

import "encoding/json"

const (
	markdown = "mrkdwn"
)

// responseToSlack assists in formatting requests to make to Slack.
type responseToSlack struct {
	// Fallback text for if the client doesn't support blocks or attachments, or if you just want to send
	// a small text snippet.
	Text string `json:"text"`
	// The type of message to send. Choices are "ephemeral", and "in_channel".
	// "ephemeral" will be sent only to "you", whereas "in_channel" also goes to the channel it was
	// summoned from.
	ResponseType string `json:"response_type"`
	// Deprecated field, don't use it unless you really have to.
	Attachments []Attachment `json:"attachments,omitempty"`
	// Read more about blocks here https://api.slack.com/messaging/composing/layouts
	Blocks []Block `json:"blocks,omitempty"`
}

type ResponseType string

const (
	inChannel ResponseType = "in_channel"
	ephemeral ResponseType = "ephemeral"
)

// RenderSlackJSON marshals a response to Slack. A fallback should be provided in the case of older Slack clients or something?
func (m *blockBuilder) RenderSlackJSON() string {
	resp, _ := json.Marshal(responseToSlack{m.text, string(m.responseType), nil, m.Build()})
	return string(resp)
}

type Attachment struct {
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Text struct {
	Text string `json:"text"`
	// "mrkdwn" or "plain_text"
	Type string `json:"type"`
	// Whether to render emoji or not in the message; only valid for "plain_text" types
	Emoji bool `json:"emoji,omitempty"`
}

type Button struct {
	// Should be "button"
	Type     string `json:"type"`
	Text     Text   `json:"text"`
	ActionID string `json:"action_id,omitempty"`
	URL      string `json:"url"`
	Style    string `json:"style,omitempty"`
}

type Accessory interface{}

// A type of Accessory
type Image struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

type Field interface{}

type Block interface{}

// https://api.slack.com/reference/messaging/block-elements
// Block kit builder https://api.slack.com/tools/block-kit-builder

type sectionBlock struct {
	Type      string    `json:"type"`
	BlockID   string    `json:"block_id"`
	Text      Text      `json:"text,omitempty"`
	Fields    []Field   `json:"fields,omitempty"`
	Accessory Accessory `json:"accessory,omitempty"`
}

type contextBlock struct {
	Type     string  `json:"type"`
	Elements []Field `json:"elements"`
	BlockID  string  `json:"block_id"`
}

type actionsBlock struct {
	Type     string  `json:"type"`
	Elements []Field `json:"elements"`
	BlockID  string  `json:"block_id"`
}

type dividerBlock struct {
	Type string `json:"type"`
}

type imageBlock struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
	Title    string `json:"title,omitempty"`
	BlockID  string `json:"block_id,omitempty"`
}

type blockBuilder struct {
	// Fallback text for if the client doesn't support blocks or attachments, or if you just want to send
	// a small text snippet.
	text string
	// The type of message to send. Choices are "ephemeral", and "in_channel".
	// "ephemeral" will be sent only to "you", whereas "in_channel" also goes to the channel it was
	// summoned from.
	responseType ResponseType
	// Read more about blocks here https://api.slack.com/messaging/composing/layouts
	blocks []Block
}

func BlockBuilder() *blockBuilder {
	return &blockBuilder{
		responseType: ephemeral,
	}
}

// SendToChannel will tell Slack to broadcast this message in the channel, rather than having it be an "ephemeral"
// message, which is one that will only show for the invoker of the command.
func (m *blockBuilder) SendToChannel() *blockBuilder {
	m.responseType = inChannel
	return m
}

// Create a new SectionBlock and add it to the slice of Blocks. blockID is not technically required, and can be omitted,
// but it's best for future enhancement possibilities if you include it. Thus it is not an optional parameter and must
// me manually omitted.
func (m *blockBuilder) SectionBlock(blockID, text string, accessory Accessory, fields ...Field) *blockBuilder {
	var iFields []Field
	for _, f := range fields {
		iFields = append(iFields, f)
	}
	b := sectionBlock{
		Type:    "section",
		BlockID: blockID,
		Text: Text{
			Type: Markdown,
			Text: text,
		},
		Fields:    iFields,
		Accessory: accessory,
	}
	m.blocks = append(m.blocks, b)
	return m
}

func (m *blockBuilder) DividerBlock() *blockBuilder {
	b := dividerBlock{
		Type: "divider",
	}
	m.blocks = append(m.blocks, b)
	return m
}

// ContextBlock - Displays message context, which can include both images and text.
func (m *blockBuilder) ContextBlock(blockID string, element Field, elements ...Field) *blockBuilder {
	iElements := []Field{element}
	for _, e := range elements {
		iElements = append(iElements, e)
	}
	b := contextBlock{
		BlockID:  blockID,
		Type:     "context",
		Elements: iElements,
	}
	m.blocks = append(m.blocks, b)
	return m
}

// ActionsBlock - A block that is used to hold interactive elements.
func (m *blockBuilder) ActionsBlock(blockID string, element Field, elements ...Field) *blockBuilder {
	iElements := []Field{element}
	for _, e := range elements {
		iElements = append(iElements, e)
	}
	b := actionsBlock{
		BlockID:  blockID,
		Type:     "actions",
		Elements: iElements,
	}
	m.blocks = append(m.blocks, b)
	return m
}

func (m *blockBuilder) ImageBlock(blockID, imageURL, altText, title string) *blockBuilder {
	b := imageBlock{
		BlockID:  blockID,
		Type:     "image",
		ImageURL: imageURL,
		AltText:  altText,
		Title:    title,
	}
	m.blocks = append(m.blocks, b)
	return m
}

func (m *blockBuilder) WithFallback(text string) *blockBuilder {
	m.text = text
	return m
}

// Build returns the slice of all the Blocks you've setup so far.
func (m *blockBuilder) Build() []Block {
	return m.blocks
}
