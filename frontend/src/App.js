import React from 'react';

import { useSelector } from 'react-redux';
import {
  Route,
  Switch,
  useRouteMatch,
} from 'react-router-dom';

import {
  CssBaseline,
  makeStyles,
  ThemeProvider,
} from '@material-ui/core';
import {
  createMuiTheme,
  lighten,
} from '@material-ui/core/styles';
import useMediaQuery from '@material-ui/core/useMediaQuery';

import Article from './components/Article/Article';
import Login from './components/Login/Login';
import Navbar from './components/Navbar/Navbar';
import Tag from './components/Tag/Tag';
import Auth from './middleware/Auth';
import { changeThemeColor } from './utils';

export default function App() {
  // 从 store 获取 theme 配置
  const themeConfig = useSelector((state) => state.siteConfig.theme);
  // 登录状态
  const isLogin = useSelector((state) => state.viewUpdate.isLogin);
  // 查询 dark 模式
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");
  // 根据 dark 模式设置主题配色
  const theme = React.useMemo(() => {
    themeConfig.palette.type = prefersDarkMode ? "dark" : "light";
    let prefer = Auth.GetPreference("theme_mode");
    if (prefer) {
      themeConfig.palette.type = prefer;
    }
    // 重新构造 theme JSON 结构
    let theme = createMuiTheme({
      ...themeConfig,
      palette: {
        ...themeConfig.palette,
        primary: {
          ...themeConfig.palette.primary,
          main:
            themeConfig.palette.type === "dark"
              ? lighten(themeConfig.palette.primary.main, 0.3)
              : themeConfig.palette.primary.main,
        },
      },
    });
    // 在 main 和 dark 之间切换：设置 meta
    changeThemeColor(
      themeConfig.palette.type === "dark"
        ? theme.palette.background.default
        : theme.palette.primary.main
    );
    return theme;
  }, [prefersDarkMode, themeConfig]);

  // css
  const useStyles = makeStyles((theme) => ({
    // container root
    container: {
      minHeight: "100%",
    },
    root: {
    //   minHeight: "100%",
      minHeight: "calc(100vh - 50px)",
      marginTop: "50px",
      display: "flex",
      justifyContent: "space-between",
      backgroundColor: "#F2F4F8",
    },
  }));

  const classes = useStyles();

  let { path } = useRouteMatch();

  return (
    <React.Fragment>
      <ThemeProvider theme={theme}>
        <div id="container" class={classes.container}>
          {/*基准样式*/}
          <CssBaseline />
          {/*导航栏*/}
          <Navbar />
          {/*主内容*/}
          <div className={classes.root}>
            <Switch>
              <Route path="/login">
                <Login />
              </Route>
              <Route path="/article">
                <Article />
              </Route>
              <Route path="/tags">
                <Tag />
              </Route>
            </Switch>
          </div>
        </div>
      </ThemeProvider>
    </React.Fragment>
  );
}
