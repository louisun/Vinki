import {
  changeViewMethod,
  enableLoadUploader,
  setSiteConfig,
  toggleSnackbar,
} from '../actions/index';
import { fixUrlHash } from '../utils/index';
import pathHelper from '../utils/page';
import API from './Api';
import Auth from './Auth';

// 初始化站点配置：用户个性化配置
export var InitSiteConfig = (rawStore) => {
    // 从「缓存」加载默认配置
    let configCache = JSON.parse(localStorage.getItem('siteConfigCache'));
    if (configCache != null) {
        rawStore.siteConfig = configCache
    }
    // 检查是否有path参数
    var url = new URL(fixUrlHash(window.location.href));
    var c = url.searchParams.get("path");
    rawStore.navigator.path = c===null?"/":c;
    // 初始化用户个性配置（主题）
    rawStore.siteConfig = initUserConfig(rawStore.siteConfig);
    // 是否登录
    rawStore.viewUpdate.isLogin = Auth.Check();

    // 更改站点标题
    document.title = rawStore.siteConfig.title;
    return rawStore
}

const initUserConfig = (siteConfig) => {
    if (siteConfig.user!==undefined && !siteConfig.user.anonymous){
        let themes = JSON.parse(siteConfig.themes);
        let user = siteConfig.user;
        delete siteConfig.user
    
        // 更换全局 theme 为用户自定义的配色
        if (user["preferred_theme"] !== "" && themes[user["preferred_theme"]] !== undefined){
            siteConfig.theme = themes[user["preferred_theme"]]
        }

        // 设置登录状态
        Auth.authenticate(user);
    }
    // 本地保存用户信息
    if (siteConfig.user!==undefined && siteConfig.user.anonymous){
        Auth.SetUser(siteConfig.user);
    }
    return siteConfig
}

export function enableUploaderLoad(){
    // 开启上传组件加载
    let user = Auth.GetUser();
    window.policyType = user!==null?user.policy.saveType : "local";
    window.uploadConfig = user!==null?user.policy:{};
    window.pathCache = [];
}

// UpdateSiteConfig 更新站点设置
export async function UpdateSiteConfig(store) {
    // 从后端获取最新配置，并更新本地数据
    API.get("/site/config").then(function(response) {
        let themes = JSON.parse(response.data.themes);
        response.data.theme = themes[response.data.defaultTheme]
        // 用户主题覆盖
        response.data = initUserConfig(response.data)
        store.dispatch(setSiteConfig(response.data));
        localStorage.setItem('siteConfigCache', JSON.stringify(response.data));

        // 偏爱的列表样式
        let preferListMethod = Auth.GetPreference("view_method");
        if(preferListMethod){
            store.dispatch(changeViewMethod(preferListMethod));
        }else{
            let path = window.location.hash.split("#");
            if(path.length >=1 && pathHelper.isSharePage(path[1])){
                store.dispatch(changeViewMethod(response.data.share_view_method));
            }else{
                store.dispatch(changeViewMethod(response.data.home_view_method));
            }
        }

    }).catch(function(error) {
        store.dispatch(toggleSnackbar("top", "right", "无法加载站点配置：" + error.message, "error"));
    }).then(function () {
        enableUploaderLoad(store);
        store.dispatch(enableLoadUploader())
    });
}