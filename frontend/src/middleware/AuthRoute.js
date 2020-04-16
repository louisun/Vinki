import React from 'react';

import {
  Redirect,
  Route,
} from 'react-router-dom';

import Auth from './Auth';

function AuthRoute({ children, ...rest }) {
    return (
      <Route
        {...rest}
        render={({ location }) =>
        Auth.Check(rest.isLogin) ? (
            children
          ) : (
            <Redirect
              to={{
                pathname: "/login",
                state: { from: location }
              }}
            />
          )
        }
      />
    );
  }

export default AuthRoute