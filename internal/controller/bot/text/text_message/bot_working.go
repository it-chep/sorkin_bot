package text_message

type TextBotMessage struct {
}

func NewTextBotMessage() TextBotMessage {
	return TextBotMessage{}
}

func (c TextBotMessage) Execute() {

}

func (c TextBotMessage) validatePhoneMessage() (valid bool) {
	return true
}
