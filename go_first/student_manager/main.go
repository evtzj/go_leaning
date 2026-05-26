package main

import "fmt"

type student struct {
	id    int
	name  string
	score float64
}

var students []student

func main() {
	for {
		showMenu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			addStudent()
		case 2:
			viewStudents()
		case 3:
			searchStudent()
		case 4:
			deleteStudent()
		case 5:
			fmt.Println("退出系统")
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}

}

func showMenu() {
	fmt.Println("========== 学生成绩管理系统 ==========")
	fmt.Println("1. 添加学生")
	fmt.Println("2. 查看所有学生")
	fmt.Println("3. 查询学生")
	fmt.Println("4. 删除学生")
	fmt.Println("5. 退出系统")
	fmt.Println("====================================")
}
func addStudent() {
	var s student
	fmt.Print("请输入学生ID: ")
	fmt.Scanln(&s.id)
	for _, student := range students {
		if student.id == s.id {
			fmt.Println("已经存在这个学生")
		}
	}
	fmt.Print("请输入学生姓名:")
	fmt.Scanln(&s.name)
	fmt.Print("请输入学生成绩:")
	fmt.Scanln(&s.score)
	students = append(students, s)
	fmt.Println("学生添加成功！")

}

func viewStudents() {
	if len(students) == 0 {
		fmt.Println("没有学生信息")
		return
	}
	for _, student := range students {
		fmt.Printf("学生学号: %d,学生姓名: %s,学生成绩: %.2f\n", student.id, student.name, student.score)
	}
}

func searchStudent() {
	var temp int
	fmt.Println("请输入你要查询的用户id")
	fmt.Scanln(&temp)
	for _, student := range students {
		if student.id == temp {
			fmt.Println("查询成功")
			fmt.Printf("学生学号: %d,学生姓名: %s,学生成绩: %.2f\n", student.id, student.name, student.score)
			return
		}
	}
	fmt.Println("查询失败,不存在该学生")
}

func deleteStudent() {
	var temp int
	fmt.Println("请输入你要删除的学生id")
	fmt.Scanln(&temp)
	for index, student := range students {
		if student.id == temp {
			students = append(students[:index], students[index+1:]...)
			fmt.Println("删除成功")
			return
		}
	}
	fmt.Println("删除失败")
}
