import React from 'react';

import {
  Redirect,
  Route,
  Switch,
} from 'react-router-dom';

import {
  CssBaseline,
  makeStyles,
} from '@material-ui/core';

import Article from './components/Article/Article';
import NotFound from './components/Common/NotFound';
import MessageBar from './components/Common/Snackbar';
import Login from './components/Login/Login';
import Navbar from './components/Navbar/Navbar';
import Tags from './components/Tag/Tag';
import Auth from './middleware/Auth';
import AuthRoute from './middleware/AuthRoute';

const useStyles = makeStyles(() => ({
    container: {
        minHeight: "100%",
    },
    root: {
        minHeight: "calc(100vh - 50px)",
        marginTop: "50px",
        display: "flex",
        justifyContent: "space-between",
        backgroundColor: "#F2F4F8",
        // boxShadow: "0 12px 15px 0 rgba(0,0,0,0.24), 0 17px 50px 0 rgba(0,0,0,0.19)"
    },
}));


export default function App() {
    const classes = useStyles();
    return (
        <React.Fragment>
            <div id="container" className={classes.container}>
                {/*基准样式*/}
                <CssBaseline />
                <MessageBar />
                {/*导航栏*/}
                <Navbar />
                {/*主内容*/}
                <div className={classes.root}>
                    <Switch>
                        <Route path="/login">
                            <Login />
                        </Route>
                        <Route path="/article/:repoName/:tagName/:articleName" render={(props) => (
                            Auth.Check() ?
                                (<Article key={`${props.match.params.repoName}-${props.match.params.tagName}-${props.match.params.articleName}`} />) :
                                (<Redirect
                                    to={{
                                        pathname: "/login",
                                    }}
                                />)
                        )} />
                        <AuthRoute path="/(|home|tags)">
                            <Tags />
                        </AuthRoute>
                        <Route path="*">
                            <NotFound msg={"页面不存在"} />
                        </Route>
                    </Switch>
                </div>
            </div>
        </React.Fragment>
    );
}
