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

func (uc *ClientUseCase) showCards(cards []viewsets.CardForList) {
	color.Yellow("Users cards:")
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, card := range cards {
		fmt.Printf("ID:%s name:%s brand:%s number:%s\n",
			yellow(card.ID),
			yellow(card.Name),
			yellow(card.Brand))
	}
	fmt.Printf("Total %s cards\n", yellow(len(cards)))
}

func (uc *ClientUseCase) showLogins(logins []viewsets.LoginForList)      {}
func (uc *ClientUseCase) showNotes(notes []viewsets.NoteForList)         {}
func (uc *ClientUseCase) showBinaries(binaries []viewsets.BinaryForList) {}
