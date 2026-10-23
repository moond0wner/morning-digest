package invest

import (
	"fmt"
	"strings"
)

type Portfolio struct {
	TotalAmount   float64    // общая стоимость (в рублях)
	ExpectedYield float64    // общая доходность за всё время (в рублях)
	DailyYield    float64    // доходность за день (в рублях)
	Positions     []Position // список позиций (акций)
}

type Position struct {
	Ticker        string
	Quantity      float64
	CurrentPrice  float64
	ExpectedYield float64
}

func (p *Portfolio) String() string {
	// TODO: Сделать лучше
	var sb strings.Builder
	sb.WriteString("Портфель\n\n")
	fmt.Fprintf(&sb, "Стоимость: %.0f ₽\n", p.TotalAmount)
	fmt.Fprintf(&sb, "Доходность: %+.3f ₽\n", p.ExpectedYield)
	fmt.Fprintf(&sb, "За день: %+.2f ₽\n\n", p.DailyYield)

	if len(p.Positions) == 0 {
		sb.WriteString("Нет открытых позиций\n")
	}

	for _, pos := range p.Positions {
		fmt.Fprintf(
			&sb,
			"• %s — %.0f шт, %.2f ₽, %+.0f ₽\n",
			pos.Ticker, pos.Quantity, pos.CurrentPrice, pos.ExpectedYield,
		)
	}
	return sb.String()
}
