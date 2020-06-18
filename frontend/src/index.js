import './assets/css/index.css';

// import './assets/js/prism';
// import 'nord-highlightjs';
import React from 'react';
import ReactDOM from 'react-dom';

import { Provider } from 'react-redux';
import {
  HashRouter as Router,
  Route,
  Switch,
} from 'react-router-dom';
import { createStore } from 'redux';

import App from './App';
import {
  InitConfig,
  UpdateConfig,
} from './middleware/Init';
import vinkiApp from './reducers';
import * as serviceWorker from './serviceWorker';

// make app work offline and load faster
serviceWorker.register();

// 初始化配置
const defaultConfig = InitConfig()

let store = createStore(
  vinkiApp,
  defaultConfig,
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

// 从服务端更新配置，并保存到全局状态中
UpdateConfig(store)

ReactDOM.render(
  <Provider store={store}>
    <Router>
      <Switch>
        <Route exact path="">
          <App />
        </Route>
      </Switch>
    </Router>
  </Provider>,
  document.getElementById("root")
);
