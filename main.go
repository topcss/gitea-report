package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Gitea 配置
var (
	GITEA_BASE_URL string
	API_TOKEN      string
)

func main() {
	fmt.Print("请输入 Gitea 实例地址（例如：http://10.0.0.10/api/v1）: ")
	fmt.Scanln(&GITEA_BASE_URL)
	fmt.Print("请输入您的 API Token: ")
	fmt.Scanln(&API_TOKEN)

	if GITEA_BASE_URL == "" {
		GITEA_BASE_URL = "http://10.0.0.10/api/v1"
	}
	if API_TOKEN == "" {
		API_TOKEN = "c142a2fa4b1bdd1f5fac76d9336c982d1337dcf5"
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	// 创建工作表
	repoSheet := "仓库信息"
	branchSheet := "分支信息"
	collabSheet := "协作者信息"

	f.NewSheet(repoSheet)
	f.NewSheet(branchSheet)
	f.NewSheet(collabSheet)
	f.DeleteSheet("Sheet1")

	// 设置表头
	setHeaders(f, repoSheet, []string{"仓库全名", "拥有者", "创建时间", "最后更新时间", "仓库描述"})
	setHeaders(f, branchSheet, []string{"所属仓库", "分支名称", "最新提交", "受保护", "合并白名单用户"})
	setHeaders(f, collabSheet, []string{"所属仓库", "协作者"})

	// 收集数据并写入
	repos := getAllRepos()
	// 设置行号
	repoRow, branchRow, collabRow := 2, 2, 2

	// 遍历仓库数据
	for _, repo := range repos {
		fullName := repo["full_name"].(string)
		owner := repo["owner"].(map[string]interface{})["login"].(string)
		created := formatTime(repo["created_at"])
		updated := formatTime(repo["updated_at"])
		description := ""
		if desc, ok := repo["description"].(string); ok {
			description = desc
		}

		// 写入仓库信息
		f.SetCellValue(repoSheet, fmt.Sprintf("A%d", repoRow), fullName)
		f.SetCellValue(repoSheet, fmt.Sprintf("B%d", repoRow), owner)
		f.SetCellValue(repoSheet, fmt.Sprintf("C%d", repoRow), created)
		f.SetCellValue(repoSheet, fmt.Sprintf("D%d", repoRow), updated)
		f.SetCellValue(repoSheet, fmt.Sprintf("E%d", repoRow), description)
		repoRow++

		// 处理分支信息
		branches := getBranches(fullName)
		// 修改分支信息处理部分
		for _, branch := range branches {
			branchName := branch["name"].(string)
			commitID := branch["commit"].(map[string]interface{})["id"].(string)

			protected := "否"
			whitelist := "" // 新增白名单变量
			if protections := getBranchProtections(fullName); len(protections) > 0 {
				for _, p := range protections {
					if p["branch_name"].(string) == branchName {
						protected = "是"
						// 获取白名单用户
						if users, ok := p["merge_whitelist_usernames"].([]interface{}); ok {
							var names []string
							for _, u := range users {
								names = append(names, u.(string))
							}
							whitelist = strings.Join(names, ", ")
						}
						break
					}
				}
			}

			f.SetCellValue(branchSheet, fmt.Sprintf("A%d", branchRow), fullName)
			f.SetCellValue(branchSheet, fmt.Sprintf("B%d", branchRow), branchName)
			f.SetCellValue(branchSheet, fmt.Sprintf("C%d", branchRow), commitID)
			f.SetCellValue(branchSheet, fmt.Sprintf("D%d", branchRow), protected)
			f.SetCellValue(branchSheet, fmt.Sprintf("E%d", branchRow), whitelist)
			branchRow++
		}

		// 处理协作者信息
		collaborators := getCollaborators(fullName)
		for _, collab := range collaborators {
			f.SetCellValue(collabSheet, fmt.Sprintf("A%d", collabRow), fullName)
			f.SetCellValue(collabSheet, fmt.Sprintf("B%d", collabRow), collab["login"].(string))
			collabRow++
		}
	}

	// 自动调整列宽
	setAutoWidth(f, repoSheet, []string{"A", "B", "C", "D", "E"})
	setAutoWidth(f, branchSheet, []string{"A", "B", "C", "D", "E"})
	setAutoWidth(f, collabSheet, []string{"A", "B"})

	// 保存文件
	fileName := fmt.Sprintf("gitea-report-%s.xlsx", time.Now().Format("20060102-150405"))
	if err := f.SaveAs(fileName); err != nil {
		fmt.Printf("保存文件失败: %v\n", err)
		return
	}

	// 获取完整保存路径
	absPath, _ := filepath.Abs(fileName)
	fmt.Printf("\n报表已生成: \n%s\n", absPath)

	// 避免闪退
	fmt.Print("按回车键退出...")
	fmt.Scanln()
}

// 通用表头设置函数
func setHeaders(f *excelize.File, sheet string, headers []string) {
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i) // 'A' 的ASCII码是65，当i=1时：'A' + 1 = 66 → 对应字符'B' 生成字符串"B1"
		f.SetCellValue(sheet, cell, h)
	}
}

// 通用自动调整列宽函数
func setAutoWidth(f *excelize.File, sheet string, cols []string) {
	for _, col := range cols {
		width, _ := f.GetColWidth(sheet, col)
		if width < 20 {
			_ = f.SetColWidth(sheet, col, col, 20)
		}
	}
}

// 获取所有仓库信息
func getAllRepos() []map[string]interface{} {
	var allRepos []map[string]interface{}
	page := 1
	limit := 50

	for {
		url := fmt.Sprintf("%s/repos/search?page=%d&limit=%d", GITEA_BASE_URL, page, limit)
		resp, err := makeRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("第%d页请求失败: %v\n", page, err)
			break
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("第%d页JSON解析失败: %v\n", page, err)
			break
		}

		dataSlice, ok := result["data"].([]interface{})
		if !ok || len(dataSlice) == 0 {
			break
		}

		for _, item := range dataSlice {
			if repo, ok := item.(map[string]interface{}); ok {
				allRepos = append(allRepos, repo)
			}
		}

		if len(dataSlice) < limit {
			break
		}

		page++
	}

	return allRepos
}

// 获取仓库分支信息
func getBranches(fullName string) []map[string]interface{} {
	url := fmt.Sprintf("%s/repos/%s/branches", GITEA_BASE_URL, fullName)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("获取分支列表失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var branches []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&branches)
	return branches
}

// 获取仓库协作者信息
func getCollaborators(fullName string) []map[string]interface{} {
	url := fmt.Sprintf("%s/repos/%s/collaborators", GITEA_BASE_URL, fullName)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("获取协作者列表失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var collaborators []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&collaborators)
	return collaborators
}

// 获取分支保护信息
func getBranchProtections(fullName string) []map[string]interface{} {
	url := fmt.Sprintf("%s/repos/%s/branch_protections", GITEA_BASE_URL, fullName)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("获取分支保护信息失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var protections []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&protections); err != nil {
		fmt.Printf("解析分支保护信息失败: %v\n", err)
		return nil
	}
	return protections
}

// 格式化时间
func formatTime(t interface{}) string {
	timeStr, ok := t.(string)
	if !ok {
		return "时间信息无效"
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "时间格式错误"
	}

	return parsedTime.Local().Format("2006-01-02 15:04:05")
}

// 发起HTTP请求
func makeRequest(method, url string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", API_TOKEN))
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
