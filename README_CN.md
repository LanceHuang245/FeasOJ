[English](README.md) | 简体中文
<p align="center">
    <a href="https://github.com/LanceHuang245/FeasOJ">
        <img src="images/logo.png" height="200"/>
    </a>
</p>

> FeasOJ [Backend](https://github.com/LanceHuang245/FeasOJ-Backend) 和 [JudgeCore](https://github.com/LanceHuang245/FeasOJ-JudgeCore)现已迁移到该项目中！\
> 现在FeasOJ正在进行全面优化、重构，请等待下一个正式版发布后使用，敬请期待！

# FeasOJ
### 项目简介
FeasOJ 是一个基于 Vue 和 Golang 的在线编程练习平台，支持多国语言、讨论区、竞赛等功能，旨在为用户提供一个方便、高效的学习和练习环境。

### 相关链接
[ImageGuard](https://github.com/LanceHuang245/ImageGuard)\
[Profanity Detector](https://github.com/LanceHuang245/ProfanityDetector)\
[配置文件说明文档](/docs/CONFIG_README_CN.md)\

### 项目特性
- 多语言支持：支持多种语言，包括英语、西班牙语、法语、意大利语、日语、简体中文等
- 多编程语言支持：C++、Java、Python、Rust、PHP、Pascal
- 用户认证：支持用户注册、登录、注销等功能
- 题目管理：支持题目上传、编辑、删除等功能
- 讨论区：支持用户发表、回复、删除评论等功能
- 竞赛：支持创建、参加、结束竞赛等功能
- 代码编辑器：支持 Markdown 编辑器，方便用户编写题解和评论
- 代码高亮：支持代码高亮，方便用户查看和编辑代码。
- 代码提交：支持用户提交代码，并在Docker 容器池中编译运行返回结果
- 实时通知：支持判题结果、竞赛消息的实时通知 (SSE)
- ...

### 项目结构
```
FeasOJ
│ 
├─images
├─docs
├─services       # 后端与JudeCore代码
│  ├─cmd
│  │ ├─app
│  │ │  ├─backend
│  │ │  └─judgecore
│  │ └─pkg
│  ├─go.mod
│  ├─go.sum
│  └─scripts
├─web            # 前端代码
│  ├─public
│  ├─src
│  ├─index.html
│  ├─package-lock.json
│  ├─package.json
│  └─vite.config.js
```

### 环境
- Vue 3
- Golang 1.25.1
- MySQL/PostgreSQL
- Docker
- Redis
- npm
- 最新版本的Chromium系浏览器或Firefox

### 如何运行
1. 克隆此库
2. 在项目目录下执行`./scripts/deps_update.sh` 更新后端与JudgeCore的依赖
3. 确保Docker Desktop(Windows/MacOS)或Docker Engine已启动
4. 在`web`目录下运行 `npm install` 安装依赖项
5. 配置`/web/src/utils/axios.js` 中的 `apiUrl` 为你FeasOJ-Bakcend服务器地址
6. 运行 `npm run dev` 启动前端服务器

### 注意
这是我第一次用Vue + Golang写大项目，所以代码会一坨，不过我会一直去改进它！
如果你找到任何Bug请发布Issue告诉我。

### 本土化
- 阿拉伯语
- 英语
- 西班牙语
- 法语
- 意大利语
- 日语
- 葡萄牙语
- 俄语
- 简体中文
- 繁體中文

如果您想要增加语言翻译或优化当前的翻译，请按照以下步骤：
- [Fork](https://github.com/LanceHuang245/FeasOJ/fork) 该项目仓库
- 复制 `/web/src/plugins/locales/en.js` 以及 `/services/cmd/app/backend/internal/utils/locales/en.json` 文件并以您想要的语言代码作为文件名称，分别粘贴到 `/web/src/plugins/locales`和`/services/cmd/app/backend/internal/utils/locales`，或者直接修改您想优化的文件
- 翻译文件中的内容
- 创建一个 [pull request](https://github.com/LanceHuang245/FeasOJ/pulls) 即可

### 项目截图
![Main](/images/Main.png)
![Login](/images/Login.png)
![Problem](/images/Problem.png)
![Profile](/images/Profile.png)
更多图片可在 [images](/images) 中查看。

### 致谢
- [Vue](https://github.com/vuejs/vue)
- [Vuetify](https://github.com/vuetifyjs/vuetify)
- [vue-avatar-cropper](https://github.com/overtrue/vue-avatar-cropper)
- [vue3-ace-editor](https://github.com/CarterLi/vue3-ace-editor)
- [md-editor-v3](https://github.com/imzbf/md-editor-v3)
- [vue-i18n](https://github.com/intlify/vue-i18n)
- [axios](https://github.com/axios/axios)