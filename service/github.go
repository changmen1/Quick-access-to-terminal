package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

// 读取 GitHub Token
func getGitHubToken() string {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ 未找到 .env 文件，请确保 .env 文件存在")
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("❌ .env 文件中未找到 GITHUB_TOKEN")
		os.Exit(1)
	}

	return token
}

// Repository 结构体
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

func GithubMain() {
	// 获取 GitHub 仓库列表（包括私有仓库）
	repos, err := fetchGitHubRepositories()
	if err != nil {
		fmt.Println("❌ 获取 GitHub 项目失败:", err)
		os.Exit(1)
	}

	// 让用户选择一个仓库
	selectedRepo, err := selectRepository(repos)
	if err != nil {
		fmt.Println("❌ 选择仓库失败:", err)
		os.Exit(1)
	}

	// 显示选中项目信息
	fmt.Printf("\n✅ 你选择了: %s\n", selectedRepo.Name)
	fmt.Printf("📌 项目地址: %s\n", selectedRepo.URL)

	// 在浏览器打开 GitHub 项目
	openURL(selectedRepo.URL)

	fmt.Println("🚀 进入系统...")
}

// 获取 GitHub 仓库（支持私有仓库）
func fetchGitHubRepositories() ([]Repository, error) {
	url := "https://api.github.com/user/repos?per_page=100"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 获取 GitHub Token
	token := getGitHubToken()

	// 添加身份认证，确保能访问私有仓库
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API 请求失败: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// 选择仓库
func selectRepository(repos []Repository) (Repository, error) {
	var names []string
	for _, repo := range repos {
		names = append(names, repo.Name)
	}

	prompt := promptui.Select{
		Label: "请选择一个项目",
		Items: names,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return Repository{}, err
	}

	return repos[index], nil
}

// 打开 URL
func openURL(url string) {
	var cmd *exec.Cmd

	switch os := os.Getenv("OS"); os {
	case "Windows_NT":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ 打开浏览器失败:", err)
	}
}
