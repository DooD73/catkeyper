package main

/*
#cgo darwin LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdint.h>

extern int goKeyboardDecision(uint16_t keycode, uint64_t flags, uint32_t eventType);
extern void goEventTapReenabled(uint32_t eventType);

static CFMachPortRef keyboardEventTap = NULL;

static CGEventRef keyboardTapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
	if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
		if (keyboardEventTap != NULL) {
			CGEventTapEnable(keyboardEventTap, true);
			goEventTapReenabled((uint32_t)type);
		}
		return event;
	}

	if (type != kCGEventKeyDown && type != kCGEventKeyUp && type != kCGEventFlagsChanged) {
		return event;
	}

	uint16_t keycode = (uint16_t)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
	uint64_t flags = (uint64_t)CGEventGetFlags(event);

	// Returning NULL from a Quartz event tap callback consumes the event. This
	// is the low-level macOS suppression point: while Cat Mode is locked, normal
	// key presses never continue to the active application.
	if (goKeyboardDecision(keycode, flags, (uint32_t)type) != 0) {
		return NULL;
	}

	return event;
}

static int startKeyboardTap(void) {
	CGEventMask mask = CGEventMaskBit(kCGEventKeyDown) |
	                   CGEventMaskBit(kCGEventKeyUp) |
	                   CGEventMaskBit(kCGEventFlagsChanged);

	keyboardEventTap = CGEventTapCreate(kCGSessionEventTap,
	                                    kCGHeadInsertEventTap,
	                                    kCGEventTapOptionDefault,
	                                    mask,
	                                    keyboardTapCallback,
	                                    NULL);
	if (keyboardEventTap == NULL) {
		return 0;
	}

	CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, keyboardEventTap, 0);
	CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
	CGEventTapEnable(keyboardEventTap, true);
	CFRunLoopRun();

	CFRunLoopRemoveSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
	CFRelease(source);
	CFRelease(keyboardEventTap);
	keyboardEventTap = NULL;
	return 1;
}

static int accessibilityTrusted(void) {
	return AXIsProcessTrusted() ? 1 : 0;
}

static int accessibilityTrustedWithPrompt(void) {
	const void *keys[] = { kAXTrustedCheckOptionPrompt };
	const void *values[] = { kCFBooleanTrue };
	CFDictionaryRef options = CFDictionaryCreate(kCFAllocatorDefault,
	                                             keys,
	                                             values,
	                                             1,
	                                             &kCFCopyStringDictionaryKeyCallBacks,
	                                             &kCFTypeDictionaryValueCallBacks);
	Boolean trusted = AXIsProcessTrustedWithOptions(options);
	CFRelease(options);
	return trusted ? 1 : 0;
}
*/
import "C"

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"catkeyper/pkg/lockmanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/oksvg"
	"github.com/srwiley/rasterx"
)

//go:embed assets/openmoji-cat-face.svg
var openMojiCatSVG []byte

//go:embed assets/openmoji-cat-face-outline.svg
var openMojiCatOutlineSVG []byte

var openMojiCat = fyne.NewStaticResource("openmoji-cat-face.svg", openMojiCatSVG)

const appName = "CatKeyper"

var (
	trayIconActive   = newTrayIconResource("catkeyper-tray-active.png", 0xff)
	trayIconInactive = newTrayIconResource("catkeyper-tray-inactive.png", 0x8c)
)

var (
	appVersion = "1.0.0"

	keyboard            = lockmanager.New()
	unlockNotifications = make(chan struct{}, 1)
	appLogger           = newLogger()
)

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	if debugLoggingEnabled() {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}))
}

func debugLoggingEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("CATKEYPER_DEBUG")))
	return v == "1" || v == "true" || v == "yes" || v == "debug"
}

func newTrayIconResource(name string, alpha uint8) fyne.Resource {
	const size = 72

	icon, err := oksvg.ReadIconStream(bytes.NewReader(openMojiCatOutlineSVG))
	if err != nil {
		return theme.NewThemedResource(fyne.NewStaticResource("openmoji-cat-face-outline.svg", openMojiCatOutlineSVG))
	}

	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	icon.SetTarget(0, 0, size, size)
	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	dasher := rasterx.NewDasher(size, size, scanner)
	icon.Draw(dasher, 1)

	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i] = 0
		img.Pix[i+1] = 0
		img.Pix[i+2] = 0
		img.Pix[i+3] = uint8(uint16(img.Pix[i+3]) * uint16(alpha) / 0xff)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return theme.NewThemedResource(fyne.NewStaticResource("openmoji-cat-face-outline.svg", openMojiCatOutlineSVG))
	}
	return theme.NewThemedResource(fyne.NewStaticResource(name, buf.Bytes()))
}

//export goKeyboardDecision
func goKeyboardDecision(keycode C.uint16_t, flags C.uint64_t, eventType C.uint32_t) C.int {
	suppress, unlocked := keyboard.KeyboardDecision(uint16(keycode), uint64(flags), uint32(eventType))
	if unlocked {
		count := keyboard.SuppressedCount()
		appLogger.Info("global unlock sequence accepted",
			"suppressed_events", count,
			"event_type", lockmanager.EventTypeName(uint32(eventType)),
		)
		select {
		case unlockNotifications <- struct{}{}:
		default:
			appLogger.Debug("unlock notification already pending")
		}
	}
	if suppress {
		return 1
	}
	return 0
}

//export goEventTapReenabled
func goEventTapReenabled(eventType C.uint32_t) {
	appLogger.Warn("macOS keyboard event tap was disabled and re-enabled",
		"reason", lockmanager.EventTypeName(uint32(eventType)),
	)
}

func startKeyboardHook(status chan<- string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	appLogger.Info("starting macOS keyboard hook",
		"tap", "cg_session_event_tap",
		"thread_locked", true,
	)

	if C.accessibilityTrusted() == 0 {
		appLogger.Warn("accessibility permission is not granted; requesting user approval")
		C.accessibilityTrustedWithPrompt()
		status <- "Accessibility permission is required. Enable it in System Settings, then restart CatKeyper."
		appLogger.Warn("keyboard hook startup paused until Accessibility permission is granted")
		return
	}

	appLogger.Info("accessibility permission verified")
	status <- "Keyboard guard ready."
	if C.startKeyboardTap() == 0 {
		appLogger.Error("failed to create macOS keyboard event tap")
		status <- "Could not create the macOS keyboard event tap. Check Accessibility/Input Monitoring permissions."
		return
	}
	appLogger.Info("keyboard event tap run loop exited")
}

type catTheme struct{}

func (catTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 255, G: 246, B: 226, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 239, G: 131, B: 47, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 224, G: 198, B: 166, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 72, G: 45, B: 24, A: 255}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 255, G: 228, B: 184, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 215, B: 148, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 255, G: 252, B: 244, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 139, G: 96, B: 54, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 205, G: 95, B: 28, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 232, G: 103, B: 31, A: 255}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 188, G: 125, B: 70, A: 255}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 255, G: 190, B: 103, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 92, G: 54, B: 25, A: 80}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (catTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (catTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (catTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 12
	case theme.SizeNameText:
		return 15
	case theme.SizeNameHeadingText:
		return 25
	case theme.SizeNameInlineIcon:
		return 20
	}
	return theme.DefaultTheme().Size(name)
}

type catScene struct {
	root       *fyne.Container
	cat        *canvas.Image
	badge      *canvas.Text
	zzz        *canvas.Text
	lockBand   *canvas.Rectangle
	lockBandHi *canvas.Rectangle
	locked     bool
	frame      int
	frameLock  sync.Mutex
}

func newCatScene() *catScene {
	c := &catScene{
		cat:        canvas.NewImageFromResource(openMojiCat),
		badge:      canvas.NewText("READY", color.NRGBA{R: 255, G: 246, B: 226, A: 255}),
		zzz:        canvas.NewText("Zzz", color.NRGBA{R: 101, G: 67, B: 38, A: 255}),
		lockBand:   canvas.NewRectangle(color.NRGBA{R: 101, G: 67, B: 38, A: 210}),
		lockBandHi: canvas.NewRectangle(color.NRGBA{R: 255, G: 205, B: 123, A: 230}),
	}

	c.cat.FillMode = canvas.ImageFillContain
	for _, txt := range []*canvas.Text{c.badge, c.zzz} {
		txt.Alignment = fyne.TextAlignCenter
		txt.TextStyle = fyne.TextStyle{Bold: true}
	}
	c.badge.TextSize = 18
	c.zzz.TextSize = 28

	bg := canvas.NewRectangle(color.NRGBA{R: 255, G: 235, B: 196, A: 255})
	accent := canvas.NewRectangle(color.NRGBA{R: 255, G: 128, B: 0, A: 255})
	soft := canvas.NewRectangle(color.NRGBA{R: 255, G: 246, B: 226, A: 255})
	checkerA := canvas.NewRectangle(color.NRGBA{R: 255, G: 221, B: 159, A: 255})
	checkerB := canvas.NewRectangle(color.NRGBA{R: 255, G: 246, B: 226, A: 255})

	c.root = container.NewWithoutLayout(bg, checkerA, checkerB, accent, soft, c.cat, c.lockBand, c.lockBandHi, c.badge, c.zzz)
	bg.Resize(fyne.NewSize(280, 280))
	checkerA.Resize(fyne.NewSize(72, 72))
	checkerA.Move(fyne.NewPos(0, 0))
	checkerB.Resize(fyne.NewSize(72, 72))
	checkerB.Move(fyne.NewPos(208, 208))
	accent.Resize(fyne.NewSize(280, 58))
	accent.Move(fyne.NewPos(0, 222))
	soft.Resize(fyne.NewSize(210, 210))
	soft.Move(fyne.NewPos(35, 20))
	c.cat.Resize(fyne.NewSize(178, 178))
	c.cat.Move(fyne.NewPos(51, 42))
	c.lockBand.Resize(fyne.NewSize(190, 30))
	c.lockBand.Move(fyne.NewPos(45, 120))
	c.lockBandHi.Resize(fyne.NewSize(140, 8))
	c.lockBandHi.Move(fyne.NewPos(70, 131))
	c.badge.Resize(fyne.NewSize(180, 32))
	c.badge.Move(fyne.NewPos(50, 236))
	c.zzz.Resize(fyne.NewSize(88, 36))
	c.zzz.Move(fyne.NewPos(182, 26))
	c.root.Resize(fyne.NewSize(280, 280))

	c.setLocked(false)
	return c
}

func (c *catScene) setLocked(isLocked bool) {
	c.frameLock.Lock()
	defer c.frameLock.Unlock()

	c.locked = isLocked
	if isLocked {
		c.badge.Text = "KEYS BLOCKED"
		c.cat.Translucency = 0.12
		c.lockBand.Show()
		c.lockBandHi.Show()
		c.zzz.Show()
	} else {
		c.badge.Text = "ALL CLEAR"
		c.cat.Translucency = 0
		c.lockBand.Hide()
		c.lockBandHi.Hide()
		c.zzz.Hide()
	}
	canvas.Refresh(c.root)
}

func (c *catScene) animate() {
	ticker := time.NewTicker(450 * time.Millisecond)
	for range ticker.C {
		c.frameLock.Lock()
		c.frame++
		frame := c.frame
		isLocked := c.locked
		c.frameLock.Unlock()

		fyne.Do(func() {
			if isLocked {
				c.zzz.Text = []string{"Z", "Zz", "Zzz", "zZz"}[frame%4]
				c.zzz.Move(fyne.NewPos(182+float32((frame%2)*5), 26-float32((frame%3)*3)))
			} else if frame%2 == 0 {
				c.cat.Move(fyne.NewPos(51, 39))
			} else {
				c.cat.Move(fyne.NewPos(51, 43))
			}
			canvas.Refresh(c.root)
		})
	}
}

var supportLinkColor = color.NRGBA{R: 139, G: 90, B: 43, A: 210}

type supportLink struct {
	widget.BaseWidget
	hovered bool
}

func newSupportLink() *supportLink {
	l := &supportLink{}
	l.ExtendBaseWidget(l)
	return l
}

func (l *supportLink) MouseIn(*desktop.MouseEvent)    { l.hovered = true; l.Refresh() }
func (l *supportLink) MouseMoved(*desktop.MouseEvent) {}
func (l *supportLink) MouseOut()                      { l.hovered = false; l.Refresh() }

func (l *supportLink) Tapped(*fyne.PointEvent) {
	u, _ := url.Parse("mailto:catkeyper.support@gmail.com")
	fyne.CurrentApp().OpenURL(u)
}

type supportLinkRenderer struct {
	link      *supportLink
	text      *canvas.Text
	underline *canvas.Rectangle
}

func (r *supportLinkRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
	r.text.Move(fyne.NewPos(0, 0))
	r.underline.Resize(fyne.NewSize(size.Width, 1))
	r.underline.Move(fyne.NewPos(0, size.Height-1))
}

func (r *supportLinkRenderer) MinSize() fyne.Size { return r.text.MinSize() }

func (r *supportLinkRenderer) Refresh() {
	if r.link.hovered {
		r.underline.Show()
	} else {
		r.underline.Hide()
	}
	canvas.Refresh(r.text)
	canvas.Refresh(r.underline)
}

func (r *supportLinkRenderer) Destroy() {}

func (r *supportLinkRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text, r.underline}
}

func (l *supportLink) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText("Contact Support", supportLinkColor)
	text.TextSize = 12
	text.Alignment = fyne.TextAlignCenter
	underline := canvas.NewRectangle(supportLinkColor)
	underline.Hide()
	return &supportLinkRenderer{link: l, text: text, underline: underline}
}

func main() {
	appLogger.Info("application starting",
		"name", appName,
		"version", appVersion,
		"goos", runtime.GOOS,
		"goarch", runtime.GOARCH,
		"debug_logs", debugLoggingEnabled(),
	)

	statusUpdates := make(chan string, 4)
	go startKeyboardHook(statusUpdates)

	catApp := app.New()
	catApp.Settings().SetTheme(catTheme{})
	catApp.SetIcon(trayIconInactive)

	win := catApp.NewWindow(appName)
	win.Resize(fyne.NewSize(440, 560))
	win.SetFixedSize(true)
	appLogger.Debug("main window configured", "width", 440, "height", 560, "fixed_size", true)

	scene := newCatScene()
	go scene.animate()

	title := canvas.NewText("CatKeyper", color.NRGBA{R: 86, G: 48, B: 22, A: 255})
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 30
	title.TextStyle = fyne.TextStyle{Bold: true}

	stateText := canvas.NewText("UNLOCKED", color.NRGBA{R: 48, G: 121, B: 82, A: 255})
	stateText.Alignment = fyne.TextAlignCenter
	stateText.TextSize = 24
	stateText.TextStyle = fyne.TextStyle{Bold: true}
	stateBg := canvas.NewRectangle(color.NRGBA{R: 236, G: 255, B: 242, A: 255})
	statePill := container.NewGridWrap(fyne.NewSize(260, 44), container.NewStack(stateBg, container.NewCenter(stateText)))

	helpText := widget.NewLabel("Unlock: hold Shift, then press C, A, T.")
	helpText.Wrapping = fyne.TextWrapWord
	helpText.Alignment = fyne.TextAlignCenter

	hookStatus := widget.NewLabel("Starting keyboard guard...")
	hookStatus.Wrapping = fyne.TextWrapWord
	hookStatus.Alignment = fyne.TextAlignCenter

	lockButton := widget.NewButtonWithIcon("Lock Keyboard", theme.VisibilityOffIcon(), nil)
	unlockButton := widget.NewButtonWithIcon("Unlock Keyboard", theme.ConfirmIcon(), nil)

	var (
		trayApp        desktop.App
		trayMenu       *fyne.Menu
		trayLockItem   *fyne.MenuItem
		trayUnlockItem *fyne.MenuItem
	)
	if d, ok := catApp.(desktop.App); ok {
		trayApp = d
		trayLockItem = fyne.NewMenuItem("Lock Keyboard", nil)
		trayLockItem.Icon = theme.VisibilityOffIcon()
		trayUnlockItem = fyne.NewMenuItem("Unlock Keyboard", nil)
		trayUnlockItem.Icon = theme.ConfirmIcon()
		trayMenu = fyne.NewMenu(appName,
			trayLockItem,
			trayUnlockItem,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Show CatKeyper", func() {
				win.Show()
				win.RequestFocus()
			}),
		)
		trayApp.SetSystemTrayMenu(trayMenu)
		trayApp.SetSystemTrayIcon(trayIconInactive)
	}

	setLocked := func(isLocked bool) {
		previous := keyboard.IsLocked()
		keyboard.SetLocked(isLocked)
		if previous != isLocked {
			appLogger.Info("keyboard lock state changed",
				"locked", isLocked,
				"suppressed_events_total", keyboard.SuppressedCount(),
			)
		} else {
			appLogger.Debug("keyboard lock state refreshed", "locked", isLocked)
		}

		scene.setLocked(isLocked)
		if isLocked {
			win.SetTitle(appName + " - Locked")
			stateText.Text = "LOCKED"
			stateText.Color = color.NRGBA{R: 171, G: 57, B: 31, A: 255}
			stateBg.FillColor = color.NRGBA{R: 255, G: 235, B: 224, A: 255}
			lockButton.Disable()
			unlockButton.Enable()
			if trayLockItem != nil {
				trayLockItem.Disabled = true
				trayUnlockItem.Disabled = false
				trayApp.SetSystemTrayIcon(trayIconActive)
			}
		} else {
			win.SetTitle(appName + " - Unlocked")
			stateText.Text = "UNLOCKED"
			stateText.Color = color.NRGBA{R: 48, G: 121, B: 82, A: 255}
			stateBg.FillColor = color.NRGBA{R: 236, G: 255, B: 242, A: 255}
			lockButton.Enable()
			unlockButton.Disable()
			if trayLockItem != nil {
				trayLockItem.Disabled = false
				trayUnlockItem.Disabled = true
				trayApp.SetSystemTrayIcon(trayIconInactive)
			}
		}
		if trayMenu != nil {
			trayMenu.Refresh()
		}
		stateText.Refresh()
		stateBg.Refresh()
	}

	lockButton.OnTapped = func() {
		appLogger.Info("lock requested from UI")
		setLocked(true)
	}
	unlockButton.OnTapped = func() {
		appLogger.Info("unlock requested from UI")
		setLocked(false)
	}
	unlockButton.Disable()
	if trayLockItem != nil {
		trayLockItem.Action = func() {
			appLogger.Info("lock requested from menu bar")
			setLocked(true)
		}
		trayUnlockItem.Action = func() {
			appLogger.Info("unlock requested from menu bar")
			setLocked(false)
		}
	}
	setLocked(false)

	sourceText := widget.NewLabel("Cat asset: OpenMoji, CC BY-SA 4.0")
	sourceText.Alignment = fyne.TextAlignCenter
	sourceText.TextStyle = fyne.TextStyle{Italic: true}

	buttons := container.NewGridWithColumns(2, lockButton, unlockButton)
	art := container.NewGridWrap(fyne.NewSize(280, 280), scene.root)
	panel := container.NewVBox(
		title,
		container.NewCenter(statePill),
		container.NewCenter(art),
		helpText,
		buttons,
		hookStatus,
		sourceText,
		container.NewCenter(newSupportLink()),
	)

	background := canvas.NewRectangle(color.NRGBA{R: 255, G: 246, B: 226, A: 255})
	win.SetContent(container.NewStack(background, container.NewPadded(panel)))

	go func() {
		for {
			select {
			case <-unlockNotifications:
				appLogger.Debug("processing global unlock notification on UI thread")
				fyne.Do(func() { setLocked(false) })
			case msg := <-statusUpdates:
				text := msg
				appLogger.Info("keyboard guard status updated", "status", text)
				fyne.Do(func() {
					hookStatus.SetText(text)
				})
			}
		}
	}()

	win.SetCloseIntercept(func() {
		appLogger.Info("application closing", "suppressed_events_total", keyboard.SuppressedCount())
		keyboard.SetLocked(false)
		win.Close()
	})

	_ = unsafe.Sizeof(C.int(0))
	appLogger.Info("showing main window")
	win.ShowAndRun()
	appLogger.Info("application stopped", "suppressed_events_total", keyboard.SuppressedCount())
}
