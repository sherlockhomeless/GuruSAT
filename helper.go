package main

func makeIntAbsolute(x int) int{
	if x < 0{
		return x * -1
	}
	return x
}
