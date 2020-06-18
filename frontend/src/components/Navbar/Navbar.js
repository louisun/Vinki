import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import {
  Link,
  withRouter,
} from 'react-router-dom';

import { IconButton } from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
// import IconButton from '@material-ui/core/IconButton';
// import InputBase from '@material-ui/core/InputBase';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import StyledMenu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import {
  fade,
  withStyles,
  withTheme,
} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import AccountCircle from '@material-ui/icons/AccountCircle';
import ClassRoundedIcon from '@material-ui/icons/ClassRounded';
import ExitToApp from '@material-ui/icons/ExitToApp';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
// import MenuIcon from '@material-ui/icons/Menu';
// import SearchIcon from '@material-ui/icons/Search';
import StarRoundedIcon from '@material-ui/icons/StarRounded';

import {
  setCurrentRepo,
  setLoginStatus,
  setRepos,
  toggleSnackbar,
} from '../../actions';
import API from '../../middleware/Api';
import Auth from '../../middleware/Auth';

const mapStateToProps = state => {
    return {
        currentRepo: state.repo.currentRepo,
        repos: state.repo.repos,
        isLogin: state.isLogin
    }
}

const mapDispatchToProps = dispatch => {
    return {
        setRepos: repos => {
            dispatch(setRepos(repos))
        },
        setCurrentRepo: currentRepo => {
            dispatch(setCurrentRepo(currentRepo))
        },
        toggleSnackbar: (vertical, horizontal, msg, color) => {
            dispatch(toggleSnackbar(vertical, horizontal, msg, color));
        },
        setLoginStatus: status => {
            dispatch(setLoginStatus(status));
        },
    }
}
const menuStyle = {
    paper: {
        border: '1px solid #d3d4d5',
    },
    list: {
        backgroundColor: '#F2F4F8',
        color: '#4C566A',
    }
}
const RepoMenu = withStyles(menuStyle)((props) => (
    <StyledMenu
        elevation={0}
        getContentAnchorEl={null}
        disableScrollLock={true}
        anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
        }}
        transformOrigin={{
            vertical: 'top',
            horizontal: 'center',
        }}
        {...props}
    />
));

const UserMenu = withStyles(menuStyle)((props) => (
    <StyledMenu
        elevation={0}
        getContentAnchorEl={null}
        disableScrollLock={true}
        anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
        }}
        transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
        }}
        {...props}
    />
));

const styles = (theme) => ({
    root: {
        display: "flex",
    },
    appbar: {
        backgroundColor: "#4E5668",
        flexGrow: 1,
        height: "50px",
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    brand: {
        color: "#86C0D2",
        fontWeight: "bold",
        textAlign: "center",
        marginLeft: "0.5em"
    },
    repo: {
        fontSize: "1.1rem",
        color: "#8FBCBB",
        fontWeight: "bold",
        background: '#4E5668',
        boxShadow: 'none',
        textTransform: 'none',
        '&:hover': {
            backgroundColor: '#81A1C1',
            color: '#E5E9F1',
            borderColor: '#4E5668',
            boxShadow: 'none',
        },
        '&:active': {
            color: '#E5E9F1',
            backgroundColor: '#81A1C1',
            borderColor: '#4E5668',
            boxShadow: 'none',
        },
    },
    grow: {
        flexGrow: 1,
    },
    search: {
        position: 'relative',
        borderRadius: theme.shape.borderRadius,
        backgroundColor: fade("#F2F4F8", 0.15),
        '&:hover': {
            backgroundColor: fade(theme.palette.common.white, 0.25),
        },
        marginLeft: 0,
        marginRight: "30px",
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            marginLeft: theme.spacing(1),
            width: 'auto',
        },
        [theme.breakpoints.down('sm')]: {
            display: 'none',
        },
    },
    searchIcon: {
        padding: theme.spacing(0, 2),
        height: '100%',
        position: 'absolute',
        pointerEvents: 'none',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    inputRoot: {
        color: 'inherit',
    },
    inputInput: {
        padding: theme.spacing(1, 1, 1, 0),
        // vertical padding + font size from searchIcon
        paddingLeft: `calc(1em + ${theme.spacing(4)}px)`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            width: '12ch',
            '&:focus': {
                width: '20ch',
            },
        },
    },

});

class NavbarComponent extends Component {
    constructor(props) {
        super(props)
        this.state = {
            repoAnchorEl: null,
            profileAnchorEl: null
        }
    }

    getRepos = () => {
        API.get("/repos").then(response => {
            if (response.data.length > 0) {
                this.props.setRepos(response.data)
                this.props.setCurrentRepo(response.data[0])
            }
        }).catch(() => { })
    }

    componentDidMount() {
        this.getRepos()
    }

    componentDidUpdate(prevProps, prevState) {
        if (!prevProps.isLogin && this.props.isLogin) {
            this.getRepos()
        }
    }

    handleRepoMenuOpen = (event) => {
        this.setState({ repoAnchorEl: event.currentTarget })
    };

    handleRepoMenuClose = () => {
        this.setState({ repoAnchorEl: null });
    };

    handleRepoMenuItemClick = (event) => {
        let repoName = event.currentTarget.attributes["repo_name"].nodeValue
        for (let i = 0; i < this.props.repos.length; i++) {
            if (this.props.repos[i] === repoName) {
                this.props.setCurrentRepo(repoName)
                break
            }
        }
        this.handleRepoMenuClose()
        this.props.history.push({
            pathname: "/tags",
        })
    }

    handleProfileMenuOpen = (event) => {
        this.setState({ profileAnchorEl: event.currentTarget })
    };


    handleProfileMenuClose = () => {
        this.setState({ profileAnchorEl: null });
    };


    logout = () => {
        API.post("/user/logout").then(response => {
            this.props.toggleSnackbar(
                "top",
                "center",
                "已退出登录",
                "success"
            );
            Auth.Signout()
            window.location.reload();
            this.props.setLoginStatus(false);
        }).catch(error => {
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "warning"
            );
        }).then(() => {
            this.handleProfileMenuClose()
        })
    }



    render() {
        const { classes } = this.props
        const profileMenuID = 'user-menu'

        const generateRepos = (repoList) => {
            let l = []
            for (let i = 0; i < repoList.length; i++) {
                let className = ""
                if (this.props.currentRepo === repoList[i]) {
                    className = classes.repoSelected
                } else {
                    className = classes.repoMenu
                }
                l.push(
                    <MenuItem className={className} onClick={this.handleRepoMenuItemClick} repo_name={repoList[i]}>
                        <ListItemIcon>
                            <ClassRoundedIcon fontSize="small" />
                        </ListItemIcon>
                        <ListItemText primary={repoList[i]} />
                    </MenuItem>
                )
            }
            return l
        }
        return (
            <div className={classes.root}>
                <AppBar position="fixed" className={classes.appbar}>
                    <Toolbar variant="dense">
                        {/* <IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu">
                            <MenuIcon />
                        </IconButton> */}
                        <Link to="/tags" style={{ textDecoration: "none", }}>
                            <Typography variant="h4" className={classes.brand}>
                                Vinki
                            </Typography>
                        </Link>
                        <div className={classes.grow} />
                        {/* <div className={classes.search}>
                            <div className={classes.searchIcon}>
                                <SearchIcon />
                            </div>
                            <InputBase
                                placeholder="搜索"
                                classes={{
                                    root: classes.inputRoot,
                                    input: classes.inputInput,
                                }}
                                inputProps={{ 'aria-label': 'search' }}
                            />
                        </div> */}
                        <div>
                            {this.props.isLogin === true ? (
                                <div>
                                    <Button
                                        aria-controls="customized-menu"
                                        aria-haspopup="true"
                                        variant="contained"
                                        size="small"
                                        onClick={this.handleRepoMenuOpen}
                                        startIcon={<StarRoundedIcon />}
                                        endIcon={<KeyboardArrowDownIcon />}
                                        className={classes.repo}
                                        disableRipple
                                    >
                                        {this.props.currentRepo !== "" ? this.props.currentRepo : '无仓库'}
                                    </Button>
                                    <RepoMenu
                                        id="repo-menu"
                                        anchorEl={this.state.repoAnchorEl}
                                        keepMounted
                                        open={Boolean(this.state.repoAnchorEl)}
                                        onClose={this.handleRepoMenuClose}
                                    >
                                        {generateRepos(this.props.repos)}
                                    </RepoMenu>
                                </div>
                            ) : (<span style={{ fontSize: "1.2em", color: "#E5E9F0", fontWeight: "bold" }}>未登录</span>)}
                        </div>
                        {this.props.isLogin ? (
                            <div>
                                <IconButton
                                    edge="end"
                                    aria-label="account of current user"
                                    aria-controls={profileMenuID}
                                    aria-haspopup="true"
                                    onClick={this.handleProfileMenuOpen}
                                    color="inherit"
                                >
                                    <AccountCircle style={{ fontSize: "1.2em", color: "#8AABCE" }} />
                                </IconButton>
                                <UserMenu
                                    id={"user-menu"}
                                    anchorEl={this.state.profileAnchorEl}
                                    keepMounted
                                    transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                                    open={Boolean(this.state.profileAnchorEl)}
                                    onClose={this.handleProfileMenuClose}
                                >
                                    <MenuItem onClick={this.logout}>
                                        <ListItemIcon>
                                            <ExitToApp fontSize="small" />
                                        </ListItemIcon>
                                        <ListItemText primary={"登出"} />
                                    </MenuItem>
                                </UserMenu>
                            </div>
                        ) : <div />}
                    </Toolbar>
                </AppBar>
            </div>
        )
    }
}

NavbarComponent.propTypes = {
    classes: PropTypes.object.isRequired,
    theme: PropTypes.object.isRequired
};

const Navbar = connect(
    mapStateToProps,
    mapDispatchToProps
)(withTheme(withStyles(styles)(withRouter(NavbarComponent))));

export default Navbar;
