package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

const dataFile = "todos.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println(`  go run main.go add "学习 Go"`)
		fmt.Println("  go run main.go list")
		fmt.Println("  go run main.go done 1")
		fmt.Println("  go run main.go delete 1")
		return
	}

	todos, err := loadTodos()
	if err != nil {
		fmt.Println("读取数据失败", err)
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("请提供代办内容")
			return
		}
		title := os.Args[2]
		todos = addTodo(todos, title)
		if err := saveTodos(todos); err != nil {
			fmt.Println("保存失败", err)
			return
		}
		fmt.Println("添加成功")

	case "list":
		listTodos(todos)

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("请输入要完成的ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID必须是数字")
			return
		}
		ok := doneTodo(todos, id)
		if !ok {
			fmt.Println("没有找到这个代办")
			return
		}
		if err := saveTodos(todos); err != nil {
			fmt.Println("保存失败")
			return
		}
		fmt.Println("已标记完成")

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("请提供要删除的代办ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID必须是数字")
			return
		}
		newtodos := deleteTodo(todos, id)
		if len(newtodos) == len(todos) {
			fmt.Println("没有找到这个代办")
			return
		}
		if err := saveTodos(newtodos); err != nil {
			fmt.Println("保存失败", err)
			return
		}
		fmt.Println("删除成功")
	case "clear":
		clearedTodo := clearTodos(todos)
		if len(clearedTodo) == len(todos) {
			fmt.Println("没有要清理的已完成的任务")
			return
		}
		if err := saveTodos(clearedTodo); err != nil {
			fmt.Println("保存失败", err)
			return
		}
		fmt.Println("已清理所有已完成的任务")
	default:
		fmt.Println("未知指令")
	}

}

func loadTodos() ([]todo, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []todo{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []todo{}, nil
	}

	var todos []todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}
	return todos, nil

}

func saveTodos(todos []todo) error {
	data, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func nextID(todos []todo) int {
	maxID := 0
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	return maxID + 1
}

func addTodo(todos []todo, title string) []todo {
	todo := todo{
		ID:    nextID(todos),
		Title: title,
		Done:  false,
	}
	return append(todos, todo)
}

func listTodos(todos []todo) {
	if len(todos) == 0 {
		fmt.Println("暂无代办")
		return
	}

	for _, todo := range todos {
		status := "未完成"
		if todo.Done {
			status = "已完成"
		}
		fmt.Printf("%d %s [%s]\n", todo.ID, todo.Title, status)
	}
}

func doneTodo(todos []todo, id int) bool {
	for i := range todos {
		if todos[i].ID == id {
			if todos[i].Done {
				fmt.Println("该任务之前完成啦")
				return true
			}
			todos[i].Done = true
			return true
		}
	}
	return false
}

func deleteTodo(todos []todo, id int) []todo {
	for i := range todos {
		if todos[i].ID == id {
			return append(todos[:i], todos[i+1:]...)
		}
	}
	return todos
}

func clearTodos(todos []todo) []todo {
	var newtodos []todo
	for _, Todo := range todos {
		if !Todo.Done {
			newtodos = append(newtodos, Todo)
		}
	}
	return newtodos
}
