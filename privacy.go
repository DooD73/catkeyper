package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const privacyTitle = "Your keystrokes stay private"

const privacyStatement = "CatKeyper works locally on your Mac. It checks keyboard events only while locked, so it can block key presses and recognize Shift + CAT. It does not record, store, analyze, or transmit your keystrokes or any other data."

func newPrivacyWindow(catApp fyne.App) fyne.Window {
	window := catApp.NewWindow("CatKeyper — Privacy & Help")
	window.Resize(fyne.NewSize(440, 560))
	window.SetFixedSize(true)
	window.SetCloseIntercept(window.Hide)

	title := canvas.NewText("Privacy", color.NRGBA{R: 86, G: 48, B: 22, A: 255})
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 30
	title.TextStyle = fyne.TextStyle{Bold: true}

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
		title,
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
	window.SetContent(container.NewStack(background, container.NewPadded(container.NewVScroll(panel))))
	return window
}

func newPrivacyMenuItem(show func()) *fyne.MenuItem {
	return fyne.NewMenuItem("Privacy & Help…", show)
}
