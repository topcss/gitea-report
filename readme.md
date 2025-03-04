# 🔥 Gitea项目全景雷达 - 用Excel征服你的代码仓库混乱

**「当你还在用石器时代的方法一个个翻仓库的时候，我们已经把整个Gitea装进了Excel」**

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/topcss/gitea-report/blob/main/LICENSE)
![Go Version](https://img.shields.io/badge/Go-1.18%2B-brightgreen)
[![GitHub stars](https://img.shields.io/github/stars/topcss/gitea-report)](https://github.com/topcss/gitea-report/stargazers)

🚀 专为Gitea打造的军火级仓库审计工具，1分钟生成全景作战地图，让那些失控的分支权限、裸奔的主干分支、幽灵协作账号无所遁形！

## 💥 痛点粉碎机
- **被当猴耍的日子结束了**：还在手动检查每个仓库的last commit时间？我们直接给你按「最后更新时间」排序的战场热力图
- **分支保护？不存在的！**：自动揪出那些没有protected branch的裸奔项目，标记出开了保护但白名单乱给权限的憨憨配置
- **权限黑洞终结者**：5秒扫描全平台协作账号，把「谁在什么项目有上帝权限」给你拍在Excel脸上


## ⚡️ 功能核爆点
```bash
./gitea-report
```
| 功能模块        | 杀伤力                          | 输出示例                  |
|-----------------|---------------------------------|--------------------------|
| 仓库全景扫描    | 500+项目元数据闪电收集          | 创建时间/最后更新精确到秒 |
| 分支安全检测    | 识别裸奔分支/无效保护/权限泄露  | 受保护分支标记+白名单审计 |
| 协作权限矩阵    | 用户-仓库权限关系可视化         | 协作者映射表              |
| 项目活跃热图    | 最近开发动态标记                | 更新时间倒序排序          |

## 🛠️ 安装即战
```bash
# 1. 下载军火
git clone https://github.com/topcss/gitea-report.git

# 2. 进入战斗位置
cd gitea-report

# 3. 执行斩首行动（Win/Mac/Linux通杀）
go run main.go
```

或者，直接下载绿色版，直接运行。


## 🤔 你为什么需要这个？
- **CTO/技术总监**：告别被蒙在鼓里的感觉，全平台项目健康度一表掌控
- **DevOps工程师**：3分钟完成原本需要两周的手动审计
- **安全团队**：权限漏洞检测自动化，合规审计报告生成器
- **接盘侠**：新接手Gitea实例时的生存工具包

## 🌟 技术暴力美学
```go
// 这是我们的战斗宣言
func main() {
    fmt.Println("开始轰炸式数据采集...")
    // 50并发请求+智能缓存
    repos := getAllRepos() 
    
    // 分支保护策略逆向工程
    checkBranchSecurity(repos)  
    
    // 权限矩阵降维打击
    buildPermissionMatrix(repos)  
    
    // 生成Excel核弹头
    launchExcelReport()  
}
```

## 📈 SEO关键词轰炸
#Gitea审计 #仓库安全扫描 #分支权限检测 #DevOps自动化 #项目健康度分析 #代码仓库治理 #权限矩阵可视化 #Gitea数据导出 #分支保护规范 #仓库活跃度分析

---

**⚠️ 警告：使用本工具可能导致你突然看清公司代码管理的混乱现状，请做好心理建设后再运行！**
