English | [简体中文](README_CN.md)
<p align="center">
    <a href="https://github.com/LanceHuang245/FeasOJ">
        <img src="images/logo.png" height="200"/>
    </a>
</p>

> FeasOJ [Backend](https://github.com/LanceHuang245/FeasOJ-Backend) and [JudgeCore](https://github.com/LanceHuang245/FeasOJ-JudgeCore) has been migrated to this project.\
> FeasOJ is now undergoing comprehensive optimization and reconstruction. Please wait for the next stable version to be released and use it. Stay tuned!

# FeasOJ
### Project Description
FeasOJ is an online programming practice platform based on Vue and Golang, supporting multi-languages, discussion forums, contests and other features, aiming to provide users with a convenient and efficient learning and practice environment.

### Related Links
[ImageGuard](https://github.com/LanceHuang245/ImageGuard)\
[Profanity Detector](https://github.com/LanceHuang245/ProfanityDetector)\
[Config Document](/docs/CONFIG_README_EN.md)\

### Features
- Multi-language support: Support multiple languages, including English, Spanish, French, Italian, Japanese, Simplified Chinese etc
- Multi Programming Language Support: C++, Java, Python, Rust, PHP, Pascal
- User Authentication: Supports user registration, login, logout and other functions
- Topic Management: Supports topic uploading, editing, deleting, etc
- Discussion Forum: supports users to post, reply, delete comments and so on
- Contests: support the functions of creating, participating and ending contests
- Code Editor: Supports Markdown editor, which is convenient for users to write questions and comments
- Code Highlighting: Support code highlighting, convenient for users to view and edit code
- Code Submission: Support users to submit code and compile and run in the sandbox to return the result
- Real-time notification: Support real-time notification of question results and contest messages (SSE)

### Project Structure
```
FeasOJ
│ 
├─images
├─docs
├─services       # Back-end and JudgeCore src
│  ├─cmd
│  │ ├─app
│  │ │  ├─backend
│  │ │  └─judgecore
│  │ └─pkg
│  ├─go.mod
│  ├─go.sum
│  └─scripts
├─web            # Front-end src
│  ├─public
│  ├─src
│  ├─index.html
│  ├─package-lock.json
│  ├─package.json
│  └─vite.config.js
```

### Environment
- Vue 3
- Golang 1.25.1
- Docker
- MySQL/PostgreSQL
- Redis
- npm
- The lastest version of Chromium or Firefox

### How to run
1. Clone this repository
2. Run `./scripts/deps_update.sh` under the project directory to update dependencies for the backend and JudgeCore
3. Ensure Docker Desktop (Windows/MacOS) or Docker Engine is started
4. Run `go run main.go` in the `services/cmd/app/backend` and `services/cmd/app/judgecore` directories respectively to start the backend service. Note that you must start the JudgeCore service before the Backend service can be started
5. On the first start of the backend service, `config.toml` will be generated in the directory. You will need to modify it and restart the service
6. Run `npm install` under the `web` directory to install dependencies
7. Configure the `apiUrl` in `/web/src/utils/axios.js` to the address of your FeasOJ-Backend server
8. Run `npm run dev` to start the frontend server

### Notice
This is the first time I've written a big project with Vue + Golang, so the code is going to be terrible, but I'll keep going to improve it!
If you find any bugs, please open an issue.

### Localization
- Arabic
- English
- Espanish
- French
- Italian
- Japanese
- Portuguese
- Russian
- Simplified Chinese
- Traditional Chinese

If you want to contribute adding new language or improving existing language, follow this step:
- [Fork](https://github.com/LanceHuang245/FeasOJ/fork) this repository
- Copy `/web/src/plugins/locales/en.js` and `/services/cmd/app/backend/internal/utils/locales/en.json` into `/web/src/plugins/locales` and `/services/cmd/app/backend/internal/utils/locales` with a new language code as the file name or edit the existing language file
- Translate all the keys in the new language file
- Create a [pull request](https://github.com/LanceHuang245/FeasOJ/pulls)

### Screenshots
![Main](/images/Main.png)
![Problem](/images/Problem.png)
![Profile](/images/Profile.png)
More screenshots can be found in the [images](/images) folder.

### Thanks
- [Vue](https://github.com/vuejs/vue)
- [Vuetify](https://github.com/vuetifyjs/vuetify)
- [vue-avatar-cropper](https://github.com/overtrue/vue-avatar-cropper)
- [vue3-ace-editor](https://github.com/CarterLi/vue3-ace-editor)
- [md-editor-v3](https://github.com/imzbf/md-editor-v3)
- [vue-i18n](https://github.com/intlify/vue-i18n)
- [axios](https://github.com/axios/axios)
