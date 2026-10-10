package invest

import (
	"fmt"
	"sort"
	"strings"
)

type Portfolio struct {
	TotalAmount   float64 // общая стоимость
	ExpectedYield float64 // доходность (в валюте)
	DailyYield    float64 // доходность за день
	Positions     []Position
}

type Position struct {
	Ticker        string
	Quantity      float64
	CurrentPrice  float64
	ExpectedYield float64
}

func (p *Portfolio) String() string {
	sort.Slice(p.Positions, func(i, j int) bool {
		return p.Positions[i].CurrentPrice*p.Positions[i].Quantity >
			p.Positions[j].CurrentPrice*p.Positions[j].Quantity
	})
	var sb strings.Builder
	sb.WriteString("Портфель\n\n")
	sb.WriteString(fmt.Sprintf("Стоимость: %.0f RUB\n", p.TotalAmount))
	sb.WriteString(fmt.Sprintf("Доходность: %.0f RUB\n", p.ExpectedYield))
	sb.WriteString(fmt.Sprintf("За день: %.0f RUB\n\n", p.DailyYield))
	for _, position := range p.Positions {
		sb.WriteString(fmt.Sprintf(
			"• %s — %.0f шт, %.0f ₽, %+.0f ₽\n",
			position.Ticker, position.Quantity, position.CurrentPrice, position.ExpectedYield,
		))
	}
	return sb.String()
}
