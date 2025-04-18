package models

type WireSetGenTemplateModel struct {
	PackageName string
	Imports     []*ImportTemplate
	WireSets    []*WireSet
}
