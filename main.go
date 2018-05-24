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
	solveDPLLnaive(satFormula,0)
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
			fmt.Printf("Found %d variables and %d clauses\n", sat.varsNum, sat.clauseNum)
			sat.values = make([]int, sat.varsNum +1)
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

// Modifies the Clauseset sat.clauses according to value given to literal
func ModifyClauses(sat *sat, literal int) {
	// holds all the clauses that have to be deleted to represent new literal interpretation
	//removelistClauses := make([]int, sat.clauseNum/3)
	//removelistVariables := make([]int, sat.varsNum/10)
	var removelistClauses []int
	var removelistVariables []int
	for clauseNumber, clause := range sat.clauses {
		for literalNumber, clauseLiteral := range clause {
			// Remove clause because it is made true through literal
			if clauseLiteral == literal {
				removelistClauses = append(removelistClauses, clauseNumber)
				break // if one literal of disjunction true => whole clause is true
			}
			// Remove literal with opposite polarity since can't be part of the solution
			if clauseLiteral == literal*-1 {
				removelistVariables = append(removelistVariables, literalNumber +1)
			}
		}
		for _, removeVariable := range removelistVariables {
			if removeVariable + 1 <= len(clause) {
				clause = append(clause[:removeVariable], clause[removeVariable+1:]...)
			} else {
				clause = clause[:len(clause)-1]
			}
		}
	}
	for _, removeClause := range removelistClauses{
		if removeClause+1 < len(sat.clauses) {
			sat.clauses = append(sat.clauses[:removeClause], sat.clauses[removeClause+1:]...)
		} else {
			sat.clauses = sat.clauses[:len(sat.clauses)-1]
		}
	}

	fmt.Printf("Literal set to %d\n", literal)
	sat.values[makeIntAbsolute(literal)] = literal

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

/*func testSolve(){
	sat := sat{}
	sat.readFormula("/formulas/test_0")
	solveDPLLnaive(sat, 0)
}*/