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
	var sb strings.Builder

	yieldEmoji := "🟢"
	if p.ExpectedYield < 0 {
		yieldEmoji = "🔴"
	}

	dailyEmoji := "📈"
	if p.DailyYield < 0 {
		dailyEmoji = "📉"
	}

	sb.WriteString("💼 Портфель\n\n")
	fmt.Fprintf(&sb, "Стоимость: %.0f ₽\n", p.TotalAmount)
	fmt.Fprintf(&sb, "%s Доходность: %+.0f ₽\n", yieldEmoji, p.ExpectedYield)
	fmt.Fprintf(&sb, "%s За день: %+.0f ₽\n\n", dailyEmoji, p.DailyYield)

	if len(p.Positions) == 0 {
		sb.WriteString("_Нет открытых позиций_ 🕳")
		return sb.String()
	}

	for _, pos := range p.Positions {
		posEmoji := "🟩"
		if pos.ExpectedYield < 0 {
			posEmoji = "🟥"
		}

		fmt.Fprintf(
			&sb,
			"%s %s\n   %.0f шт × %.2f ₽ = %+.0f ₽\n",
			posEmoji, pos.Ticker, pos.Quantity, pos.CurrentPrice, pos.ExpectedYield,
		)
	}
	return strings.TrimSpace(sb.String())
}
