package maze

import (
	"fmt"
	"testing"
)

func TestMaze(t *testing.T)  {
	//注意有坑，读取数据时，Windows系统每个换行符,都会被解析为0，导致迷宫数据错乱
	maze := readMaze("./maze.in")
	for _,row := range maze{
		for _,val:=range row{
			fmt.Printf("%3d",val)
		}
		fmt.Println()
	}
	fmt.Println()
	steps := walk(maze,point{0,0},point{len(maze)-1,len(maze[0])-1})
	for _,row :=range steps{
		for _,val :=range row{
			fmt.Printf("%3d",val)
		}
		fmt.Println()
	}
}