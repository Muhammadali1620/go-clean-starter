package dto

// MultiLang represents standard multilingual text in 3 supported languages.
type MultiLang struct {
	Uz string `json:"uz" validate:"required,min=1,max=255"`
	Ru string `json:"ru" validate:"required,min=1,max=255"`
	En string `json:"en,omitempty" validate:"omitempty,max=255"`
}

// MultiLangText represents longer multilingual descriptions.
type MultiLangText struct {
	Uz string `json:"uz" validate:"required,min=1"`
	Ru string `json:"ru" validate:"required,min=1"`
	En string `json:"en,omitempty"`
}
