package main

import (
	"fmt"
)

var solvedSAT *Sat

func SolveDPLLnaive(satProblem Sat, try int) bool {
	// Priorities: Solved - Fail - Backtrack - UP - PL - S
	if try != 0 {
		ModifyClauses(&satProblem, try)
	}
	if solveRule(&satProblem) {
		solvedSAT = &satProblem
		return true
	}
	if failRule(&satProblem) {
		return false
	}
	if unitPropagationRule(&satProblem) {
		return SolveDPLLnaive(satProblem, 0)
	}
	if PureLiteralRule(&satProblem) {
		return SolveDPLLnaive(satProblem, 0)
	}
	// Split Rule
	literal := satProblem.clauses[0][0]
	satProblem.values[makeIntAbsolute(literal)] = literal
	if SolveDPLLnaive(satProblem, literal) {
		return true
	}
	satProblem.values[makeIntAbsolute(literal)] = -1 * literal
	if SolveDPLLnaive(satProblem, literal*-1) {
		return true
	}

	// Backtrack-Rule; True if all variables are set and clauses are still left OR there are empty clauses which are not satisfiable anymore
	if clausesContainEmptyClause(satProblem.clauses) {
		return false
	}
	return SolveDPLLnaive(satProblem, 0)
}

//Solved-Rule
func solveRule(satProblem *Sat) bool {
	if len(satProblem.clauses) == 0 {
		if DEBUG {
			fmt.Printf("SAT solved with interpretation %v\n", satProblem.values)
		}
		return true
	}
	return false
}

// Fail-Rule
func failRule(sat *Sat) bool {
	// Solution has failed if empty clause is contained in clause-set
	/*for _, clause := range sat.clauses {
		if len(clause) == 0 {
			return true
		}
	}*/
	return false
}

func unitPropagationRule(sat *Sat) bool {
	for _, clause := range sat.clauses {
		if len(clause) == 1 {
			ModifyClauses(sat, clause[0])
			if DEBUG {
				fmt.Printf("UP used for clause %v\n", clause)
			}
			return true
		}
	}
	return false
}

// Pure-Literal-Rule
func PureLiteralRule(satProblem *Sat) bool {
	//reset counter, TODO: implement changes clauses in modifyClauses
	counterPositive, counterNegative := make([]int, satProblem.varCount+1), make([]int, satProblem.varCount+1)
	satProblem.counter[0] = counterPositive
	satProblem.counter[1] = counterNegative
	var pureLiteral int

	// count occurances of literals
	for _, clause := range satProblem.clauses {
		for _, literal := range clause {
			if literal > 0 {
				satProblem.counter[0][makeIntAbsolute(literal)]++
			} else {
				satProblem.counter[1][makeIntAbsolute(literal)]++
			}
		}
	}
	for index := range satProblem.counter[0] {
		posCounter, negCounter := satProblem.counter[0][index], satProblem.counter[1][index]
		// Positive or negative more then zero, the other 0
		if posCounter > 0 && negCounter == 0 {
			pureLiteral = index
			break
		} else if negCounter > 0 && posCounter == 0 {
			pureLiteral = index * -1
			break
		}

	}
	if pureLiteral != 0 {
		ModifyClauses(satProblem, pureLiteral)
		return true
	} else {
		return false
	}

}

func clausesContainEmptyClause(clauses [][]int) bool {
	for _, clause := range clauses {
		if len(clause) == 0 {
			return true
		}
	}
	return false
}
