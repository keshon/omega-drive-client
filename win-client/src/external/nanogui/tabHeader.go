package nanogui

import (
	"app/src/external/nanovgo"
)

type TabHeader struct {
	WidgetImplement
	tabButtons               []*TabButton
	activeTab                int
	callback                 func(index int)
	visibleStart, visibleEnd int
	overflowing bool
}

func NewTabHeader(parent Widget) *TabHeader {
	header := &TabHeader{}

	InitWidget(header, parent)
	return header
}

func (t *TabHeader) SetActiveTab(index int) {
	if index < 0 || index >= len(t.tabButtons) {
		return
	}
	t.activeTab = index

	if t.callback == nil {
		return
	}
	t.callback(index)
}

func (t *TabHeader) ActiveTab() int {
	return t.activeTab
}

func (t *TabHeader) isVisibleTab(index int) bool {
	return index >= t.visibleStart && index < t.visibleEnd
}

func (t *TabHeader) AddTab(index int, label string) {
	if index > len(t.tabButtons) {
		index = len(t.tabButtons)
	}

	tab := NewTabButton(t, label)
	t.tabButtons = append(t.tabButtons, nil)
	copy(t.tabButtons[index+1:], t.tabButtons[index:])
	t.tabButtons[index] = tab

	t.SetActiveTab(index)
}

func (t *TabHeader) RemoveTab(index int) {
	copy(t.tabButtons[index:], t.tabButtons[index+1:])
	t.tabButtons[len(t.tabButtons)-1] = nil
	t.tabButtons = t.tabButtons[:len(t.tabButtons)-1]
}

func (t *TabHeader) tabIndex(label string) (index int, ok bool) {
	for i := range t.tabButtons {
		if t.tabButtons[i].Label == label {
			return i, true
		}
	}
	return -1, false
}

func (t *TabHeader) TabLabelAt(index int) string {
	if index < 0 || index >= len(t.tabButtons) {
		return ""
	}
	return t.tabButtons[index].Label
}

func (t *TabHeader) visibleButtonWidth() (width int) {
	if t.visibleStart == t.visibleEnd {
		return 0
	}
	width = t.x + t.Theme().TabControlWidth
	for _, b := range t.tabButtons {
		width += b.w
	}
	 return width
}

func (t *TabHeader) activeButtonWidth() (width int) {
	if t.visibleStart == t.visibleEnd || t.activeTab < t.visibleStart || t.activeTab >= t.visibleEnd {
		return 0
	}
	width = t.x + t.Theme().TabControlWidth
	for _, b := range t.tabButtons {
		width += b.w
	}
	return width
}

func (t *TabHeader) ensureTabVisible() {
//TODO
}

func (t *TabHeader) calculateVisibleEnd() {
	curPos := t.Theme().TabControlWidth
	lastPos := t.w - t.Theme().TabControlWidth

	for i, b :=range t.tabButtons {
		curPos += b.w
		if curPos > lastPos {
			t.visibleEnd = i
			return
		}
	}
}

func (t *TabHeader) OnPerformLayout(self Widget, ctx *nanovgo.Context) {
	t.WidgetImplement.OnPerformLayout(self, ctx)

	for _, b :=range t.tabButtons {
		prefW, _ := b.PreferredSize(ctx)
		prefW = clampI(t.Theme().TabMinButtonWidth, t.Theme().TabMaxButtonWidth, prefW)
		b.w = prefW
		b.calculateVisibleString(ctx)
	}

	t.calculateVisibleEnd()

	t.overflowing = t.visibleStart !=0 || t.visibleEnd != len(t.tabButtons) -1
}

func (t *TabHeader) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	ctx.SetFontFace("sans")
	ctx.SetFontSize(20)
	ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
	w := 2*t.theme.TabControlWidth
	h := 0
	for _, b :=range t.tabButtons {
		prefW, prefH := b.PreferredSize(ctx)
		prefW = clampI(t.Theme().TabMinButtonWidth, t.Theme().TabMaxButtonWidth, prefW)
		w += prefW
		h = maxI(h, prefH)
	}

	return w, h
}

func (t *TabHeader) Draw(self Widget, ctx *nanovgo.Context) {
	t.WidgetImplement.Draw(self, ctx)

	if t.overflowing {
		t.drawControls(ctx)
	}

}

func (t *TabHeader) drawControls(ctx *nanovgo.Context) {
	var arrowColor nanovgo.Color

	//Left Button
	ctx.BeginPath()
	iconLeft := IconLeftBold
	fontSize := t.theme.ButtonFontSize   //TODO: Consider to handle specific FontSize
	ih := float32(fontSize) * 1.5
	ctx.SetFontSize(ih)
	ctx.SetFontFace(t.theme.FontIcons)

	if t.visibleStart != 0 {
		arrowColor = t.theme.TextColor
	} else {
		arrowColor = t.theme.ButtonGradientBotPushed
	}
	ctx.SetFillColor(arrowColor)
	ctx.SetTextAlign(nanovgo.AlignMiddle | nanovgo.AlignLeft)

	xScale := float32(0.2)
	iconPosX := float32(t.x) + xScale*float32(t.theme.TabControlWidth)
	iconPosY := float32(t.y) + float32(t.h)/2.0 + 1
	ctx.TextRune(iconPosX, iconPosY, []rune{rune(iconLeft)})

	// Right Button
	if t.visibleEnd != len(t.tabButtons)-1{
		arrowColor = t.theme.TextColor
	} else {
		arrowColor = t.theme.ButtonGradientBotPushed
	}
	ctx.SetFillColor(arrowColor)

	iconRight := IconRightBold
	iconPosX = float32(t.x + t.w - t.theme.TabControlWidth) - xScale*float32(t.theme.TabControlWidth)
	ctx.TextRune(iconPosX, iconPosY, []rune{rune(iconRight)})



}