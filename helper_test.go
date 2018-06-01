package main

import "testing"

func TestSat_DeepCopySAT(t *testing.T) {
	s := Sat{}
	s.ReadFormula("formulas/medium_satisfiable")
	n := *s.DeepCopySAT()
	//checks if clauses are identical
	for clauseIndex, clause := range s.clauses {
		for literalIndex, _ := range clause {
			if s.clauses[clauseIndex][literalIndex] != n.clauses[clauseIndex][literalIndex] {
				t.Fail()
			}

		}
	}
	//check if value change affects both
	n.values[1]++
	if s.values[1] == n.values[1] {
		t.Fail()
	}
	// checks if clause change affects both
	n.clauses = n.clauses[:1]
	if len(s.clauses) == len(n.clauses) {
		t.Fail()
	}
}

func Test_CheckSolution(t *testing.T) {
	satProbelmSatisfiable, satProblemUnsatisfiable := Sat{}, Sat{}
	satProbelmSatisfiable.ReadFormula("formulas/test_0")
	satProblemUnsatisfiable.ReadFormula("formulas/test_unsatisfiable")
	SolveDPLLnaive(satProbelmSatisfiable, 0)
	if !CheckSolution(&satProbelmSatisfiable, solvedSAT.values) {
		t.Fail()
	}
	//TODO: Stupid test --> if no solution, &solvedSat = nil
	SolveDPLLnaive(satProblemUnsatisfiable, 0)
	if CheckSolution(&satProblemUnsatisfiable, solvedSAT.values) {
		t.Fail()
	}

}
