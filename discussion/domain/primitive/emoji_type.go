package primitive

type EmojiType interface {
	EmojiType() int
}

func NewEmojiType(v int) (EmojiType, error) {
	return emojiType(v), nil
}

func CreateEmojiType(v int) EmojiType {
	return emojiType(v)
}

type emojiType int

func (t emojiType) EmojiType() int {
	return int(t)
}

func (t emojiType) IsEqual() {

}
