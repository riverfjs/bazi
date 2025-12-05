# 🎋 八字排盘系统

基于 Go 语言开发的在线八字排盘计算器，提供美观的 Web 界面和完整的八字计算功能。

## ✨ 功能特点

- 📅 完整的公历/农历日期转换
- 🔮 精确的四柱八字计算
- 🌟 大运起运时间推算
- 👤 支持性别和出生地信息
- 📋 一键复制格式化结果
- 🎨 现代化的响应式界面
- ⚡ 快速部署到云平台

## 🚀 快速开始

### 本地运行

```bash
# 克隆项目
git clone <your-repo-url>
cd bazi

# 安装依赖
go mod download

# 运行服务器
go run main.go

# 访问
打开浏览器访问 http://localhost:8080
```

### Docker 运行

```bash
# 构建镜像
docker build -t bazi-paipan .

# 运行容器
docker run -p 8080:8080 bazi-paipan
```

## 🌐 免费部署

支持一键部署到以下平台（完全免费）：

### 推荐：Railway（最简单）
1. 访问 https://railway.app
2. 用 GitHub 登录
3. 选择此仓库，自动部署

### Render
1. 访问 https://render.com
2. 连接 GitHub 仓库
3. 选择 Web Service，自动配置

### Fly.io
```bash
flyctl launch
flyctl deploy
```

详细部署教程请查看 [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) 或 [QUICK_DEPLOY.md](./QUICK_DEPLOY.md)

## 📖 API 文档

### POST /api/bazi

计算八字信息

**请求参数：**
```json
{
  "year": 2000,
  "month": 1,
  "day": 1,
  "hour": 12,
  "minute": 0,
  "second": 0,
  "sex": 1
}
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "solarDate": "2000年1月1日 12:0:0",
    "lunarDate": "农历己卯年十一月廿五",
    "siZhu": "己卯 丙子 甲寅 庚午",
    "daYun": "...",
    "qiYunDate": "..."
  }
}
```

## 🎨 界面预览

- 简洁优雅的表单设计
- 实时验证输入
- 动画效果提升体验
- 响应式布局支持移动端

## 📦 项目结构

```
.
├── main.go                # 主程序
├── public/                # 静态文件
│   └── index.html         # Web界面
├── go.mod                 # Go依赖
├── Dockerfile             # Docker配置
├── .gitignore             # Git忽略文件
├── railway.json           # Railway配置
├── render.yaml            # Render配置
├── fly.toml               # Fly.io配置
├── DEPLOYMENT_GUIDE.md    # 详细部署指南
└── QUICK_DEPLOY.md        # 快速部署指南
```

## 🔧 技术栈

- **后端：** Go 1.21+
- **前端：** 原生 HTML/CSS/JavaScript
- **八字库：** github.com/warrially/BaziGo

## 📄 版权声明

本项目基于 BaziGo 八字库开发
- 原作者库地址：https://github.com/warrially/BaziGo
- 作者联系方式：+86-167-632-33049

作者只想保留版权，无任何使用或者发布限制。您只需要在您的发行版本中注明代码出处

## 🙏 致谢

八字部分参考的是三清宫命理
- https://weibo.com/bazishequ

日历部分参考：
- 中国日历类（Chinese Calendar Class (CCC)）v0.1
- 版权所有 (C) 2002-2003 neweroica (wy25@mail.bnu.edu.cn)
- CNPACK 作者：刘啸 (liuxiao@cnpack.org)、周劲羽(zjy@cnpack.org)

