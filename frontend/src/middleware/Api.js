import axios from 'axios';

import Auth from './Auth';

// export let baseURL = "http://localhost:6167/api/v1"
export let baseURL = "/api/v1";

export const getBaseURL = () => {
    return baseURL;
};

const instance = axios.create({
    baseURL: getBaseURL(),
    withCredentials: true,
    crossDomain: true,
});

function AppError(message, code, error) {
    this.code = code;
    this.message = message || '未知错误';
    this.message += error ? (" " + error) : "";
    this.stack = (new Error()).stack;
}

AppError.prototype = Object.create(Error.prototype);
AppError.prototype.constructor = AppError;

instance.interceptors.response.use(
    function (response) {
        response.rawData = response.data;
        response.data = response.data.data;
        if (response.rawData.code !== 200) {
            // 认证错误：设置认证状态为登出，重定向至登录页面
            if (response.rawData.code === 401) {
                Auth.Signout()
                window.location.href = "/#/login";
            }
            // 非管理员，重定向至主页
            if (response.rawData.code === 2000) {
                window.location.href = "/#/home";
            }
            // 错误都要抛出 AppError
            throw new AppError(response.rawData.msg, response.rawData.code, response.rawData.error);
        }
        return response;
    },
    function (error) {
        return Promise.reject(error);
    }
);

export default instance;
