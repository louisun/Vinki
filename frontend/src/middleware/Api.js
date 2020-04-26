import axios from 'axios';

export let baseURL = "api/v1";

export const getBaseURL = () => {
    return baseURL;
};

const instance = axios.create({
    baseURL: getBaseURL(),
    withCredentials: false,
    crossDomain: true,
    // headers: {'Access-Control-Allow-Origin': '*'}
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
        return response;
    },
    function (error) {
        return Promise.reject(error);
    }
);

export default instance;
