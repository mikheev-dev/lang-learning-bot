package translation

type TermTranslation struct {
	Pos  string
	Vars []string
}

type TranslatedTerm struct {
	Text         string
	Translations []TermTranslation
}

type TranslatedSequence struct {
	Text        string
	Translation string
}

type SequenceTranslator interface {
	Translate(string) (*TranslatedSequence, error)
}

type TermTranslator interface {
	Translate(string) (*TranslatedTerm, error)
}
