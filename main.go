package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Constants for better maintainability
const (
	MinSymbolsPerCard = 2
	MaxSymbolsPerCard = 15
	WindowWidth       = 800
	WindowHeight      = 700
	CardsListHeight   = 300
)

// Structure for tracking card states
type CardItem struct {
	Card      []int
	Processed bool
}

// UI components grouped together
type AppUI struct {
	window              fyne.Window
	symbolsPerCardEntry *widget.Entry
	resultLabel         *widget.Label
	successLabel        *widget.Label
	errorLabel          *widget.Label
	cardsList           *widget.List
	cardItems           []CardItem
}

// Business logic separated from UI
type DobbleGenerator struct {
	symbolsPerCard int
	cards          [][]int
}

func NewDobbleGenerator() *DobbleGenerator {
	return &DobbleGenerator{}
}

func (dg *DobbleGenerator) Generate(symbolsPerCard int) [][]int {
	dg.symbolsPerCard = symbolsPerCard
	dg.cards = generateDobbleByY(symbolsPerCard)
	return dg.cards
}

func (dg *DobbleGenerator) GetTotalSymbols() int {
	n := dg.symbolsPerCard - 1
	return n*n + n + 1
}

// Input validation methods
func validateInput(input string) (int, error) {
	if input == "" {
		return 0, fmt.Errorf("please enter the number of images per card")
	}
	
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("please enter a valid number")
	}
	
	if value <= 1 {
		return 0, fmt.Errorf("the number of images per card must be greater than 1")
	}
	
	if value > MaxSymbolsPerCard {
		return 0, fmt.Errorf("the number of images per card is limited to %d for performance", MaxSymbolsPerCard)
	}
	
	return value, nil
}

// UI helper methods
func (ui *AppUI) clearAllLabels() {
	ui.resultLabel.SetText("")
	ui.successLabel.SetText("")
	ui.errorLabel.SetText("")
}

func (ui *AppUI) showError(message string) {
	ui.clearAllLabels()
	ui.errorLabel.SetText(message)
}

func (ui *AppUI) showSuccess(message string) {
	ui.clearAllLabels()
	ui.successLabel.SetText(message)
}

func generateDobbleByY(y int) [][]int {
	n := y - 1
	cards := [][]int{}

	// 1. First card: [0..n]
	firstCard := []int{}
	for i := 0; i <= n; i++ {
		firstCard = append(firstCard, i+1) // Convert to 1-indexing
	}
	cards = append(cards, firstCard)

	// 2. Next n cards: each with 0 and n others
	for i := 0; i < n; i++ {
		card := []int{1} // 0 becomes 1
		for j := 0; j < n; j++ {
			symbol := (n + 1) + i*n + j + 1 // Fixed formula for 1-indexing
			card = append(card, symbol)
		}
		cards = append(cards, card)
	}

	// 3. Remaining n^2 cards
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			card := []int{i + 2} // i+1 becomes i+2 for 1-indexing
			for k := 0; k < n; k++ {
				symbol := (n + 1) + k*n + ((i*k + j) % n) + 1 // Fixed formula for 1-indexing
				card = append(card, symbol)
			}
			cards = append(cards, card)
		}
	}

	return cards
}

func countCommonSymbols(card1, card2 []int) int {
	count := 0
	for _, s1 := range card1 {
		for _, s2 := range card2 {
			if s1 == s2 {
				count++
			}
		}
	}
	return count
}

func filterValidCards(cards [][]int) [][]int {
	if len(cards) == 0 {
		return cards
	}
	
	// Use greedy approach to find maximum compatible set
	validCards := [][]int{}
	used := make([]bool, len(cards))
	
	// Add first card as base
	validCards = append(validCards, cards[0])
	used[0] = true
	
	// Try to add as many cards as possible
	for {
		added := false
		for i := 1; i < len(cards); i++ {
			if used[i] {
				continue
			}
			
			candidate := cards[i]
			isValid := true
			
			// Check that candidate has exactly 1 common symbol with each already added card
			for _, validCard := range validCards {
				if countCommonSymbols(candidate, validCard) != 1 {
					isValid = false
					break
				}
			}
			
			if isValid {
				validCards = append(validCards, candidate)
				used[i] = true
				added = true
				break // Start over to potentially find better combinations
			}
		}
		
		if !added {
			break // No more cards can be added
		}
	}
	
	return validCards
}

func validateCards(cards [][]int) bool {
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if countCommonSymbols(cards[i], cards[j]) != 1 {
				return false
			}
		}
	}
	return true
}

func formatCard(card []int) string {
	strSymbols := make([]string, len(card))
	for i, symbol := range card {
		strSymbols[i] = strconv.Itoa(symbol)
	}
	return "[" + strings.Join(strSymbols, ", ") + "]"
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dobble Card Generator")
	myWindow.Resize(fyne.NewSize(WindowWidth, WindowHeight))

	// Create UI instance
	ui := &AppUI{
		window:              myWindow,
		symbolsPerCardEntry: widget.NewEntry(),
		resultLabel:         widget.NewLabel("Enter the number of images per card. Supported values: 2, 3, 4, 5, 6, 8, 9..."),
		successLabel:        widget.NewLabel(""),
		errorLabel:          widget.NewLabel(""),
	}
	
	ui.symbolsPerCardEntry.SetPlaceHolder("Enter the number of images per card")
	ui.resultLabel.Wrapping = fyne.TextWrapWord
	
	// Create business logic instance
	generator := NewDobbleGenerator()
	
	// Setup cards list - simple approach with full refresh
	ui.cardsList = widget.NewList(
		func() int { return len(ui.cardItems) },
		func() fyne.CanvasObject {
			checkbox := widget.NewCheck("Processed", nil)
			label := widget.NewLabel("")
			return container.NewHBox(label, checkbox)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(ui.cardItems) {
				hbox := item.(*fyne.Container)
				label := hbox.Objects[0].(*widget.Label)
				checkbox := hbox.Objects[1].(*widget.Check)
				
				cardItem := ui.cardItems[id]
				
				// Update display
				if cardItem.Processed {
					label.SetText(fmt.Sprintf("Card %d", id+1))
				} else {
					label.SetText(fmt.Sprintf("Card %d: %s", id+1, formatCard(cardItem.Card)))
				}
				
				// Set checkbox without callback to avoid recursion
				checkbox.OnChanged = nil
				checkbox.SetChecked(cardItem.Processed)
				
				// Set callback that updates state and refreshes entire list
				checkbox.OnChanged = func(checked bool) {
					ui.cardItems[id].Processed = checked
					// Force full refresh of the list
					ui.cardsList.Refresh()
				}
			}
		},
	)

	// Generate button with simplified logic
	generateButton := widget.NewButton("Generate Cards", func() {
		// Validate input
		symbolsPerCard, err := validateInput(ui.symbolsPerCardEntry.Text)
		if err != nil {
			ui.showError(err.Error())
			return
		}

		// Show progress
		ui.resultLabel.SetText("Generating cards...")
		ui.successLabel.SetText("")
		ui.errorLabel.SetText("")
		
		// Generate cards
		cards := generator.Generate(symbolsPerCard)
		totalSymbols := generator.GetTotalSymbols()
		
		if len(cards) == 0 {
			ui.showError("Failed to generate cards")
			return
		}
		
		// Validate and filter cards
		allValid := validateCards(cards)
		displayCards := cards
		if !allValid {
			displayCards = filterValidCards(cards)
		}
		
		// Create card items
		ui.cardItems = make([]CardItem, len(displayCards))
		for i, card := range displayCards {
			ui.cardItems[i] = CardItem{
				Card:      card,
				Processed: false,
			}
		}
		
		// Show results
		if allValid {
			ui.showSuccess(fmt.Sprintf("Generated %d cards (%d symbols(images) per card, %d total symbols(images))", len(cards), symbolsPerCard, totalSymbols))
		} else {
			ui.showError(fmt.Sprintf("Showing %d of %d correct cards", len(displayCards), len(cards)))
		}
		
		// Show mathematical formula
		mathText := fmt.Sprintf("Formula: n=%d, symbols=n²+n+1=%d, cards=%d", symbolsPerCard-1, totalSymbols, totalSymbols)
		ui.resultLabel.SetText(mathText)

		ui.cardsList.Refresh()
	})

	// Create UI layout
	form := container.NewVBox(
		widget.NewLabel("Dobble Card Generator"),
		widget.NewCard("Instructions", "Enter the number of images per card. The application will automatically calculate the optimal number of symbols and generate all cards.", widget.NewLabel("")),
		widget.NewSeparator(),
		widget.NewLabel("Number of images per card:"),
		ui.symbolsPerCardEntry,
		generateButton,
		widget.NewSeparator(),
		ui.resultLabel,
		ui.successLabel,
		ui.errorLabel,
		widget.NewSeparator(),
		widget.NewLabel("Card List:"),
	)
	
	cardsScroll := container.NewScroll(ui.cardsList)
	cardsScroll.SetMinSize(fyne.NewSize(0, CardsListHeight))

	content := container.NewBorder(form, nil, nil, nil, cardsScroll)
	
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
} 