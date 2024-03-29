package puzzle

import (
	"encoding/json"
	"fmt"
)

const (
	Up    = 1
	Down  = 2
	Left  = 3
	Right = 4
)

type Config struct {
	mat              [][]int8
	i, j             *int
	neighbors        []*Config
	parent           *Config
	Depth            int
	Cost             float64
	HKey             float64
	Explored         bool
	key              *string
	ReverseNeighbors bool
	Move             int
}

func (c Config) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Matrix [][]int8
		Move   int
	}{
		Matrix: c.mat,
		Move:   c.Move,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func NewConfig(m [][]int8) *Config {

	return &Config{
		mat: m,
	}

}

func (c *Config) Len() int {
	return len(c.mat)
}

func (c *Config) getInvCountFinal() int {
	invCount := 0
	lenC := c.Len() * c.Len()
	for i := 0; i < lenC-1; i++ {
		for j := i + 1; j < lenC; j++ {

			a, b := c.mat[j/c.Len()][j%c.Len()], c.mat[i/c.Len()][i%c.Len()]
			if a != 0 &&
				b != 0 &&
				b > a {
				invCount++
			}
		}
	}
	return invCount
}

func (c *Config) IsSolvable() bool {

	return c.getInvCountFinal()%2 == 0
}

func (c *Config) Less(c2 *Config) bool {

	if c.key != nil && c2.key != nil {
		return *c.key < *c2.key
	}

	for i := range c.mat {
		for j := range c.mat[i] {
			if c.mat[i][j] < c2.mat[i][j] {
				return true
			}
		}
	}

	return false
}

func (c *Config) IsSame(c2 *Config) bool {
	if c.key != nil && c2.key != nil {
		return *c.key == *c2.key
	}

	if len(c.mat) != len(c2.mat) {
		return false
	}

	for i := range c2.mat {

		if len(c.mat[i]) != len(c2.mat[i]) {
			return false
		}

		for j := range c2.mat[i] {
			if c.mat[i][j] != c2.mat[i][j] {
				return false
			}
		}
	}

	return true
}

func (c *Config) FindEmpty() (int, int) {

	if c.i != nil && c.j != nil {
		return *c.i, *c.j
	}
	for i := range c.mat {
		for j := range c.mat[i] {
			if c.mat[i][j] == Empty {
				z, f := i, j

				c.i, c.j = &z, &f

				return i, j
			}
		}
	}

	panic("No Empty Evalue")
}

func (c *Config) UpdatePosition(i, j int) {
	c.i, c.j = &i, &j
}

func (c *Config) Copy() *Config {

	m := make([][]int8, len(c.mat))
	for i := range m {
		m[i] = make([]int8, len(c.mat[0]))
	}
	c.Iterate(func(i, j *int, val *int8) {

		m[*i][*j] = *val
	})

	i, j := *c.i, *c.j
	return &Config{
		mat:              m,
		Depth:            c.Depth,
		Cost:             c.Cost,
		i:                &i,
		j:                &j,
		ReverseNeighbors: c.ReverseNeighbors,
	}
}

func (c *Config) Iterate(f func(*int, *int, *int8)) {
	for i := range c.mat {
		for j := range c.mat[i] {

			f(&i, &j, &c.mat[i][j])

		}
	}
}

func (c *Config) String() string {
	str := ""
	for i := range c.mat {
		for j := range c.mat[i] {
			str += fmt.Sprintf("%d ", c.mat[i][j])
		}
		str += "\n"
	}
	return str
}

func (c *Config) Println() {
	fmt.Println(c.String())
}

func (c *Config) MoveUp() *Config {

	i, j := c.FindEmpty()

	if i > 0 {
		c = c.Copy()
		c.mat[i-1][j], c.mat[i][j] = c.mat[i][j], c.mat[i-1][j]
		c.UpdatePosition(i-1, j)
	} else {
		return nil
	}

	return c
}

func (c *Config) MoveDown() *Config {

	i, j := c.FindEmpty()

	if i < len(c.mat)-1 {
		c = c.Copy()
		c.mat[i+1][j], c.mat[i][j] = c.mat[i][j], c.mat[i+1][j]
		c.UpdatePosition(i+1, j)
	} else {
		return nil
	}

	return c
}

func (c *Config) MoveLeft() *Config {

	i, j := c.FindEmpty()

	if j > 0 {
		c = c.Copy()
		c.mat[i][j-1], c.mat[i][j] = c.mat[i][j], c.mat[i][j-1]
		c.UpdatePosition(i, j-1)
	} else {
		return nil
	}
	return c
}

func (c *Config) MoveRight() *Config {

	i, j := c.FindEmpty()

	if j < len(c.mat[i])-1 {
		c = c.Copy()
		c.mat[i][j+1], c.mat[i][j] = c.mat[i][j], c.mat[i][j+1]
		c.UpdatePosition(i, j+1)
	} else {
		return nil
	}
	return c
}

func (c *Config) Neighbors() []*Config {

	if c.neighbors != nil {
		return c.neighbors
	}

	cU, cD, cL, cR := c.MoveUp(), c.MoveDown(), c.MoveLeft(), c.MoveRight()

	neighbors := make([]*Config, 0)

	if cU != nil {
		cU.parent = c
		cU.Depth++
		cU.Cost++
		cU.Move = Up
		neighbors = append(neighbors, cU)
	}

	if cD != nil {
		cD.parent = c
		cD.Depth++
		cD.Cost++
		cD.Move = Down
		neighbors = append(neighbors, cD)
	}

	if cL != nil {
		cL.parent = c
		cL.Depth++
		cL.Cost++
		cL.Move = Left
		neighbors = append(neighbors, cL)
	}

	if cR != nil {
		cR.parent = c
		cR.Depth++
		cR.Cost++
		cR.Move = Right
		neighbors = append(neighbors, cR)
	}

	if c.ReverseNeighbors {
		for i, j := 0, len(neighbors)-1; i < j; i, j = i+1, j-1 {
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		}
	}

	c.neighbors = neighbors

	return neighbors
}

func (c *Config) Key() string {

	if c.key != nil {
		return *c.key
	}

	str := ""

	c.Iterate(func(_, _ *int, val *int8) {
		str += fmt.Sprintf("%d", *val)
	})

	c.key = &str

	return str
}

func (c *Config) Parents() []*Config {
	m := c.parent
	parents := make([]*Config, 0)

	for m != nil {
		parents = append(parents, m)
		m = m.parent
	}

	for i, j := 0, len(parents)-1; i < j; i, j = i+1, j-1 {
		parents[i], parents[j] = parents[j], parents[i]
	}
	return parents
}
