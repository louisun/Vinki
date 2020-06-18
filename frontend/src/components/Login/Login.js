import React, {
  useCallback,
  useState,
} from 'react';

import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';

import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import {
  makeStyles,
  withStyles,
} from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';

import {
  setLoginStatus,
  toggleSnackbar,
} from '../../actions';
import API from '../../middleware/Api';
import Auth from '../../middleware/Auth';

const CssTextField = withStyles({
  root: {
    '& label.Mui-focused': {
      color: '#5E81AC',
    },
    '& .MuiInput-underline:after': {
      borderBottomColor: '#5E81AC',
    },
    '& .MuiOutlinedInput-root': {
      '&:hover fieldset': {
        borderColor: '#5E81AC',
      },
      '&.Mui-focused fieldset': {
        borderColor: '#5E81AC',
      },
    },
  },
})(TextField);

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(22),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: "#86C0D2",
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
    backgroundColor: "#86C1D3",
    "&:hover": {
      backgroundColor: "#5E81AC",
      color: "#FFFFFF",
    },
  },
  text: {
    color: "#5E81AC",
    fontWeight: "bold"
  }
}));

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const dispatch = useDispatch();
  const ToggleSnackbar = useCallback(
    (vertical, horizontal, msg, color) =>
      dispatch(toggleSnackbar(vertical, horizontal, msg, color)),
    [dispatch]
  );
  const SetLoginStatus = useCallback(
    status => dispatch(setLoginStatus(status)),
    [dispatch]
  );

  let history = useHistory();
  const classes = useStyles();

  const login = e => {
    e.preventDefault();
    API.post("/user/login", {
      userName: email,
      password: password,
    })
      .then(response => {
        setLoading(false);
        // 本地保存用户登录状态和数据
        Auth.authenticate(response.data);
        // 全局状态
        SetLoginStatus(true);
        history.push("/home")
        ToggleSnackbar("top", "center", "登录成功", "success");
      })
      .catch(error => {
        setLoading(false);
        ToggleSnackbar("top", "center", error.message, "warning");
      })
  }

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5" style={{ color: "#4E5668", fontFamily: "Inter" }}>
          Login
        </Typography>
        <form className={classes.form} noValidate onSubmit={login}>
          <CssTextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="email"
            label="邮箱"
            name="email"
            autoComplete="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            autoFocus
          />
          <CssTextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            name="password"
            label="密码"
            type="password"
            id="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            autoComplete="current-password"
          />
          {/* <FormControlLabel
            control={<Checkbox value="remember" color="primary" />}
            label="记住我"
          /> */}
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
            disabled={loading}
            disableElevation
          >
            登录
          </Button>
          <Grid container >
            <Grid item xs>
              <Link href="#" variant="body2" className={classes.text}>
                忘记账号
              </Link>
            </Grid>
            <Grid item>
              <Link href="#" variant="body2" className={classes.text}>
                {"没有账号？来注册吧"}
              </Link>
            </Grid>
          </Grid>
        </form>
      </div>
      {/* <Box mt={8}>
        <Copyright />
      </Box> */}
    </Container>
  );
}