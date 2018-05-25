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
	if (satProblem.varCount+1 == len(satProblem.values)) && (len(satProblem.clauses) > 0) || (clausesContainEmptyClause(satProblem.clauses)) {
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
func PureLiteralRule(sat *Sat) bool {
	var pureLiteral int
	var polTracker PolarityTracker
	polarityMap := map[int]PolarityTracker{}

	for _, clause := range sat.clauses {
		for _, literal := range clause {
			polTracker = polarityMap[makeIntAbsolute(literal)]
			//check if opposite polarity was found
			if polTracker.isLiteralBipolar(literal) {
				continue
			}
		}
	}

	for index:= 1; index <= sat.varCount+1; index++{
		polTracker = polarityMap[index]
		if !polTracker.Both{
			if polTracker.Pos{
				pureLiteral = index
			} else {
				pureLiteral = index * -1
			}
		}
	}
	if pureLiteral != 0 {
		ModifyClauses(sat, pureLiteral)
		return true
	} else {
		return false
	}
}



func pureLiteralRuleOld(sat *Sat) bool {
	// holds the value from the list that check occurances of literals
	var set int8
	// List to keep track of positive/negative occurence of literal
	polarityList := make([]int8, sat.varCount+1)
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
			ModifyClauses(sat, pureLiteral * int(polarityList[literalNumber]))
			return true
		}

	}
	return false
}

func clausesContainEmptyClause(clauses [][]int) bool {
	for _, clause := range clauses {
		if len(clause) == 0 {
			return true
		}
	}
	return false
}
