package models

type Language string

const (
	LangUz Language = "uz"
	LangRu Language = "ru"
	LangEn Language = "en"
)

// MultiLang represents standard multilingual text in 3 supported languages.
type MultiLang struct {
	Uz string `json:"uz" validate:"required,min=1,max=255"`
	Ru string `json:"ru,omitempty" validate:"required,omitempty,max=255"`
	En string `json:"en,omitempty" validate:"omitempty,max=255"`
}

// Resolve extracts the text for the requested language with fallback to Uzbek.
func (m MultiLang) Resolve(lang Language) string {
	switch lang {
	case LangRu:
		if m.Ru != "" {
			return m.Ru
		}
	case LangEn:
		if m.En != "" {
			return m.En
		}
	case LangUz:
		if m.Uz != "" {
			return m.Uz
		}
	}

	if m.Uz != "" {
		return m.Uz
	}
	if m.Ru != "" {
		return m.Ru
	}
	return m.En
}

// ToMap converts the MultiLang to a map[string]string.
func (m MultiLang) ToMap() map[string]string {
	return map[string]string{"uz": m.Uz, "ru": m.Ru, "en": m.En}
}

// ToMapOrNil converts the MultiLang to a map[string]string or nil.
func (m *MultiLang) ToMapOrNil() map[string]string {
	if m == nil {
		return nil
	}
	return m.ToMap()
}

// MultiLangFromMap converts a map[string]string to a MultiLang.
func MultiLangFromMap(m map[string]string) *MultiLang {
	if m == nil {
		return nil
	}
	return &MultiLang{Uz: m["uz"], Ru: m["ru"], En: m["en"]}
}

// MultiLangText represents longer multilingual descriptions.
type MultiLangText struct {
	Uz string `json:"uz" validate:"required,min=1,max=2048"`
	Ru string `json:"ru,omitempty" validate:"required,omitempty,max=2048"`
	En string `json:"en,omitempty" validate:"omitempty,max=2048"`
}

// Resolve extracts the text for the requested language with fallback to Uzbek.
func (m MultiLangText) Resolve(lang Language) string {
	switch lang {
	case LangRu:
		if m.Ru != "" {
			return m.Ru
		}
	case LangEn:
		if m.En != "" {
			return m.En
		}
	case LangUz:
		if m.Uz != "" {
			return m.Uz
		}
	}

	if m.Uz != "" {
		return m.Uz
	}
	if m.Ru != "" {
		return m.Ru
	}
	return m.En
}

// ToMap converts the MultiLangText to a map[string]string.
func (m MultiLangText) ToMap() map[string]string {
	return map[string]string{"uz": m.Uz, "ru": m.Ru, "en": m.En}
}

// ToMapOrNil converts the MultiLangText to a map[string]string or nil.
func (m *MultiLangText) ToMapOrNil() map[string]string {
	if m == nil {
		return nil
	}
	return m.ToMap()
}

// MultiLangTextFromMap converts a map[string]string to a MultiLangText.
func MultiLangTextFromMap(m map[string]string) *MultiLangText {
	if m == nil {
		return nil
	}
	return &MultiLangText{Uz: m["uz"], Ru: m["ru"], En: m["en"]}
}
