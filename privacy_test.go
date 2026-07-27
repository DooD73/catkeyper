package main

import (
	"os"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fynetest "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestPrivacyStatementMeetsAcceptanceCriteria(t *testing.T) {
	required := []string{
		"works locally on your Mac",
		"does not record, store, analyze, or transmit",
		"keystrokes or any other data",
	}

	for _, phrase := range required {
		if !strings.Contains(privacyStatement, phrase) {
			t.Errorf("privacy statement does not contain %q", phrase)
		}
	}
}

func TestPrivacyPageShowsCanonicalStatementAndBackAction(t *testing.T) {
	catApp := fynetest.NewApp()
	t.Cleanup(catApp.Quit)

	backCalled := false
	page := newPrivacyPage(func() { backCalled = true })

	for _, text := range []string{
		privacyStatement,
		"Why Accessibility permission?",
		"About CatKeyper",
	} {
		if !containsLabelText(page, text) {
			t.Fatalf("privacy page does not show %q", text)
		}
	}
	if min := page.MinSize(); min.Width > 440 || min.Height > 560 {
		t.Fatalf("privacy content minimum size %v exceeds its 440x560 window", min)
	}

	back := findButtonByIcon(page, theme.NavigateBackIcon().Name())
	if back == nil {
		t.Fatal("privacy page does not have a Back button")
	}
	if back.Text != "" {
		t.Fatalf("Back button text = %q; want icon-only control", back.Text)
	}
	if back.Importance != widget.LowImportance {
		t.Fatalf("Back button importance = %v; want low importance", back.Importance)
	}
	back.Tapped(&fyne.PointEvent{})
	if !backCalled {
		t.Fatal("Back button did not invoke its action")
	}
}

func TestPrivacyActionsOpenTheSharedPage(t *testing.T) {
	t.Run("footer link", func(t *testing.T) {
		called := false
		link := newFooterLink("Privacy", func() { called = true })

		if link.label != "Privacy" {
			t.Fatalf("footer label = %q", link.label)
		}
		link.Tapped(&fyne.PointEvent{})
		if !called {
			t.Fatal("privacy footer link did not invoke its action")
		}
	})

	t.Run("menu item", func(t *testing.T) {
		called := false
		item := newPrivacyMenuItem(func() { called = true })

		if item.Label != "Privacy & Help…" {
			t.Fatalf("menu label = %q", item.Label)
		}
		item.Action()
		if !called {
			t.Fatal("privacy menu item did not invoke its action")
		}
	})
}

func TestREADMEUsesCanonicalPrivacyStatement(t *testing.T) {
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readme), privacyStatement) {
		t.Fatal("README does not contain the canonical privacy statement")
	}
}

func containsLabelText(object fyne.CanvasObject, text string) bool {
	switch object := object.(type) {
	case *widget.Label:
		return object.Text == text
	case *widget.Card:
		return containsLabelText(object.Content, text)
	case *container.Scroll:
		return containsLabelText(object.Content, text)
	case *fyne.Container:
		for _, child := range object.Objects {
			if containsLabelText(child, text) {
				return true
			}
		}
	}
	return false
}

func findButtonByIcon(object fyne.CanvasObject, iconName string) *widget.Button {
	switch object := object.(type) {
	case *widget.Button:
		if object.Icon != nil && object.Icon.Name() == iconName {
			return object
		}
	case *widget.Card:
		return findButtonByIcon(object.Content, iconName)
	case *container.Scroll:
		return findButtonByIcon(object.Content, iconName)
	case *fyne.Container:
		for _, child := range object.Objects {
			if button := findButtonByIcon(child, iconName); button != nil {
				return button
			}
		}
	}
	return nil
}
