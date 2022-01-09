package domain

import "github.com/cs-sysimpl/suzukake/domain/values"

type Choice struct {
	id    values.ChoiceID
	label values.ChoiceLabel
}

func NewChoice(
	id values.ChoiceID,
	label values.ChoiceLabel,
) Choice {
	return Choice{
		id:    id,
		label: label,
	}
}

func (c *Choice) GetID() values.ChoiceID {
	return c.id
}

func (c *Choice) GetLabel() values.ChoiceLabel {
	return c.label
}
