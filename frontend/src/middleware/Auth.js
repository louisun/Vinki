const Auth = {
    isAuthenticated: false,
    authenticate(user) {
        Auth.SetUser(user);
        Auth.isAuthenticated = true;
    },
    GetUser() {
        return JSON.parse(localStorage.getItem("user"));
    },
    SetUser(user) {
        localStorage.setItem("user", JSON.stringify(user))
    },
    Check() {
        if (Auth.isAuthenticated) {
            return true;
        }
        // 缓存中有user信息，暂且视为已登录
        if (localStorage.getItem("user") !== null) {
            return Auth.GetUser().id !== 0;
        }
        return false
    },
    Signout() {
        Auth.isAuthenticated = false;
        let oldUser = Auth.GetUser();
        oldUser.id = 0;
        localStorage.setItem("user", JSON.stringify(oldUser));
    }
}

export default Auth;