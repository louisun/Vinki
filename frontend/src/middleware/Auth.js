const Auth = {
    isAuthenticated: false,
    // authenticate 登录，在本地 SetUser，并设置 isAuthenticated 为 true
    authenticate(cb) {
        Auth.SetUser(cb);
        Auth.isAuthenticated = true;
    },
    // GetUser 在本地缓存中获取 user 的 JSON 信息
    GetUser(){
        return JSON.parse(localStorage.getItem("user"))
    },
    // SetUser 在本地缓存中添加 user 的 JSON 信息
    SetUser(newUser){
        localStorage.setItem("user", JSON.stringify(newUser));
    },
    // Check 检查是否已认证
    Check() {
        if (Auth.isAuthenticated) {
            return true;
        }
        if (localStorage.getItem("user") !== null){
            return !Auth.GetUser().anonymous;
        }
        return false

    },
    // signout 登出
    signout() {
        Auth.isAuthenticated = false;
        let oldUser = Auth.GetUser();
        oldUser.id = 0;
        localStorage.setItem("user", JSON.stringify(oldUser));
    },
    // SetPreference 在本地保存用户个性化设置的 key-value
    SetPreference(key,value){
        let preference = JSON.parse(localStorage.getItem("user_preference"));
        preference = (preference == null) ? {} : preference;
        preference[key] = value;
        localStorage.setItem("user_preference", JSON.stringify(preference));
    },
    // GetPreference 获取个性化设置中 key 对应的 value
    GetPreference(key){
        let preference = JSON.parse(localStorage.getItem("user_preference"));
        if (preference && preference[key]){
            return preference[key];
        }
        return null;
    },
};

export default Auth;
