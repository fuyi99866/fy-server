package maze

import (
	"fmt"
	"os"
)

/**
使用广度优先算法，搜索最短路径
 */

var dirs = [4]point{
	{-1,0},
	{1,0},
	{0,-1},
	{0,1},
}

type point struct {
	i,j int
}

func readMaze(filename string) [][] int {
	file,err :=os.Open(filename)
	if err != nil{
		panic(err)
	}

	var row,col int
	fmt.Fscanf(file,"%d %d",&row, &col)
	maze :=make([][]int,row)
	for i :=range maze{
		maze[i] = make([]int,col)
		for j :=range maze[i]{
			fmt.Fscanf(file,"%d",&maze[i][j])
		}
	}
	return maze
}

func (p point)add(r point) point  {
	return point{p.i+r.i,p.j+r.j}
}

func (p point) at(grid [][]int) (int ,bool){
	if p.i<0||p.i>=len(grid){
		return 0,false
	}
	if p.j<0||p.j>=len(grid[p.i]) {
		return 0,false
	}

	return grid[p.i][p.j],true
}

func walk(maze [][]int,start,end point) [][]int {
	steps := make([][]int ,len(maze))
	for i:=range steps{
		steps[i] = make([]int ,len(maze[i]))
	}
	Q :=[]point{start}
	for len(Q)>0 {
		cur := Q[0]
		Q = Q[1:]
		if cur == end{
			break
		}

		for _,dir :=range dirs{
			next := cur.add(dir)
			val ,ok := next.at(maze)
			if !ok || val == 1{
				continue
			}

			//fmt.Println("*****  ",k,dir.i,dir.j,val)

			val,ok = next.at(steps)
			if !ok||val != 0{
				continue
			}
			if next ==start{
				continue
			}
			curSteps,_:=cur.at(steps)
			steps[next.i][next.j] = curSteps +1
			Q = append(Q,next)
		}

	}
	return steps
}

