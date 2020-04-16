import React from 'react';
import ReactDOM from 'react-dom';

import { Provider } from 'react-redux';
import {
  HashRouter as Router,
  Route,
  Switch,
} from 'react-router-dom';
import { createStore } from 'redux';
import './assets/css/index.css'

import App from './App';
import {
  InitSiteConfig,
  UpdateSiteConfig,
} from './middleware/Init';
import vinkiApp from './reducers';
import * as serviceWorker from './serviceWorker';

// make app work offline and load faster
serviceWorker.register();

// 初始化站点配置
const defaultStatus = InitSiteConfig(
  // rawStore
  { siteConfig: {
      title: "Vinki",
      loginCaptcha: false,
      regCaptcha: false,
      forgetCaptcha: false,
      emailActive: false,
      QQLogin: false,
      themes:null,
      authn:false,
      // 主题
      theme: {
          palette: {
              common: { black: "#000", white: "#fff" },
              background: { paper: "#fff", default: "#fafafa" },
              primary: {
                  light: "#7986cb",
                  // main: "#3F51B6",
                  main: "#4E5668",
                  dark: "#2E3441",
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
              explorer: {
                  filename: "#474849",
                  icon: "#8f8f8f",
                  bgSelected: "#D5DAF0",
                  emptyIcon: "#e8e8e8"
              }
          }
      }
  },
  navigator: {
      path: "/",
      refresh: true
  },
  viewUpdate: {
      isLogin:false,
      loadUploader:false,
      open: false,
      explorerViewMethod: "icon",
      sortMethod: "timePos",
      subTitle:null,
      contextType: "none",
      menuOpen: false,
      navigatorLoading: true,
      navigatorError: false,
      navigatorErrorMsg: null,
      modalsLoading: false,
      storageRefresh: false,
      userPopoverAnchorEl: null,
      shareUserPopoverAnchorEl: null,
      modals: {
          createNewFolder: false,
          rename: false,
          move: false,
          remove: false,
          share: false,
          music: false,
          remoteDownload: false,
          torrentDownload: false,
          getSource: false,
          copy:false,
          resave: false,
          compress:false,
          decompress:false,
      },
      snackbar: {
          toggle: false,
          vertical: "top",
          horizontal: "center",
          msg: "",
          color: ""
      }
  },
  explorer: {
      dndSignal:false,
      dndTarget:null,
      dndSource:null,
      fileList: [],
      dirList: [],
      selected: [],
      selectProps: {
          isMultiple: false,
          withFolder: false,
          withFile: false
      },
      imgPreview: {
          first: null,
          other: []
      },
      keywords: null
  }
});
// redux Store 对象
let store = createStore(vinkiApp, defaultStatus,window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());
// 从后端获取最新配置，并更新本地数据，更新 Store 对象
UpdateSiteConfig(store);



ReactDOM.render(
  <Provider store={store}>
    <Router>
      <Switch>
        {/* <Route path="/admin">
          <Suspense fallback={"Loading..."}>
            <Admin />
          </Suspense>
        </Route> */}
        <Route exact path="">
          <App />
        </Route>
      </Switch>
    </Router>
    </Provider>,
  document.getElementById("root")
);
