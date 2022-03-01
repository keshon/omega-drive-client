package ui

import (
	"app/src/external/nanogui"
	"app/src/external/nanovgo"
)

func GridPanel(screen *nanogui.Screen, x int, y int) {
	window := nanogui.NewWindow(screen, "Grid of small widgets")
	window.SetPosition(x, y)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	{
		nanogui.NewLabel(window, "Regular text :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "Text")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Floating point :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50.0")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("GiB")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	}
	{
		nanogui.NewLabel(window, "Positive integer :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	{
		nanogui.NewLabel(window, "Float box :").SetFont("sans-bold")
		floatBox := nanogui.NewFloatBox(window, 10.0)
		floatBox.SetEditable(true)
		floatBox.SetFixedSize(100, 20)
		floatBox.SetUnits("GiB")
		floatBox.SetDefaultValue(0.0)
		floatBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Int box :").SetFont("sans-bold")
		intBox := nanogui.NewIntBox(window, true, 50)
		intBox.SetEditable(true)
		intBox.SetFixedSize(100, 20)
		intBox.SetUnits("MHz")
		intBox.SetDefaultValue(0)
		intBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Checkbox :").SetFont("sans-bold")
		checkbox := nanogui.NewCheckBox(window, "Check me")
		checkbox.SetFontSize(16)
		checkbox.SetChecked(true)
	}
	{
		nanogui.NewLabel(window, "Combobox :").SetFont("sans-bold")
		combobox := nanogui.NewComboBox(window, []string{"Item 1", "Item 2", "Item 3"})
		combobox.SetFontSize(16)
		combobox.SetFixedSize(100, 20)
	}
	{
		nanogui.NewLabel(window, "Color button :").SetFont("sans-bold")

		popupButton := nanogui.NewPopupButton(window, "")
		popupButton.SetBackgroundColor(nanovgo.RGBA(255, 120, 0, 255))
		popupButton.SetFontSize(16)
		popupButton.SetFixedSize(100, 20)
		popup := popupButton.Popup()
		popup.SetLayout(nanogui.NewGroupLayout())

		colorWheel := nanogui.NewColorWheel(popup)
		colorWheel.SetColor(popupButton.BackgroundColor())

		colorButton := nanogui.NewButton(popup, "Pick")
		colorButton.SetFixedSize(100, 25)
		colorButton.SetBackgroundColor(colorWheel.Color())

		colorWheel.SetCallback(func(color nanovgo.Color) {
			colorButton.SetBackgroundColor(color)
		})

		colorButton.SetChangeCallback(func(pushed bool) {
			if pushed {
				popupButton.SetBackgroundColor(colorButton.BackgroundColor())
				popupButton.SetPushed(false)
			}
		})
	}
}
