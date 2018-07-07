package main

import (
	"github.com/fatih/color"
)

var solvedSAT *Sat

func SolveDPLLnaive(satProblem Sat, try int) bool {
	if DEBUG {
		color.Yellow("Values set: %v\n", satProblem.values)
		color.Yellow("Clauses left: %d\n", len(satProblem.clauses))
	}
	// Priorities: Solved - Fail - Backtrack - UP - PL - S
	if try != 0 {
		ModifyClauses(&satProblem, try)
	}
	if solveRule(&satProblem) {
		solvedSAT = &satProblem
		return true
	}
	// Backtrack-Rule;
	if clausesContainEmptyClause(satProblem.clauses) {
		if DEBUG {
			color.Red("Starting Backtracking\n")
		}
		return false
	}

	if unitPropagationRule(&satProblem) {
		return SolveDPLLnaive(satProblem, 0)
	}
	if PureLiteralRule(&satProblem) {
		return SolveDPLLnaive(satProblem, 0)
	}

	return splitRules[CUR_SPLIT_RULE](&satProblem)

}

func splitRuleChronological(satProblem *Sat) bool {
	literal := satProblem.clauses[0][0]
	satProblem.values[makeIntAbsolute(literal)] = literal
	satPositive, satNegative := satProblem.DeepCopySAT(), satProblem.DeepCopySAT()
	return SolveDPLLnaive(*satPositive, literal) || SolveDPLLnaive(*satNegative, literal*-1)
}

func SplitRuleWithCoutingOfLiteralOccurances(satProblem *Sat) bool {
	var max, literal int
	var positiveLiteral bool
	// array that contains the absolute value of the sum of the positive and negative occurances of a literal
	// adjustedLiteralOccurances := make([]int, satProblem.varCount)
	for index := range satProblem.counter[0] {
		num := satProblem.counter[0][index] + satProblem.counter[1][index]
		if num < 0 {
			num = num * -1
			positiveLiteral = false
		}
		if num > max {
			max = num
			if positiveLiteral {
				literal = index
			} else {
				literal = index * -1
			}
		}
		positiveLiteral = true
	}
	satPositive, satNegative := satProblem.DeepCopySAT(), satProblem.DeepCopySAT()
	//TODO: If backtracking is needed here, must unsuitable next interpretation is chosen
	return SolveDPLLnaive(*satPositive, literal) || SolveDPLLnaive(*satNegative, literal*-1)

}

func SplitRuleWithCoutingOfLiteralOccurancesAndShortClausePreferation(satProblem *Sat) bool {
	// array that contains the absolute value of the sum of the positive and negative occurances of a literal
	// adjustedLiteralOccurances := make([]int, satProblem.varCount)
	clauseSort := make([][][]int, 5)
	//Sorts Clauses depended on length; e.g. clause [-12 21 9] is in clauseSort[2][120]
	for _, clause := range satProblem.clauses {
		lengthOfClause := len(clause)
		if lengthOfClause >= 5 {
			clauseSort[4] = append(clauseSort[4], clause)
		} else {
			clauseSort[lengthOfClause-1] = append(clauseSort[lengthOfClause-1], clause)
		}
	}
	//Adds weight to pure occurance counting; e.g. 2 in [2 13] counts
	multiplicator := 2
	for i := len(clauseSort) - 1; i >= 0; i-- {
		for _, clause := range clauseSort[i] {
			for _, variable := range clause {
				if variable > 0 {
					satProblem.counter[0][variable] = satProblem.counter[0][variable] + 1*multiplicator
				} else {
					varAbs := variable * -1
					satProblem.counter[1][varAbs] = satProblem.counter[1][varAbs] + 1*multiplicator
				}
			}
		}
		multiplicator *= 2
	}

	var max, literal int
	var positiveLiteral bool

	for index := range satProblem.counter[0] {
		num := satProblem.counter[0][index] + satProblem.counter[1][index]
		if num < 0 {
			num = num * -1
			positiveLiteral = false
		}
		if num > max {
			max = num
			if positiveLiteral {
				literal = index
			} else {
				literal = index * -1
			}
		}
		positiveLiteral = true
	}

	satPositive, satNegative := satProblem.DeepCopySAT(), satProblem.DeepCopySAT()
	return SolveDPLLnaive(*satPositive, literal) || SolveDPLLnaive(*satNegative, literal*-1)

}

//Solved-Rule
func solveRule(satProblem *Sat) bool {
	if len(satProblem.clauses) == 0 {
		if DEBUG {
			color.Blue("SAT solved with interpretation %v\n", satProblem.values)
		}
		return true
	}
	return false
}

func unitPropagationRule(sat *Sat) bool {
	for _, clause := range sat.clauses {
		if len(clause) == 1 {
			ModifyClauses(sat, clause[0])
			if DEBUG {
				color.Green("UP used for clause %v\n", clause)
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
		if DEBUG {
			color.Green("Pure Literal %d was found\n", pureLiteral)
		}
		return true
	} else {
		return false
	}

}
