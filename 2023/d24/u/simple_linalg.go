package u


type Two_by_two_matrix [2][2]float64
type Two_by_one_matrix [2][1]float64

func (a Two_by_two_matrix) Det() float64 {
	return a[0][0]*a[1][1] - a[0][1]*a[1][0]
}

func (a Two_by_two_matrix) Invert() Two_by_two_matrix {
	det := a.Det()
	return Two_by_two_matrix{
		{a[1][1], -a[0][1]},
		{-a[1][0], a[0][0]},
	}.Scale(1/det)
}

func (a Two_by_two_matrix) Add(b Two_by_two_matrix) Two_by_two_matrix {
	return Two_by_two_matrix{
		{a[0][0]+b[0][0], a[0][1]+b[0][1]},
		{a[1][0]+b[1][0], a[1][1]+b[1][1]},
	}
}

func (a Two_by_two_matrix) Scale(s float64) Two_by_two_matrix {
	return Two_by_two_matrix{
		{a[0][0]*s, a[0][1]*s},
		{a[1][0]*s, a[1][1]*s},
	}
}

func (a Two_by_two_matrix) Mult_by_column(b Two_by_one_matrix) Two_by_one_matrix {
	return Two_by_one_matrix{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0]},
		{a[1][0]*b[0][0] + a[1][1]*b[1][0]},
	}
}

func (a Two_by_two_matrix) Mult_by_two_by_two(b Two_by_two_matrix) Two_by_two_matrix {
	return Two_by_two_matrix{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
		{a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
	}
}

func (a Two_by_two_matrix) SubstituteColumn(b Two_by_one_matrix, col int) Two_by_two_matrix {
	if col == 0 {
		return Two_by_two_matrix{
			{b[0][0], a[0][1]},
			{b[1][0], a[1][1]},
		}
	} else {
		return Two_by_two_matrix{
			{a[0][0], b[0][0]},
			{a[1][0], b[1][0]},
		}
	}
}

type Linear_System struct {
	// A is the matrix of coefficients
	A Two_by_two_matrix
	// B is the matrix of constants
	B Two_by_one_matrix
}

func (a Linear_System) Solve() Two_by_one_matrix {
	return a.A.Invert().Mult_by_column(a.B)
}



