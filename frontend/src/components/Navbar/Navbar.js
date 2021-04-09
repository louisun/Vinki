import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import {
  Link,
  withRouter,
} from 'react-router-dom';

import { faMarkdown } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  Backdrop,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Dialog,
  IconButton,
  List,
} from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import InputBase from '@material-ui/core/InputBase';
import ListItem from '@material-ui/core/ListItem';
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
import ImportExportIcon from '@material-ui/icons/ImportExport';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import SearchIcon from '@material-ui/icons/Search';
import StarRoundedIcon from '@material-ui/icons/StarRounded';
import SyncIcon from '@material-ui/icons/Sync';
import SyncAltOutlinedIcon from '@material-ui/icons/SyncAltOutlined';

import {
  setArticleList,
  setCurrentRepo,
  setCurrentTag,
  setCurrentTopTag,
  setLoginStatus,
  setRepos,
  setSecondTags,
  setSubTags,
  setTopTags,
  toggleSnackbar,
} from '../../actions';
import API from '../../middleware/Api';
import Auth from '../../middleware/Auth';

const mapStateToProps = state => {
    return {
        currentRepo: state.repo.currentRepo,
        currentTag: state.tag.currentTag,
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
        setCurrentTopTag: (currentTopTag) => {
            dispatch(setCurrentTopTag(currentTopTag));
        },
        setCurrentTag: (currentTag) => {
            dispatch(setCurrentTag(currentTag));
        },
        setTopTags: (topTags) => {
            dispatch(setTopTags(topTags));
        },
        setSecondTags: (secondTags) => {
            dispatch(setSecondTags(secondTags));
        },
        setSubTags: (subTags) => {
            dispatch(setSubTags(subTags));
        },
        setArticleList: (articleList) => {
            dispatch(setArticleList(articleList));
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
    backdrop: {
        zIndex: theme.zIndex.drawer + 1,
        color: '#fff',
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
        [theme.breakpoints.down('xs')]: {
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
            width: '20ch',
        },
    },
    searchCard: {
        minWidth: "500px",
        padding: "20px",
        overflow: "auto",
    },
    chips: {
        display: "flex",
        justifyContent: "left",
        flexWrap: "wrap",
        padding: theme.spacing(0.5),
    },
    tagChip: {
        margin: theme.spacing(1),
        backgroundColor: "#E5E9F0",
        color: "#4C566A",
        fontSize: "15px",
        fontFamily: "Inter",
        "&:hover": {
            backgroundColor: "#86C1D3",
            color: "#FFFFFF",
            textShadow: "0 0 .9px #FFF, 0 0 .9px #FFF",
        },
        "&:focus": {
            backgroundColor: "#86C1D3",
            color: "#FFFFFF",
            textShadow: "0 0 .9px #FFF, 0 0 .9px #FFF",
        },
        // fontWeight: "bold",
    },
    articleChip: {
        backgroundColor: "#86C1D3",
        color: "#FFFFFF",
        textShadow: "0 0 .9px #FFF, 0 0 .9px #FFF",
        fontFamily: "Inter",
    },

});

class NavbarComponent extends Component {
    constructor(props) {
        super(props)
        this.state = {
            user: null,
            repoAnchorEl: null,
            profileAnchorEl: null,
            loadingStatus: false,
            searchDialog: false,
            searchKeywords: "",
            searchTagsResult: null,
            searchArticlesResult: null,
        }
    }

    getRepos = () => {
        API.get("/repos").then(repoResp => {
            if (repoResp.data.length > 0) {
                this.props.setRepos(repoResp.data)
                API.get("/admin/config/repo").then(response => {
                    if (response.data == "") {
                        // 如果数据的当前仓库为空，默认为列表第一个
                        this.props.setCurrentRepo(repoResp.data[0])
                    } else {
                        // 数据库当前仓库非空，才设置该仓库
                        this.props.setCurrentRepo(response.data)
                    }
                }).catch(error => {
                    this.props.toggleSnackbar(
                        "top",
                        "center",
                        error.message,
                        "error"
                    );
                });
            }
        }).catch(() => { })
    }

    setRepo = () => {

    }

    componentDidMount() {
        this.setState({ user: Auth.GetUser() })
        this.getRepos()
    }

    componentWillReceiveProps = (nextProps) => {
        this.setState({ user: Auth.GetUser() })
        if (!this.props.isLogin && nextProps.isLogin) {
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
                API.post("/admin/config/repo", {
                    currentRepo: repoName
                }).catch(error => {
                    this.props.toggleSnackbar(
                        "top",
                        "center",
                        error.message,
                        "error"
                    );
                });
                break
            }
        }
        this.props.setTopTags([])
        this.props.setSecondTags([])
        this.props.setSubTags([])
        this.props.setArticleList([])
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

    handleLoadingOpen = () => {
        this.setState({ loadingStatus: true })
    }

    handleLoadingClose = () => {
        this.setState({ loadingStatus: false })
    }

    handleSearchDialogOpen = () => {
        this.setState({ searchDialog: true })
    }

    handleSearchDialogClose = () => {
        this.setState({ searchDialog: false })
    }

    onSearch = (e) => {
        if (e.key === 'Enter' && e.target.value !== "") {
            this.setState({ searchKeywords: e.target.value })
            API.get("/search", {
                params: {
                    type: "tag",
                    repo: this.props.currentRepo,
                    keyword: e.target.value
                }
            }).then(response => {
                this.setState({ searchTagsResult: response.data })
            }).catch(error => {
                this.props.toggleSnackbar(
                    "top",
                    "center",
                    error.message,
                    "error"
                );
            })
            API.get("/search", {
                params: {
                    type: "article",
                    repo: this.props.currentRepo,
                    keyword: e.target.value
                }
            }).then(response => {
                this.setState({ searchArticlesResult: response.data })
            }).catch(error => {
                this.props.toggleSnackbar(
                    "top",
                    "center",
                    error.message,
                    "error"
                );
            })
            this.handleSearchDialogOpen()
            e.target.value = ""
        }
    }

    logout = () => {
        API.post("/user/logout").then(() => {
            this.props.toggleSnackbar(
                "top",
                "center",
                "已退出登录",
                "success"
            );
            Auth.Signout()
            this.props.setLoginStatus(false);
            window.location.reload();
        }).catch(error => {
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        }).then(() => {
            this.handleProfileMenuClose()
        })
    }

    refreshAll = () => {
        this.handleLoadingOpen();
        API.post("/admin/refresh/all").then(response => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                response.rawData.msg,
                "success"
            );
            window.location.reload()
        }).catch(error => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        })
    }

    refreshRepo = () => {
        this.handleLoadingOpen();
        API.post("/admin/refresh/repo", {
            repoName: this.props.currentRepo
        }).then(response => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                response.rawData.msg,
                "success"
            );
            window.location.reload()
        }).catch(error => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        })
    }

    refreshTag = () => {
        this.handleLoadingOpen();
        API.post("/admin/refresh/tag", {
            repoName: this.props.currentRepo,
            tagName: this.props.currentTag
        }).then(response => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                response.rawData.msg,
                "success"
            );
            window.location.reload()
        }).catch(error => {
            this.handleProfileMenuClose()
            this.handleLoadingClose()
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        })
    }

    handleTagClick = (event, tag) => {
        this.props.setCurrentTag(tag)
        API.get("/tag", {
            params: {
                flat: false,
                repoName: this.props.currentRepo,
                tagName: tag
            }
        }).then(response => {
            if (response.data == null) {
                return;
            }
            if (tag.split("|").length === 1) {
                this.props.setCurrentTopTag(tag);
                if (response.data.SubTags) {
                    this.props.setSecondTags(response.data.SubTags);
                } else {
                    this.props.setSecondTags([]);
                }
                this.props.setSubTags([])
            } else if (tag.split("|").length > 1) {
                this.props.setCurrentTopTag(tag.split("|")[0]);
                this.props.setSecondTags([tag]);
                this.props.setSubTags(response.data.SubTags);
            }
            this.props.setArticleList(response.data.ArticleInfos)
            this.handleSearchDialogClose()
            this.props.history.push({
                pathname: "/tags",
            })
        }).catch(error => {
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        });
    }

    handleArticleClick = (event, tag, articleName) => {
        this.props.setCurrentTag(tag);
        API.get("/tag", {
            params: {
                flat: false,
                repoName: this.props.currentRepo,
                tagName: tag
            }
        }).then(response => {
            if (response.data == null) {
                return;
            }
            if (tag.split("|").length === 1) {
                this.props.setCurrentTopTag(tag);
                if (response.data.SubTags) {
                    this.props.setSecondTags(response.data.SubTags);
                } else {
                    this.props.setSecondTags([]);
                }
                this.props.setSubTags([])
            } else if (tag.split("|").length > 1) {
                this.props.setSubTags(response.data.SubTags);
            }
            this.props.setArticleList(response.data.ArticleInfos)
            this.handleSearchDialogClose()
            this.props.history.push(`/article/${this.props.currentRepo}/${tag}/${articleName}`);
        }).catch(error => {
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        });
    };

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
                    <MenuItem className={className} onClick={this.handleRepoMenuItemClick} repo_name={repoList[i]} key={"nav_repo" + i}>
                        <ListItemIcon>
                            <ClassRoundedIcon fontSize="small" />
                        </ListItemIcon>
                        <ListItemText primary={repoList[i]} />
                    </MenuItem>
                )
            }
            return l
        }

        const generateSearchTags = (tagList) => {
            let l = [];
            for (let i = 0; i < tagList.length; i++) {
                l.push(
                    <Chip
                        key={"nav_search_tag" + i}
                        label={tagList[i].split("|").join(" | ")}
                        className={classes.tagChip}
                        clickable={true}
                        onClick={(event) => this.handleTagClick(event, tagList[i])}
                    />
                );
            }
            return l;
        }


        const generateSearchArticles = (articleList) => {
            let l = [];
            for (let i = 0; i < articleList.length; i++) {
                l.push(
                    <ListItem
                        button={true}
                        onClick={(event) => {
                            this.handleArticleClick(event, articleList[i].tag, articleList[i].article);
                        }}
                    >
                        <ListItemIcon>
                            <FontAwesomeIcon
                                icon={faMarkdown}
                                style={{ color: "#4E5668", fontSize: "1.3em" }}
                            />
                        </ListItemIcon>
                        <ListItemText
                            primary={articleList[i].article}
                            style={{
                                color: "#4C566A",
                                textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1",
                            }}
                            primaryTypographyProps={{ variant: "subtitle1" }}
                        />
                        <Chip
                            key={"nav_search_article_tag" + i}
                            label={articleList[i].tag.split("|").join(" | ")}
                            className={classes.articleChip}
                        />
                    </ListItem>
                );
            }
            return <List> {l} </List>;
        };

        return (
            <div>
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
                            {(this.state.user != null && this.state.user.is_admin && this.props.currentRepo !== "") ? (
                                <div className={classes.search}>
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
                                        onKeyPress={this.onSearch}
                                    />
                                </div>

                            ) : (<div />)}
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
                                        {(this.state.user != null && this.state.user.is_admin) ? (
                                            <div>
                                                <MenuItem onClick={this.refreshTag}>
                                                    <ListItemIcon>
                                                        <SyncAltOutlinedIcon fontSize="small" />
                                                    </ListItemIcon>
                                                    <ListItemText primary={"标签同步"} />
                                                </MenuItem>
                                                <MenuItem onClick={this.refreshRepo}>
                                                    <ListItemIcon>
                                                        <ImportExportIcon fontSize="small" />
                                                    </ListItemIcon>
                                                    <ListItemText primary={"仓库同步"} />
                                                </MenuItem>
                                                <MenuItem onClick={this.refreshAll}>
                                                    <ListItemIcon>
                                                        <SyncIcon fontSize="small" />
                                                    </ListItemIcon>
                                                    <ListItemText primary={"全局同步"} />
                                                </MenuItem>
                                            </div>
                                        ) : (<div />)}
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
                <Backdrop className={classes.backdrop} open={this.state.loadingStatus}>
                    <CircularProgress color="inherit" />
                </Backdrop>
                <Dialog open={this.state.searchDialog} onClose={this.handleSearchDialogClose} aria-labelledby="form-dialog-title">
                    <Card className={classes.searchCard}>
                        <CardContent>
                            <Typography variant="body2" color="textSecondary" component="p">搜索 {this.state.searchKeywords} 结果</Typography>
                        </CardContent>
                        {this.state.searchTagsResult !== null ? (
                            <div className={classes.chips}>
                                {generateSearchTags(this.state.searchTagsResult)}
                            </div>
                        ) : (<div />)}
                        {this.state.searchArticlesResult !== null ? (
                            <div>
                                {generateSearchArticles(this.state.searchArticlesResult)}
                            </div>
                        ) : (<div />)}
                    </Card>
                </Dialog>
            </div >
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
