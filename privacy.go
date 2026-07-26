package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const privacyTitle = "Your keystrokes stay private"

const privacyStatement = "CatKeyper works locally on your Mac. It checks keyboard events only while locked, so it can block key presses and recognize Shift + CAT. It does not record, store, analyze, or transmit your keystrokes or any other data."

func newPrivacyPage(goBack func()) fyne.CanvasObject {
	title := canvas.NewText("Privacy", color.NRGBA{R: 86, G: 48, B: 22, A: 255})
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 30
	title.TextStyle = fyne.TextStyle{Bold: true}
	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), goBack)
	header := container.NewStack(
		container.NewCenter(title),
		container.NewBorder(nil, nil, back, nil),
	)

	headline := widget.NewLabel(privacyTitle)
	headline.Alignment = fyne.TextAlignCenter
	headline.TextStyle = fyne.TextStyle{Bold: true}

	statement := widget.NewLabel(privacyStatement)
	statement.Wrapping = fyne.TextWrapWord

	permissionTitle := widget.NewLabel("Why Accessibility permission?")
	permissionTitle.TextStyle = fyne.TextStyle{Bold: true}

	permissionText := widget.NewLabel("macOS requires this permission so CatKeyper can block key presses system-wide and recognize Shift + CAT while locked.")
	permissionText.Wrapping = fyne.TextWrapWord

	helpText := widget.NewLabel("Help: Unlock from the menu bar, or hold Shift and press C, A, T in sequence.")
	helpText.Wrapping = fyne.TextWrapWord

	aboutTitle := widget.NewLabel("About CatKeyper")
	aboutTitle.TextStyle = fyne.TextStyle{Bold: true}

	aboutText := widget.NewLabel(fmt.Sprintf("%s %s\nFree and open source under the MIT license.", appName, appVersion))
	aboutText.Alignment = fyne.TextAlignCenter

	panel := container.NewVBox(
		headline,
		container.NewGridWrap(fyne.NewSize(380, 105), statement),
		widget.NewSeparator(),
		permissionTitle,
		permissionText,
		helpText,
		widget.NewSeparator(),
		aboutTitle,
		aboutText,
		container.NewCenter(newSupportLink()),
	)

	background := canvas.NewRectangle(color.NRGBA{R: 255, G: 246, B: 226, A: 255})
	page := container.NewBorder(header, nil, nil, nil, container.NewVScroll(panel))
	return container.NewStack(background, container.NewPadded(page))
}

func newPrivacyMenuItem(show func()) *fyne.MenuItem {
	return fyne.NewMenuItem("Privacy & Help…", show)
}
