package main

import "testing"

func TestModifyClauses(t *testing.T) {
	/* Test-SAT-Formula:
	   1 2 3
	   1 2 3
	--> Tests if satisfied formulas are correctly removed
	*/
	varsNum := 3
	clausesNum := 2
	clause1 := []int{1,2,3}
	clause2 := []int{1,2,3}
	clauses := [][]int{clause1, clause2}
	satProblem := sat{varsNum:varsNum, clauseNum:clausesNum, clauses:clauses, values:make([]int, varsNum) }
	ModifyClauses(&satProblem,1)
	if len(satProblem.clauses) != 0{
		t.Fail()
	}
	/*
	1 2
	 */
	varsNum = 2
	clausesNum = 1
	clause1 = []int{1,2}
	satProblem = sat{varsNum:varsNum, clauseNum:clausesNum, clauses:clauses, values:make([]int, varsNum)}
	ModifyClauses(&satProblem, 1)
	if satProblem.clauses[0][0] != 2 {
		t.Fail()
	}
	if satProblem.values[1] != 1 {
		t.Fail()
	}

}
