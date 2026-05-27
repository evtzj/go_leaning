**Todo CLI**

这是一个用 Go 实现的简单命令行待办清单工具，用于练习 Go 基本语法、文件读写和命令行参数解析。

**项目文件**
- **主程序**: [go_leaning/todo-cli/main.go](go_leaning/todo-cli/main.go#L1-L400)
- **数据文件**: todos.json（运行时在同级目录生成）

**快速开始**
- **运行（临时）**: 使用 `go run` 快速执行：

```bash
go run main.go add "学习 Go"
go run main.go list
go run main.go done 1
go run main.go delete 1
```

- **构建二进制并运行**:

```bash
go build -o todo
./todo add "学习 Go"
./todo list
```

**支持的命令**
- **add <标题>**: 添加一条待办，例如 `add "学习 Go"`。
- **list**: 列出所有待办并显示完成状态。
- **done <ID>**: 将指定 ID 标记为已完成。
- **unDone <ID>**: 取消完成状态（将已完成标记为未完成）。
- **delete <ID>**: 删除指定 ID 的待办。
- **clear**: 清理所有已完成的待办。
- **update <ID> <新标题>**: 更新指定待办的标题。
- **listDone**: 仅列出已完成的待办。
- **listUndone**: 仅列出未完成的待办。
- **stats**: 显示任务统计（总数、已完成、未完成）。

示例：

```bash
go run main.go add "买书"
go run main.go list
go run main.go done 2
go run main.go update 2 "买 Go 书"
go run main.go delete 3
```

**数据存储**
- 程序把所有待办保存在当前目录下的 `todos.json`，以 JSON 格式序列化。第一次运行若无该文件会自动创建。
- 数据结构（简要）：每条待办包含 ID、Title、Done 字段。

**实现要点**
- 使用 `os.Args` 解析命令行参数，使用 `encoding/json` 进行序列化/反序列化。
- ID 由当前最大 ID + 1 生成，删除不会重用旧 ID，保证 ID 唯一性。
- 每次对数据的增删改操作后会调用保存函数持久化到 `todos.json`。

**注意事项 & 建议**
- 当前实现为单文件、同步读写，适合学习与个人使用；若多进程同时访问需加锁或升级为数据库存储。
- 标题参数使用简单的 `os.Args` 读取：如果标题含有空格，请用引号包裹（示例已给出）。
- 若想改进：可加入子命令解析库（例如 cobra）、增加交互式界面或添加单元测试。

**贡献与扩展**
- 欢迎自行拓展功能：例如按优先级排序、支持到期时间、保存到用户目录（例如 `$HOME/.config/todo-cli/`）等。
- 如果你想我帮你把某个改进实现到代码里，告诉我具体需求，我可以直接修改代码并验证。

**许可证**
- 你可以按个人练习用途自由使用和修改本代码。
