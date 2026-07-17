package models

type Settings struct {
	KEY   string         `bun:"id,pk,notnull"`
	Value map[string]any `bun:"options,type:json"`
}
