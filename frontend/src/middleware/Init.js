import API from '../middleware/Api';
import Auth from './Auth';

var config = {
    // 0. 站点设置
    siteConfig: {
        title: "Vinki",
        theme: {
            common: { black: "#000", white: "#fff" },
            background: { paper: "#fff", default: "#fafafa" },
            primary: {
                light: "#7986cb",
                main: "#4E5668",
                dark: "#303f9f",
                contrastText: "#fff"
            },
            secondary: {
                light: "#ff4081",
                main: "#f50057",
                dark: "#c51162",
                contrastText: "#fff"
            },
            error: {
                light: "#e57373",
                main: "#f44336",
                dark: "#d32f2f",
                contrastText: "#fff"
            },
            text: {
                primary: "rgba(0, 0, 0, 0.87)",
                secondary: "rgba(0, 0, 0, 0.54)",
                disabled: "rgba(0, 0, 0, 0.38)",
                hint: "rgba(0, 0, 0, 0.38)"
            },
        },
    },
    // 1. 仓库信息
    repo: {
        repos: [],           // 仓库列表
        currentRepo: ""      // 当前仓库信息
    },
    // 2. 标签相关信息
    tag: {
        currentTag: "",      // 当前标签信息
        currentTopTag: "",   // 当前一级标签
        topTags: [],         // 一级标签列表
        secondTags: [],      // 二级标签列表
        subTags: [],         // 子标签
        articleList: [],     // 文章信息列表
    },
    snackbar: {
        toggle: false,
        vertical: "top",
        horizontal: "center",
        msg: "",
        color: ""
    },
    isLogin: false
}

export function InitConfig() {
    // 先要去本地缓存中查找相应信息
    let user = Auth.GetUser()
    if (user !== null && user.id !== 0) {
        Auth.authenticate(user)
    }
    // 初始化全局登录状态（与Auth一致）
    config.isLogin = Auth.Check()
    return config
}

export function UpdateConfig(store) {
    API.get("/site/config").then(response => {
        if (response.data !== undefined) {
            // 更新登录状态
            Auth.authenticate(response.data)
        }
    }).catch(error => { })
}