package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	satFormula := sat{}
	satFormula.readFormula("./formulas/" + "easy")
	solveDPLL(satFormula)
}

func (sat *sat) readFormula(path string) {
	formulaFile, err := os.Open(path)
	defer formulaFile.Close()
	check(err)
	scanner := bufio.NewScanner(formulaFile)
	scanner.Split(bufio.ScanLines)
	formulaInxex := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'c' {
			continue
		} else if line[0] == 'p' {
			split := strings.Split(line, " ")
			sat.varsNum, _ = strconv.Atoi(split[2])
			sat.clauseNum, _ = strconv.Atoi(split[3])
			fmt.Printf("Found %d variables and %d clauses", sat.varsNum, sat.clauseNum)
			sat.values = make([]int, 0) // set variables
			sat.clauses = make([][]int, sat.clauseNum)
		} else {
			var varValue int
			line = strings.Replace(line, " 0", "", 1)
			split := strings.Split(line, " ")
			// super akward formula importing loop
			for _, variable := range split {
				varValue, _ = strconv.Atoi(variable)
				sat.clauses[formulaInxex] = append(sat.clauses[formulaInxex], varValue)
			}
			formulaInxex++
		}
	}

}

func solveDPLL(sat sat) bool {
	// Priorities: Solved - Fail - Backtrack - UP - PL - S

	//Solved-Rule
	if len(sat.clauses) == 0 {
		fmt.Print(sat.values)
		return true
	}

	// Fail-Rule
	if false {
		fmt.Println("Formula unsatisfiable")
	}

	// Backtrack-Rule
	if (sat.varsNum == len(sat.values)) && (len(sat.clauses) > 0) {
		return false
	}

	// Unit-Propagation necessary?

	// Pure-Literal-Rule
	polarityList := make([]int8, sat.varsNum+1)
	for _, clause := range sat.clauses {
		for literal := range clause {
			//TODO: weiter
			set := polarityList[literal]
		}
	}
	return false

}

type sat struct {
	varsNum, clauseNum int
	clauses            [][]int // conjunctive clause set
	values             []int   // variables set
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
