package adapter

type AttachmentType string

const (
	AttachmentTypeText   AttachmentType = "text"
	AttachmentTypeImage  AttachmentType = "image"
	AttachmentTypeFolder AttachmentType = "folder"
)

type Attachment struct {
	Name     string
	Path     string // Original file path if available
	Data     string // Base64 if image, plain text if text
	Type     AttachmentType
	MimeType string
}

type Provider interface {
	Name() string
	Generate(prompt string, attachments []Attachment) (string, error)
	GenerateStream(prompt string, attachments []Attachment, tokens chan<- string) error
	SetModel(model string)
	SetSystemPrompt(prompt string)
}
