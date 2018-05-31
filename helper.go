package main

import "fmt"

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

func makeIntAbsolute(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func constructSATProblem(clauses ...*[]int) *Sat {
	varCount, clauseCount := 0, 0
	varAlreadyCounter := make(map[int]bool)
	var clauseList [][]int
	for _, clause := range clauses {
		clauseList = append(clauseList[:], *clause)
		clauseCount++
		for _, variable := range *clause {
			if !varAlreadyCounter[makeIntAbsolute(variable)] {
				varAlreadyCounter[makeIntAbsolute(variable)] = true
				varCount++
			}
		}
	}
	valList := make([]int, varCount+1)
	counter := [2][]int{}
	pos, neg := make([]int, varCount+1), make([]int, varCount+1)
	counter[0], counter [1] = pos, neg
	return &Sat{varCount: varCount, clauseCount: clauseCount, clauses: clauseList, values: valList, counter:counter}
}

type Sat struct {
	varCount, clauseCount int
	clauses               [][]int // conjunctive clause set
	values                []int   // variables set
	counter               [2][]int
}

func (s *Sat) DeepCopySAT() *Sat{

	// Initing empty Structure
	newSATProblem := Sat{}
	newSATProblem.varCount = s.varCount
	newSATProblem.clauseCount = s.clauseCount
	newSATProblem.clauses = make([][]int,len(s.clauses))
	newSATProblem.values = make([]int, s.varCount+1)
	newSATProblem.counter = [2][]int{}
	pos, neg := make([]int, s.varCount+1), make([]int, s.varCount+1)
	newSATProblem.counter[0], newSATProblem.counter[1] = pos, neg

	// Copying values
	//copy(newSATProblem.clauses, s.clauses)
	for indexC, clause := range s.clauses{
		//copy(newSATProblem.clauses[index], s.clauses[index])
		for _, literal:= range clause{
			newSATProblem.clauses[indexC] = append(newSATProblem.clauses[indexC], literal )
		}
	}
	copy(newSATProblem.values, s.values)

	return &newSATProblem

}

