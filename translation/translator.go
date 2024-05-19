package translation

type Translator interface {
	translate(string) (string, error)
}

