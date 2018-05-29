package main

import "testing"

func TestModifyClauses(t *testing.T) {
	/* Test-SAT-Formula:
	set --> 1
	   1 2 3
	   1 2 3
	--> Tests if satisfied formulas are correctly removed
	*/
	satProblem := constructSATProblem(&[]int{1, 2, 3}, &[]int{1, 2, 3})
	ModifyClauses(satProblem, 1)

	if len(satProblem.clauses) != 0 {
		t.Log("Testcase 1: not all clauses deleted")
		t.Fail()
	}

	/*
	set --> 1
		1 -2
		-2
	*/

	satProblem = constructSATProblem(&[]int{1, -2}, &[]int{-2})
	ModifyClauses(satProblem, 1)
	if satProblem.clauses[0][0] != -2 {
		t.Log("Testcase 2: Failed to delete first clause")
		t.Fail()
	}
	if satProblem.values[1] != 1 {
		t.Log("Testcase 2: Failed to set value correctly")
		t.Fail()
	}
}

func TestSolveDPLLnaive(t *testing.T) {
	/*
	1 2
	-1 -2
	 */
	 satProblem := constructSATProblem(&[]int{1,2}, &[]int{-1,-2})
	 if !SolveDPLLnaive(*satProblem,0){
	 	t.Log("Couldn not solve simple SAT")
	 	t.Fail()
	 }
}

func TestPureLiteralRule(t *testing.T) {
	satProblem := constructSATProblem(&[]int{1, 2}, &[]int{1, -2})
	if !PureLiteralRule(satProblem){
		t.Fail()
	}
}
