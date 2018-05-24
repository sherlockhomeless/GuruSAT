package main

import (
	"fmt"
)

func solveDPLLnaive(sat sat, try int) bool {
	// Priorities: Solved - Fail - Backtrack - UP - PL - S
	if try != 0 {
		ModifyClauses(&sat, try)
	}
	if solveRule(&sat) {
		return true
	}
	if failRule(&sat) {
		return false
	}
	unitPropagationRule(&sat)
	pureLiteralRule(&sat)
	splitRule(&sat)

	// Backtrack-Rule
	if (sat.varsNum == len(sat.values)) && (len(sat.clauses) > 0) {
		return false
	}
	return false
}

//Solved-Rule
func solveRule(sat *sat) bool {
	if len(sat.clauses) == 0 {
		fmt.Print(sat.values)
		return true
	}
	return false
}

// Fail-Rule
//TODO: Implementieren
func failRule(sat *sat) bool {
	if false {
		fmt.Println("Formula unsatisfiable")
		return false
	}
	return false
}

func unitPropagationRule(sat *sat) {
	for _, clause := range sat.clauses {
		if len(clause) == 1 {
			ModifyClauses(sat, clause[0])
		}
	}
}

// Pure-Literal-Rule
func pureLiteralRule(sat *sat) bool {
	// holds the value from the list that check occurances of literals
	var set int8
	// List to keep track of positive/negative occurence of literal
	polarityList := make([]int8, sat.varsNum+1)
	// checking all clauses and literals for their polarity
	for _, clause := range sat.clauses {
		for _, literal := range clause {
			set = polarityList[makeIntAbsolute(literal)]
			if set == 0 {
				if literal > 0 {
					polarityList[makeIntAbsolute(literal)] = 1
				} else {
					polarityList[makeIntAbsolute(literal)] = -1
				}
				// Literal has only one polarity over all clauses
			} else if (set > 0 && literal < 0) || (set < 0 && literal > 0) {
				polarityList[makeIntAbsolute(literal)] = -2
			}
		}
	}
	pureLiteral := 0
	// checking polarity list for pure Literals
	for literalNumber, literal := range polarityList {
		if (literal == 1) || (literal == -1) {
			pureLiteral = literalNumber
			ModifyClauses(sat, pureLiteral)
			break
		}

	}
	return false
}

func splitRule(satProblem *sat) bool {
	var literal int
	literal = satProblem.clauses[0][0]
	satProblem.values[makeIntAbsolute(literal)] = literal
	if solveDPLLnaive(*satProblem, literal) {
		return true
	}
	satProblem.values[makeIntAbsolute(literal)] = -1 * literal
	if solveDPLLnaive(*satProblem, literal*-1) {
		return true
	}
	/*

		splitSAT.values[makeIntAbsolute(literal)] = literal
		if solveDPLLnaive(splitSAT) {
			return true
		}
		splitSAT.values[makeIntAbsolute(literal)] = -1 * literal
		if solveDPLLnaive(splitSAT) {
			return true
		}
	*/
	return false
}
