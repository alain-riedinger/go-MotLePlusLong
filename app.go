package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ArgError signals an argument error
type ArgError struct {
	arg string
	err string
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("Argument: \"%s\" %s", e.arg, e.err)
}

func help() {
	fmt.Println()
	fmt.Println("+------------------+")
	fmt.Println("| Mot Le Plus Long |")
	fmt.Println("+------------------+")
	fmt.Println("A game that consists in finding a word with a set of 10 random letters")
	fmt.Println()
	fmt.Println("Syntax:")
	fmt.Println("-------")
	fmt.Println("MotLePlusLong [dico <dico> --start <ln> --end <ln>]")
	fmt.Println("- to play a game")
	fmt.Println("  no options needed, the application will do")
	fmt.Println()
	fmt.Println("- to generate the adapted dictionary from an 'unmunch' seed")
	fmt.Println("  dico")
	fmt.Println("  <dico>, path to the 'unmunch' dictionary file")
	fmt.Println("  --start <ln>, optional, line of parsing start, preceding are skipped")
	fmt.Println("  --end <ln>, optional, line of parsing end, following are skipped")
	fmt.Println()
}

func parseArgs(args []string) (string, int, int, error) {
	var dicoPath string // By default, set to nil
	lineStart := 0
	lineEnd := 0

	i := 1
	for i < len(args) {
		switch strings.ToLower(args[i]) {
		case "dico":
			if i+1 < len(args) {
				dicoPath = args[i+1]
				i += 2
			} else {
				return dicoPath, lineStart, lineEnd, &ArgError{"dico", "additional path missing"}
			}
		case "--start":
			if i+1 < len(args) {
				lineStart, _ = strconv.Atoi(args[i+1])
				i += 2
			} else {
				return dicoPath, lineStart, lineEnd, &ArgError{"dico", "additional line number missing"}
			}
		case "--end":
			if i+1 < len(args) {
				lineEnd, _ = strconv.Atoi(args[i+1])
				i += 2
			} else {
				return dicoPath, lineStart, lineEnd, &ArgError{"dico", "additional line number missing"}
			}
		default:
			return dicoPath, lineStart, lineEnd, &ArgError{args[i], "unknown"}
		}
	}
	return dicoPath, lineStart, lineEnd, nil
}

func main() {
	// Parse arguments
	dicoPath, lineStart, lineEnd, err := parseArgs(os.Args)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		help()
		return
	}

	if dicoPath != "" {
		// Construction of the dictionary
		parseUnmunchedDico(dicoPath, lineStart, lineEnd)
	} else {
		// The plain game

		// Load the dictionary
		dico := loadStrictDico("fr-mlpl-flat-strict.txt")

		// For unit tests only
		// idx := calcIndex("merde")
		// word, ok := dico[idx]
		// fmt.Println(word, ok)
		// countup(30, 30)
		// nb := promptVoyelles(4)
		// fmt.Println(nb)

		// Default number of voyelles
		nbVoyelles := 4

		for {
			clearScreen()
			displayTitle()
			// Launch a game with 30 seconds of think time
			nbVoyelles = newGame(dico, nbVoyelles, 30, 30)
			prompt("Want another game ?")
		}
	}
}

func findSolution(mot *Mot, dico map[[14]byte][]string, plaques string, sol chan *Solution) {
	// Initialize the recursive search root structure
	solution := NewSolution()
	solution.Current = plaques

	// Start the recursive resolution
	found := mot.SolveTirage(dico, *solution)

	// Send solution to the channel, while execution in parallel
	sol <- found
}

func newGame(dico map[[14]byte][]string, previous int, chars int, seconds int) int {
	nbVoyelles := promptVoyelles(previous)
	mot := NewMot()
	plaques := mot.GetPlaques(nbVoyelles)
	s := " "
	for i := 0; i < len(plaques); i++ {
		s += fmt.Sprintf(" %c ", plaques[i])
	}
	fmt.Println(" +----------------------------+")
	fmt.Println(strings.ToUpper(s))
	fmt.Println(" +----------------------------+")

	// Wait for some seconds of think time
	countup(chars, seconds)

	// Solution is searched during chrono time (it's shorter so no cheating)
	sol := make(chan *Solution)
	go findSolution(mot, dico, plaques, sol)
	found := <-sol
	close(sol)

	prompt("Want a solution ?")

	if found != nil {
		fmt.Printf("Best words found of: %d letters\n", found.BestLen)
		for _, w := range found.Best {
			fmt.Println(strings.ToUpper(w))
		}
	} else {
		fmt.Println(" No acceptable word has been found...")
	}
	fmt.Println()
	// To avoid typing this number in the next game
	return nbVoyelles
}

func displayTitle() {
	fmt.Println("  +---------------------------+")
	fmt.Println("  |  Le Mot le plus Long !    |")
	fmt.Println("  +---------------------------+")
}

func promptVoyelles(previous int) int {
	nb := previous
	fmt.Printf("How many voyelles ? [%d] (press Enter)\n", previous)
	var s string
	fmt.Scanln(&s)
	if s != "" {
		nb, _ = strconv.Atoi(s)
	}
	return nb
}

func prompt(msg string) {
	fmt.Printf("%s (press Enter)\n", msg)
	var s string
	fmt.Scanln(&s)
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
