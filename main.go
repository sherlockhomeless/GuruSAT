package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/fatih/color"
)

var DEBUG bool
var CUR_SPLIT_RULE int
var formulas = []string{"easy", "test_0", "medium_satisfiable", "flat200-1.cnf", "test_unsatisfiable"}
var splitRules []func(*Sat)bool

func main() {
	DEBUG = true
	satFormula := Sat{}
	// choose which SAT-Formula to check
	if len(os.Args)>1{
		satPath := os.Args[1]
		satFormula.ReadFormula(satPath)
	} else {
		satFormula.ReadFormula("formulas/" + formulas[3])

	}
	// choose which split-Rule
	splitNaiv := splitRuleChronological
	splitCounting := SplitRuleWithCoutingOfLiteralOccurances
	splitRules = []func(sat *Sat)bool{splitCounting,splitNaiv}
	for index := 0; index < len(splitRules); index++{
		start := time.Now()
		CUR_SPLIT_RULE = index
		if SolveDPLLnaive(satFormula, 0) {
			color.Blue("SAT satisfiable with interpretation: %v\n", solvedSAT.values)
		} else {
			color.Red("SAT unsatisfiable")
		}
		elapsed := time.Since(start)
		fmt.Printf("Solving took %s with Split-Rule %d\n", elapsed, index)
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
		if line[0] == 'c' {
			continue
		} else if line[0] == 'p' {
			split := strings.Split(line, " ")
			sat.varCount, _ = strconv.Atoi(split[2])
			sat.clauseCount, _ = strconv.Atoi(split[3])
			fmt.Printf("Found %d variables and %d clauses\n", sat.varCount, sat.clauseCount)
			sat.values = make([]int, sat.varCount+1)
			sat.clauses = make([][]int, sat.clauseCount)
			counterPositive, counterNegative := make([]int, sat.varCount+1), make([]int, sat.varCount+1)
			sat.counter = [2][]int{}
			sat.counter[0] = counterPositive
			sat.counter[1] = counterNegative
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
				sat.clauses[formulaIndex] = append(sat.clauses[formulaIndex], varValue)
			}
			formulaIndex++
		}
	}

}

