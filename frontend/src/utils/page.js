const pageHelper = {
    isHomePage(path){
        return path === "/home"
    },
    isAdminPage(path){
        return path && path.startsWith("/admin")
    },
    isMobile(){
        return window.innerWidth < 600;
    },
}
export default pageHelper