package nanogui

import (
	"app/src/external/nanovgo"
)

type TabButton struct {
	Header *TabHeader
	Label string
	w, h int
}

func NewTabButton(header *TabHeader, label string) *TabButton {
	button := &TabButton{
		Label:label,
		Header:header,
	}

	return button
}

func (tb *TabButton) SetSize(w, h int){
	tb.w,tb.h = w, h
}

func (tb *TabButton) PreferredSize(ctx *nanovgo.Context) (int, int){
	labelWidth, bounds := ctx.TextBounds(0,0, tb.Label)
	buttonWidth := int(labelWidth) + 2 * tb.Header.Theme().TabButtonHorizontalPadding
	buttonHeight := int(bounds[3]) - int(bounds[1]) + 2 * tb.Header.Theme().TabButtonVerticalPadding

	return buttonWidth, buttonHeight
}

func (tb *TabButton) calculateVisibleString(ctx *nanovgo.Context) {
	//TODO
}

func (tb *TabButton) drawAtPosition(ctx *nanovgo.Context, xPos,yPos float32, active bool) {
	w, h := tb.Header.Size()
	width := float32(w)
	height := float32(h)
	theme := tb.Header.Theme()

	ctx.Save()
	ctx.IntersectScissor(xPos, yPos, width + 1, height)

	if !active {
		gradtop := theme.ButtonGradientTopPushed
		gradbot := theme.ButtonGradientBotPushed

		ctx.BeginPath()
		ctx.RoundedRect(xPos +1, yPos + 1, width -1, height +1, float32(theme.ButtonCornerRadius))
		backgroundColor := nanovgo.LinearGradient(xPos,yPos,xPos,yPos+height,gradtop,gradbot)

		ctx.SetFillPaint(backgroundColor)
		ctx.Fill()

		ctx.BeginPath()
		ctx.RoundedRect(xPos + 0.5, yPos + 1.5, width, height, float32(theme.ButtonCornerRadius))
		ctx.SetStrokeColor(theme.BorderDark)
		ctx.Stroke()

	} else {
		ctx.BeginPath()
		ctx.SetStrokeWidth(1.0)
		ctx.RoundedRect(xPos + 0.5, yPos + 1.5, width, height+1, float32(theme.ButtonCornerRadius))
		ctx.SetStrokeColor(theme.BorderLight)
		ctx.Stroke()

		ctx.BeginPath()
		ctx.RoundedRect(xPos + 0.5, yPos + 0.5, width, height+1, float32(theme.ButtonCornerRadius))
		ctx.SetStrokeColor(theme.BorderDark)
		ctx.Stroke()
	}
	ctx.ResetScissor()  //TODO: Check this
	ctx.Restore()

	textX := xPos + float32(theme.TabButtonHorizontalPadding)
	textY := yPos + float32(theme.TabButtonVerticalPadding)
	textColor := theme.TextColor
	ctx.BeginPath()
	ctx.SetFillColor(textColor)
	ctx.Text(textX, textY,tb.Label )  //TODO: dots handling
}

//TODO: No longer used?
//func (tb *TabButton) drawInactiveBorderAt(ctx *nanovgo.Context, xPos,yPos, offset float32, color nanovgo.Color) {
//
//}
//
//func (tb *TabButton) drawActiveBorderAt(ctx *nanovgo.Context, xPos,yPos, offset float32, color nanovgo.Color) {
//
//}

