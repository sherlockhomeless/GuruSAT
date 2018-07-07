package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var DEBUG bool
var CUR_SPLIT_RULE int
var formulas = []string{"easy", "test_0", "medium_satisfiable", "flat200-1.cnf", "test_unsatisfiable", "uf20-01.cnf"}
var splitRules []func(*Sat) bool

func main() {
	DEBUG = false
	CUR_SPLIT_RULE = 0
	if len(os.Args) > 2{
		CUR_SPLIT_RULE,_ = strconv.Atoi(os.Args[2])
	}

	splitRules = []func(sat *Sat) bool{splitRuleChronological, SplitRuleWithCoutingOfLiteralOccurances, SplitRuleWithCoutingOfLiteralOccurancesAndShortClausePreferation}

	files,_ := ioutil.ReadDir(".")
	formulaCount := len(files)
	satFormulas := make([]Sat, formulaCount)
	// choose which SAT-Formula to check
	if len(os.Args) > 1 {
		// this solves all formulas in the cwd
		if os.Args[1] == "." {
			formulas, _ := ioutil.ReadDir(".")
			for index, f := range formulas {
				satFormulas[index].ReadFormula(f.Name())
			}
		} else {
			satPath := os.Args[1]
			satFormulas[0].ReadFormula(satPath)
		}
	} else {
		satFormulas[0].ReadFormula("formulas/" + formulas[0])

	}
	color.Red("Start Solving")
	for _, formula := range satFormulas {
		if formula.varCount == 0 {
			break
		}
		start := time.Now()
		if SolveDPLLnaive(formula, 0) {
			color.Green("SAT satisfiable with interpretation: %v\n", solvedSAT.values[1:])
		} else {
			color.Red("SAT unsatisfiable\n")
		}
		elapsed := time.Since(start)
		color.Magenta("Solving took %s seconds\n", elapsed)
	}

}

func (sat *Sat) ReadFormula(path string) {
	formulaFile, err := os.Open(path)
	defer formulaFile.Close()
	check(err)
	scanner := bufio.NewScanner(formulaFile)
	scanner.Split(bufio.ScanLines)
	formulaIndex := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == 'c' || line[0] == '%' || line[0] == '0' {
			continue
		} else if line[0] == 'p' {
			line = strings.Replace(line, "  ", " ", 20)
			splits := strings.Split(line, " ")
			sat.varCount, _ = strconv.Atoi(splits[2])
			sat.clauseCount, _ = strconv.Atoi(splits[3])
			fmt.Printf("Found %d variables and %d clauses\n", sat.varCount, sat.clauseCount)
			sat.values = make([]int, sat.varCount+1)
			sat.clauses = make([][]int, sat.clauseCount)
			counterPositive, counterNegative := make([]int, sat.varCount+1), make([]int, sat.varCount+1)
			sat.counter = [2][]int{}
			sat.counter[0] = counterPositive
			sat.counter[1] = counterNegative
		} else if lineNotValid(line) {
			continue
		} else {
			var varValue int
			line = strings.Replace(line, " 0", "", 1)
			split := strings.Split(line, " ")
			// super akward formula importing loop
			for _, variable := range split {
				varValue, err = strconv.Atoi(variable)
				if err != nil {
					continue
				}
				if formulaIndex == len(sat.clauses) {
					print(line + "\n")
				}
				sat.clauses[formulaIndex] = append(sat.clauses[formulaIndex], varValue)
			}
			formulaIndex++
		}
	}
}

func lineNotValid(l string) bool {
	validChars := "1234567890- !"
	for _, char := range l {
		for _, valChar := range validChars {
			if char == valChar {
				break
			} else if valChar == '!' {
				return true
			}
		}
	}
	return false
}
