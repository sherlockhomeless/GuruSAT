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
	solveDPLLnaive(satFormula)
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

func solveDPLLnaive(sat sat) bool {
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
	var pureLiteral int
	polarityList := make([]int8, sat.varsNum+1)
	for _, clause := range sat.clauses {
		for _, literal := range clause {
			set := polarityList[literal]
			if set == 0 {
				if literal > 0 {
					polarityList[literal] = 1
				} else {
					polarityList[literal] = -1
				}
				// Literal has only one polarity over all clauses
			} else if (set > 0 && literal < 0) || (set < 0 && literal > 0) {
				polarityList[literal] = -2
			}
		}
	}
	for literalNumber, literal := range polarityList {
		if (literal == 1) || (literal == 0) {
			pureLiteral = literalNumber
			//TODO: Literal einsetzen
			break
		}

	}
	return false

}

func modifyClauses(sat sat, literal int){
	removelistClauses := make([]int, sat.clauseNum/3)
	removelistVariables := make([]int, sat.varsNum/10)
	for clauseNumber, clause := range sat.clauses{
		for _, clauseLiteral := range clause{
			if clauseLiteral == literal{
				removelistClauses = append(removelistClauses, clauseNumber)
			}
		}
	}

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
