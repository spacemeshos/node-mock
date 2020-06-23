package utils

import (
	"encoding/hex"
	"fmt"

	v1 "github.com/spacemeshos/api/release/go/spacemesh/v1"
)

func printLayers(layers *[]v1.Layer) {
	fmt.Println("-- Layers --")

	for _, v := range *layers {
		fmt.Printf("%d - %s - %s\n", v.GetNumber(), hex.EncodeToString(v.GetHash()), v.GetStatus().String())
	}

	fmt.Println("-- ------ --")
}

func printAccounts(accounts *[]v1.Account) {
	fmt.Println("- Accounts -")

	for _, v := range *accounts {
		fmt.Printf("TT %s\n", hex.EncodeToString(v.Address.Address))
	}

	fmt.Println("- -------- -")
}
