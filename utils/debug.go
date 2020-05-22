package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/spacemeshos/node-mock/spacemesh"
)

func printLayers(layers *[]spacemesh.Layer) {
	fmt.Println("-- Layers --")

	for _, v := range *layers {
		fmt.Printf("%d - %s - %s\n", v.GetNumber(), hex.EncodeToString(v.GetHash()), v.GetStatus().String())
	}

	fmt.Println("-- ------ --")
}

func printAccounts(accounts *[]spacemesh.Account) {
	fmt.Println("- Accounts -")

	for _, v := range *accounts {
		fmt.Printf("TT %s\n", hex.EncodeToString(v.Address.Address))
	}

	fmt.Println("- -------- -")
}
