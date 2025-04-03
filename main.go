package main

import (
	"fmt"
	"os"
	github "zxl-boos/service"
)

func main() {
	data, err := os.ReadFile("./config/logo.txt")
	if err != nil {
		fmt.Println("读取logo失败:", err)
		return
	}
	fmt.Println(string(data))
	github.GithubMain() // 调用 `github` 包的 `GithubMain` 方法
	fmt.Println("按回车键退出...")
	fmt.Scanln() // 等待用户输入，防止窗口自动关闭
}
