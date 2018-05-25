package main

func makeIntAbsolute(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func constructSATProblem(clauses ...*[]int) *Sat {
	varCount, clauseCount := 0, 0
	varAlreadyCounter := make (map[int]bool)
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
	return &Sat{varCount:varCount, clauseCount:clauseCount, clauses:clauseList, values:valList}
}

type Sat struct {
	varCount, clauseCount int
	clauses               [][]int // conjunctive clause set
	values                []int   // variables set
}

type PolarityTracker struct {
	Pos, Neg, Both bool
}

func (p *PolarityTracker) isLiteralBipolar(literal int) bool{
	if literal > 0{
		p.Pos = true
	} else {
		p.Neg = true
	}
	if p.Both || p.Pos && p.Neg {
		p.Both = true
		return true
	} else {
		return false
	}
}

