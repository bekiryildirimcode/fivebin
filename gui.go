package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GUI manages the desktop application interface
type GUI struct {
	app        fyne.App
	window     fyne.Window
	dictionary *Dictionary
	meanings   *UserMeanings

	// UI components
	wordLabel     *widget.RichText
	typeLabel     *widget.Label
	levelLabel    *widget.Label
	phoneticLabel *widget.Label
	examplesList  *widget.List
	examplesData  []string
	meaningEntry  *widget.Entry
	saveButton    *widget.Button
	newWordButton *widget.Button
	searchEntry   *widget.Entry
	searchButton  *widget.Button
	currentWord   string // Track current word for saving meanings
	statusLabel   *widget.Label
}

// NewGUI creates and initializes the GUI
func NewGUI(dict *Dictionary, meanings *UserMeanings) *GUI {
	g := &GUI{
		app:        app.NewWithID("com.vocabulary.learner"),
		dictionary: dict,
		meanings:   meanings,
	}

	g.window = g.app.NewWindow("Vocabulary Learner")
	g.window.Resize(fyne.NewSize(1000, 800))
	g.window.CenterOnScreen()

	// Set application icon (will be created)
	g.setIcon()

	g.buildUI()

	return g
}

// setIcon sets the application icon
func (g *GUI) setIcon() {
	icon := getAppIcon()
	g.window.SetIcon(icon)
	g.app.SetIcon(icon)
}

// buildUI constructs the user interface with modern styling
func (g *GUI) buildUI() {
	// Search bar at the top
	g.searchEntry = widget.NewEntry()
	g.searchEntry.SetPlaceHolder("Search for a word...")
	g.searchEntry.OnSubmitted = func(text string) {
		g.onSearchClick()
	}

	g.searchButton = widget.NewButtonWithIcon("Search", theme.SearchIcon(), g.onSearchClick)
	g.searchButton.Importance = widget.HighImportance

	// Make search entry expand to fill available width
	searchContainer := container.NewBorder(
		nil, nil,
		nil,            // Left: empty
		g.searchButton, // Right: button
		g.searchEntry,  // Center: entry (takes full width)
	)

	searchCard := container.NewBorder(
		widget.NewLabelWithStyle("🔍 Search Word", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		searchContainer,
	)
	searchCard = container.NewPadded(container.NewBorder(nil, nil, nil, nil, searchCard))

	// Large, centered word display
	g.wordLabel = widget.NewRichText()
	g.wordLabel.Wrapping = fyne.TextWrapOff

	// Word info labels with modern styling
	g.typeLabel = widget.NewLabel("")
	g.typeLabel.Alignment = fyne.TextAlignCenter

	g.levelLabel = widget.NewLabel("")
	g.levelLabel.Alignment = fyne.TextAlignCenter

	g.phoneticLabel = widget.NewLabel("")
	g.phoneticLabel.Alignment = fyne.TextAlignCenter

	// Status label for feedback
	g.statusLabel = widget.NewLabel("")
	g.statusLabel.Alignment = fyne.TextAlignCenter
	g.statusLabel.Wrapping = fyne.TextWrapWord

	// Word display card
	wordCenter := container.NewCenter(g.wordLabel)
	wordInfo := container.NewVBox(
		wordCenter,
		g.typeLabel,
		g.levelLabel,
		g.phoneticLabel,
		g.statusLabel,
	)
	wordCard := container.NewPadded(container.NewCenter(wordInfo))

	// Examples list with modern styling
	g.examplesList = widget.NewList(
		func() int {
			return len(g.examplesData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(g.examplesData) {
				label.SetText(g.examplesData[id])
			}
		},
	)

	examplesCard := container.NewBorder(
		widget.NewLabelWithStyle("📚 Examples", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		g.examplesList,
	)
	examplesCard = container.NewPadded(examplesCard)

	// User meaning input field
	g.meaningEntry = widget.NewMultiLineEntry()
	g.meaningEntry.SetPlaceHolder("Enter your meaning, translation, or notes here...")
	g.meaningEntry.Wrapping = fyne.TextWrapWord

	g.saveButton = widget.NewButtonWithIcon("💾 Save Meaning", theme.DocumentSaveIcon(), g.onSaveMeaningClick)
	g.saveButton.Importance = widget.MediumImportance

	meaningCard := container.NewBorder(
		widget.NewLabelWithStyle("✏️ Your Notes", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(g.saveButton),
		nil, nil,
		g.meaningEntry,
	)
	meaningCard = container.NewPadded(meaningCard)

	// Action buttons
	g.newWordButton = widget.NewButtonWithIcon("🎲 New Word", theme.ContentAddIcon(), g.onNewWordClick)
	g.newWordButton.Importance = widget.HighImportance

	buttonContainer := container.NewHBox(
		container.NewBorder(nil, nil, nil, nil, g.newWordButton),
	)
	buttonContainer = container.NewCenter(buttonContainer)

	// Main layout with modern split views
	leftSection := container.NewVSplit(
		wordCard,
		examplesCard,
	)
	leftSection.SetOffset(0.4) // Word takes 40%, examples 60%

	rightSection := meaningCard

	mainContent := container.NewHSplit(leftSection, rightSection)
	mainContent.SetOffset(0.6) // Left 60%, right 40%

	// Final layout
	content := container.NewBorder(
		searchCard,
		buttonContainer,
		nil, nil,
		mainContent,
	)

	// Add padding to entire content
	finalContent := container.NewPadded(content)
	g.window.SetContent(finalContent)
}

// onSearchClick handles the search button click
func (g *GUI) onSearchClick() {
	searchText := g.searchEntry.Text
	if searchText == "" {
		g.statusLabel.SetText("Please enter a word to search")
		return
	}

	g.searchButton.SetText("Searching...")
	g.searchButton.Disable()
	g.statusLabel.SetText("Searching...")

	go func() {
		word, err := g.dictionary.FindWord(searchText)

		if err != nil {
			g.statusLabel.SetText("❌ " + err.Error())
			g.wordLabel.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:  "Word not found",
					Style: widget.RichTextStyle{ColorName: theme.ColorNameError},
				},
			}
			g.wordLabel.Refresh()
			g.typeLabel.SetText("")
			g.levelLabel.SetText("")
			g.phoneticLabel.SetText("")
			g.examplesData = []string{}
			g.examplesList.Refresh()
			g.meaningEntry.SetText("")
		} else {
			g.statusLabel.SetText("✅ Word found!")
			g.updateWordDisplay(word)
		}

		g.searchButton.SetText("Search")
		g.searchButton.Enable()
	}()
}

// onNewWordClick handles the "New Word" button click
func (g *GUI) onNewWordClick() {
	g.newWordButton.SetText("Loading...")
	g.newWordButton.Disable()
	g.statusLabel.SetText("Loading random word...")

	go func() {
		word, err := g.dictionary.GetRandomWord()

		if err != nil {
			g.statusLabel.SetText("❌ " + err.Error())
			g.wordLabel.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:  "Error: " + err.Error(),
					Style: widget.RichTextStyle{ColorName: theme.ColorNameError},
				},
			}
			g.wordLabel.Refresh()
			g.typeLabel.SetText("")
			g.levelLabel.SetText("")
			g.phoneticLabel.SetText("")
			g.examplesData = []string{}
			g.examplesList.Refresh()
		} else {
			g.statusLabel.SetText("")
			g.updateWordDisplay(word)
		}

		g.newWordButton.SetText("🎲 New Word")
		g.newWordButton.Enable()
	}()
}

// updateWordDisplay updates all UI elements with the new word data
func (g *GUI) updateWordDisplay(word *Root) {
	if word == nil {
		return
	}

	g.currentWord = word.Value.Word

	// Update word label with large, bold text
	g.wordLabel.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  word.Value.Word,
			Style: widget.RichTextStyle{SizeName: theme.SizeNameHeadingText, TextStyle: fyne.TextStyle{Bold: true}},
		},
	}
	g.wordLabel.Refresh()

	// Format type with icon
	typeText := "📝 " + word.Value.Type
	g.typeLabel.SetText(typeText)

	// Format level with icon
	levelText := "📊 Level: " + word.Value.Level
	g.levelLabel.SetText(levelText)

	// Format phonetics (US / UK)
	phoneticText := "🔊 US: " + word.Value.Phonetics.US
	if word.Value.Phonetics.UK != "" {
		phoneticText += " | UK: " + word.Value.Phonetics.UK
	}
	g.phoneticLabel.SetText(phoneticText)

	// Update examples
	g.examplesData = word.Value.Examples
	g.examplesList.Refresh()

	// Load saved meaning if it exists
	savedMeaning := g.meanings.Get(word.Value.Word)
	g.meaningEntry.SetText(savedMeaning)
}

// onSaveMeaningClick handles the "Save Meaning" button click
func (g *GUI) onSaveMeaningClick() {
	if g.currentWord == "" {
		g.statusLabel.SetText("No word selected")
		return
	}

	meaning := g.meaningEntry.Text
	if err := g.meanings.Set(g.currentWord, meaning); err != nil {
		g.statusLabel.SetText("❌ Error saving: " + err.Error())
		return
	}

	g.statusLabel.SetText("✅ Meaning saved!")
	g.saveButton.SetText("💾 Saved!")
	g.saveButton.Disable()

	go func() {
		time.Sleep(2 * time.Second)
		g.saveButton.SetText("💾 Save Meaning")
		g.saveButton.Enable()
		if g.currentWord != "" {
			g.statusLabel.SetText("")
		}
	}()
}

// ShowAndRun displays the window and starts the application
func (g *GUI) ShowAndRun() {
	// Load an initial word
	g.onNewWordClick()

	g.window.ShowAndRun()
}
