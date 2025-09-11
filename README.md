English | [简体中文](README_CN.md)
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

# FeasOJ
### Project Description
FeasOJ is an online programming practice platform based on Vue and Golang, supporting multi-languages, discussion forums, contests and other features, aiming to provide users with a convenient and efficient learning and practice environment.
<br>
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)
[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
[API Document(CN)](https://claret-feasoj.apifox.cn)
[API Document(EN)](https://claret-feasoj.apifox.cn/en/)

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
├─public
├─src
│  ├─assets
│  ├─components
│  ├─pages
│  ├─plugins
│  ├─router
│  └─utils
```

### Environment
- Vue 3
- pnpm
- The lastest version of Chromium or Firefox

### How to run
1. Clone repository.
2. Run `pnpm i` Install dependencies.
3. Config `src/utils/axios.js` to your [FeasOJ-Backend](https://github.com/ClaretWheel1481/FeasOJ-Backend) address.
4. Run `pnpm dev` to start the front-end server.

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
- [Fork](https://github.com/ClaretWheel1481/FeasOJ/fork) this repository
- Copy `src/plugins/locales/en.js` file into `src/plugins/locales` with a new language code as the file name or edit the existing language file
- Translate all the keys in the new language file
- Create a [pull request](https://github.com/ClaretWheel1481/FeasOJ/pulls)

### Screenshots
![Main](/assets/Main.png)
![Login](/assets/Login.png)
![Problem](/assets/Problem.png)
![Profile](/assets/Profile.png)
More screenshots can be found in the [assets](/assets) folder.

### Thanks
- [Vue](https://github.com/vuejs/vue)
- [Vuetify](https://github.com/vuetifyjs/vuetify)
- [vue-avatar-cropper](https://github.com/overtrue/vue-avatar-cropper)
- [vue3-ace-editor](https://github.com/CarterLi/vue3-ace-editor)
- [md-editor-v3](https://github.com/imzbf/md-editor-v3)
- [vue-i18n](https://github.com/intlify/vue-i18n)
- [axios](https://github.com/axios/axios)
