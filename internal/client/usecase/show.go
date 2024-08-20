package usecase

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/nextlag/keeper/internal/client/usecase/viewsets"
)

const (
	showAllData  = "a"
	showCards    = "c"
	showLogins   = "l"
	showNotes    = "n"
	showBinaries = "b"
)

// ShowVault displays the user's vault contents based on the specified option.
func (uc *ClientUseCase) ShowVault(userPassword, showVaultOption string) {
	if !uc.verifyPassword(userPassword) {
		return
	}

	switch showVaultOption {
	case showAllData:
		uc.showCards(uc.repo.LoadCards())
		uc.showLogins(uc.repo.LoadLogins())
		uc.showNotes(uc.repo.LoadNotes())
		uc.showBinaries(uc.repo.LoadBinaries())
	case showCards:
		uc.showCards(uc.repo.LoadCards())
	case showLogins:
		uc.showLogins(uc.repo.LoadLogins())
	case showNotes:
		uc.showNotes(uc.repo.LoadNotes())
	case showBinaries:
		uc.showBinaries(uc.repo.LoadBinaries())
	}
}

// showCards prints out a list of cards to the console.
func (uc *ClientUseCase) showCards(cards []viewsets.CardForList) {
	color.Yellow("Users cards:")
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, card := range cards {
		fmt.Printf("ID: %s name: %s brand: %s\n",
			yellow(card.ID),
			yellow(card.Name),
			yellow(card.Brand))
	}
	fmt.Printf("Total %s cards\n", yellow(len(cards)))
}

// showLogins prints out a list of logins to the console.
func (uc *ClientUseCase) showLogins(logins []viewsets.LoginForList) {
	color.Yellow("Users logins:")
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, login := range logins {
		fmt.Printf("ID: %s name: %s uri: %s\n",
			yellow(login.ID),
			yellow(login.Name),
			yellow(login.URI))
	}
	fmt.Printf("Total %s logins\n", yellow(len(logins)))
}

// showNotes prints out a list of notes to the console.
func (uc *ClientUseCase) showNotes(notes []viewsets.NoteForList) {
	color.Yellow("Users notes:")
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, note := range notes {
		fmt.Printf("ID: %s name: %s\n",
			yellow(note.ID),
			yellow(note.Name))
	}
	fmt.Printf("Total %s notes\n", yellow(len(notes)))
}

// showBinaries prints out a list of binary files to the console.
func (uc *ClientUseCase) showBinaries(binaries []viewsets.BinaryForList) {
	color.Yellow("Users files:")
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, binary := range binaries {
		fmt.Printf("ID: %s name: %s file_name: %s\n",
			yellow(binary.ID),
			yellow(binary.Name),
			yellow(binary.FileName))
		fmt.Printf("Total %s binaries\n", yellow(len(binaries)))
	}
}
