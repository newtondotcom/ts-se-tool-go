package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func displayMainMenu() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    MAIN MENU                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println("1. Set money amount")
	fmt.Println("2. Set XP points")
	fmt.Println("3. Set all skills to max")
	fmt.Println("4. Buy all garages")
	fmt.Println("5. Upgrade all garages")
	fmt.Println("6. Populate all garages with random trucks")
	fmt.Println("7. Recruit employees and populate all trucks")
	fmt.Println("8. Save and exit")
	fmt.Println("9. Exit without saving")
	fmt.Print("\nSelect option (1-9): ")
}

func getUserChoice() int {
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		return -1
	}
	return choice
}

func promptMoney() (int64, error) {
	fmt.Print("Enter money amount: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	amount, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %v", err)
	}
	return amount, nil
}

func promptXP() (uint32, error) {
	fmt.Print("Enter XP points: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	xp, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid XP: %v", err)
	}
	return uint32(xp), nil
}

func confirmContinue() bool {
	fmt.Print("\nDo you want to make another change? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
