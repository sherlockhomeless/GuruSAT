package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var DEBUG bool
var formulas = []string{"easy","test_0","medium_satisfiable"}
func main() {
	DEBUG = true
	satFormula := Sat{}
	satFormula.readFormula("./formulas/" + formulas[0])
	a := SolveDPLLnaive(satFormula, 0)
	if a {
		print("cool")
	}
}

func (sat *Sat) readFormula(path string) {
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
			sat.counter = make([][]int,2,)
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

// Modifies the Clauseset Sat.clauses according to value given to literal
func ModifyClausesOld(sat *Sat, literal int) {
	// holds all the clauses that have to be deleted to represent new literal interpretation
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
				removelistVariables = append(removelistVariables, literalNumber+1)
			}
		}
		for _, removeVariable := range removelistVariables {
			if removeVariable+1 <= len(clause) {
				clause = append(clause[:removeVariable], clause[removeVariable+1:]...)
			} else {
				clause = clause[:len(clause)-1]
			}
		}
	}
	for _, removeClause := range removelistClauses {
		if removeClause+1 < len(sat.clauses) {
			sat.clauses = append(sat.clauses[:removeClause], sat.clauses[removeClause+1:]...)
		} else {
			sat.clauses = sat.clauses[:len(sat.clauses)-1]
		}
	}

	fmt.Printf("Literal set to %d\n", literal)
	sat.values[makeIntAbsolute(literal)] = literal

}

func ModifyClauses(satProblem *Sat, literalSet int) {
	deleteListClauses := make([]bool, satProblem.clauseCount)
	deleteVariableIndex := -1
	counter := 0
	for clauseIndex, clause := range satProblem.clauses {
		// Deletes clause if it was solved by literal set
		if isClauseSolved(&clause, literalSet) {
			deleteListClauses[clauseIndex] = true
			counter++
			if DEBUG {
				fmt.Printf("Deleting clause %v because literal %d was set\n", clause, literalSet)
			}
			continue
		}
		deleteVariableIndex = doesClauseContainLiteralInOpPolarity(&clause, literalSet)
		if deleteVariableIndex != -1 {
			// Delets literal from clause if it has opposite polarity to literal set
			if DEBUG {
				fmt.Printf("Deleting literal %d from clause %v, because literal was set to %d\n", satProblem.clauses[clauseIndex][deleteVariableIndex], clause, literalSet)
			}
			satProblem.clauses[clauseIndex] = deleteLiteralFromClause(clause, deleteVariableIndex)

		}
	}
	counter = 0
	for index, clauseToDelete := range deleteListClauses {
		if clauseToDelete {
			satProblem.clauses = deleteClauseFromFormula(satProblem.clauses, index-counter)
			counter++
		}

	}
	satProblem.values[makeIntAbsolute(literalSet)] = literalSet
}

// Returns if clause is solved under current interpretation
func isClauseSolved(clause *[]int, literal int) bool {
	for _, curLiteral := range *clause {
		if curLiteral == literal {
			return true
		}
	}
	return false
}

// Returns if clause contains literal in opposite polarity to set literal
func doesClauseContainLiteralInOpPolarity(clause *[]int, literal int) int {
	for index, curLiteral := range *clause {
		if curLiteral*-1 == literal {
			return index
		}
	}
	return -1
}

func deleteLiteralFromClause(clause []int, index int) []int {
	newClause := append(clause[:index], clause[index+1:]...)
	if len(newClause) == 0 {
		return nil
	}
	return newClause
}

func deleteClauseFromFormula(clauses [][]int, index int) (clausesnew [][]int) {
	clausesnew = append(clauses[:index], clauses[index+1:]...)
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
