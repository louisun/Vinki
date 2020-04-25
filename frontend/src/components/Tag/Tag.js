import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import { faMarkdown } from '@fortawesome/free-brands-svg-icons';
import { faHashtag } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Box from '@material-ui/core/Box';
import Chip from '@material-ui/core/Chip';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import {
  withStyles,
  withTheme,
} from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';

import {
  setArticleList,
  setCurrentTag,
  setCurrentTopTag,
  setSecondTags,
  setSubTags,
  setTopTags,
} from '../../actions';
import API from '../../middleware/Api';
import { isEmptyObject } from '../../utils';

const mapStateToProps = (state) => {
    return {
        currentRepo: state.repo.currentRepo,
        currentTopTag: state.tag.currentTopTag,
        currentTag: state.tag.currentTag,
        topTags: state.tag.topTags,
        secondTags: state.tag.secondTags,
        subTags: state.tag.subTags,
        articleList: state.tag.articleList,
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
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
    };
};
const styles = (theme) => ({
    mid: {
        marginTop: "30px",
        flex: "1 0",
    },
    midWrapper: {
        margin: "auto",
        maxWidth: "1000px",
        minHeight: "70vh",
        padding: "20px 60px 20px 60px",
        backgroundColor: "#FFFFFF",
        color: "#4C566A",
        boxShadow:
            "0 4px 6px rgba(184,194,215,0.25), 0 5px 7px rgba(184,194,215,0.1)",
        borderRadius: "8px",
        marginBottom: "50px",
    },
    title: {
        marginTop: "30px",
        verticalAlign: "middle",
        textIndent: "0.5em",
    },
    view: {
        marginTop: "20px",
        marginBottom: "20px",
    },
    swicher: {
        fontSize: "15px",
        float: "right",
        marginRight: "30px",
    },
    chips: {
        display: "flex",
        justifyContent: "left",
        flexWrap: "wrap",
        padding: theme.spacing(0.5),
    },
    chip: {
        margin: theme.spacing(1),
        backgroundColor: "#E5E9F0",
        color: "#4C566A",
        fontSize: "15px",
        "&:hover": {
            backgroundColor: "#86C1D3",
            color: "#FFFFFF",
        },
        "&:focus": {
            backgroundColor: "#86C1D3",
            color: "#FFFFFF",
        },
        // fontWeight: "bold",
    },
    chipSelected: {
        margin: theme.spacing(1),
        backgroundColor: "#86C1D3",
        color: "#FFFFFF",
        textShadow: "0 0 .9px #FFF, 0 0 .9px #FFF",
        fontSize: "15px",
        "&:hover": {
            // backgroundColor: "#D8DEE9",
        },
        "&:focus": {
            backgroundColor: "#86C1D3",
            color: "#FFFFFF",
        },
    },
    topChip: {
        margin: theme.spacing(1),
        backgroundColor: "#E5E9F0",
        color: "#4C566A",
        fontSize: "15px",
        "&:hover": {
            backgroundColor: "#658BB9",
            color: "#FFFFFF",
        },
        "&:focus": {
            backgroundColor: "#658BB9",
            color: "#FFFFFF",
        },
        // fontWeight: "bold",
    },
    topChipSelected: {
        margin: theme.spacing(1),
        backgroundColor: "#5E81AC",
        color: "#FFFFFF",
        textShadow: "0 0 .9px #FFF, 0 0 .9px #FFF",
        fontSize: "15px",
        "&:hover": {
            // backgroundColor: "#D8DEE9",
        },
        "&:focus": {
            backgroundColor: "#658BB9",
            color: "#FFFFFF",
        },
    },
});

const TagUnfoldSwitch = withStyles({
    switchBase: {
        color: "#FFFFFF",
        "&$checked": {
            color: "#86C2D4",
        },
        "&$checked + $track": {
            backgroundColor: "#86C2D4",
        },
    },
    checked: {},
    thumb: {},
    track: {
        backgroundColor: "#4E5668",
    },
})(Switch);

class TagsComponent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            flat: false,
        };
    }

    loadTopTags() {
        API.get(`/repos/${this.props.currentRepo.ID}/tags`).then((response) => {
            if (response.data !== null) {
                this.props.setTopTags(response.data);
                this.props.setSecondTags([]);
                this.props.setSubTags([]);
                this.props.setCurrentTag({});
                this.props.setCurrentTopTag({});
                this.props.setArticleList([]);
            }
        });
    }
    componentDidMount() {
        if (
            !isEmptyObject(this.props.currentRepo) &&
            this.props.topTags.length === 0
        ) {
            this.loadTopTags();
        }
    }
    componentDidUpdate(prevProps, prevState) {
        // 1. 当一级标签为空时，刷新内容
        if (
            !isEmptyObject(this.props.currentRepo) &&
            this.props.topTags.length === 0
        ) {
            this.loadTopTags();
            return;
        }
        // 2. 当仓库变更时，刷新内容
        let refreshTopTags = false;
        if (!isEmptyObject(this.props.currentRepo)) {
            if (
                isEmptyObject(prevProps.currentRepo) ||
                this.props.currentRepo.ID !== prevProps.currentRepo.ID
            ) {
                refreshTopTags = true;
            }
        }
        if (refreshTopTags) {
            this.loadTopTags();
        }
    }

    handleSwitch = (event) => {
        // setState({...state, [event.target.name]: event.target.checked});
        this.setState({ flat: event.target.checked });
        if (event.target.checked === true) {
            API.get(`/tags/${this.props.currentTopTag.ID}/articles?flat=true`).then(
                (response) => {
                    this.props.setSecondTags(response.data.SubTags);
                    this.props.setSubTags([]);
                    this.props.setArticleList(response.data.ArticleInfos);
                }
            );
        } else {
            API.get(`/tags/${this.props.currentTopTag.ID}/articles`).then(
                (response) => {
                    this.props.setSecondTags(response.data.SubTags);
                    this.props.setArticleList(response.data.ArticleInfos);
                }
            );
        }
    };

    handleTagClick = (event, tag, type) => {
        this.setState((state) => ({ flat: false }));
        this.props.setCurrentTag(tag);
        API.get(`/tags/${tag.ID}/articles`).then((response) => {
            if (type === "top") {
                this.props.setCurrentTopTag(tag);
                this.props.setSecondTags(response.data.SubTags);
                this.props.setSubTags([]);
            } else {
                if (this.state.flat === false) {
                    this.props.setSubTags(response.data.SubTags);
                }
            }
            this.props.setArticleList(response.data.ArticleInfos);
        });
    };

    handleArticleClick = (event, id) => {
        // API.get(`/articles/${id}`).then(response => {
        //     this.props.setArticle(response.data)
        // })
        this.props.history.push(`/tag/${this.props.currentTag.ID}/article/${id}`);
    };

    render() {
        const { classes } = this.props;

        const generateTags = (tagList, type) => {
            let l = [];
            for (let i = 0; i < tagList.length; i++) {
                let className = "";

                if (tagList[i].ID === this.props.currentTopTag.ID) {
                    className = classes.topChipSelected;
                } else if (tagList[i].ID === this.props.currentTag.ID) {
                    className = classes.chipSelected;
                } else if (type === "top") {
                    className = classes.topChip;
                } else {
                    className = classes.chip;
                }
                l.push(
                    <Chip
                        key={i}
                        label={tagList[i].Name}
                        className={className}
                        clickable={true}
                        onClick={(event) => this.handleTagClick(event, tagList[i], type)}
                    />
                );
            }
            return l;
        };

        const generateArticles = (articleList) => {
            let l = [];
            for (let i = 0; i < articleList.length; i++) {
                l.push(
                    <ListItem
                        button={true}
                        onClick={(event) => {
                            this.handleArticleClick(event, articleList[i].ID);
                        }}
                    >
                        <ListItemIcon>
                            <FontAwesomeIcon
                                icon={faMarkdown}
                                style={{ color: "#4E5668", fontSize: "1.3em" }}
                            />
                        </ListItemIcon>
                        <ListItemText
                            primary={articleList[i].Title}
                            style={{
                                color: "#4C566A",
                                textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1",
                            }}
                            primaryTypographyProps={{ variant: "subtitle1" }}
                        />
                    </ListItem>
                );
            }
            return <List> {l} </List>;
        };
        return (
            <React.Fragment>
                <div className={classes.mid}>
                    <div className={classes.midWrapper}>
                        {this.props.topTags.length > 0 && (
                            <Box className={classes.view}>
                                <Typography
                                    variant="h5"
                                    gutterBottom
                                    className={classes.title}
                                >

                                    标签
                    </Typography>
                                <div className={classes.chips}>
                                    {generateTags(this.props.topTags, "top")}
                                </div>
                            </Box>
                        )}
                        {this.props.secondTags.length > 0 && (
                            <div>
                                <Divider variant="middle" />
                                <Box className={classes.view}>
                                    <Typography
                                        variant="h5"
                                        gutterBottom
                                        className={classes.title}
                                    >
                                        子标签
                        <span className={classes.swicher}>
                                            展开
                          <TagUnfoldSwitch
                                                checked={this.state.flat}
                                                onChange={this.handleSwitch}
                                                name="flat"
                                            />
                                        </span>
                                    </Typography>
                                    <div className={classes.chips}>
                                        {generateTags(this.props.secondTags, "second")}
                                    </div>
                                </Box>
                            </div>
                        )}
                        {this.props.subTags.length > 0 && (
                            <div>
                                <Divider variant="middle" />
                                <Box className={classes.view}>
                                    <div className={classes.chips}>
                                        {generateTags(this.props.subTags, "second")}
                                    </div>
                                </Box>
                            </div>
                        )}

                        {this.props.articleList.length > 0 && (
                            <div>
                                <Divider variant="middle" />
                                <Box className={classes.view}>
                                    <Typography
                                        variant="h5"
                                        gutterBottom
                                        className={classes.title}
                                    >
                                        文档列表
                      </Typography>
                                    <div>{generateArticles(this.props.articleList)}</div>
                                </Box>
                            </div>
                        )}
                    </div>
                </div>
            </React.Fragment>
        );
    }
}

TagsComponent.propTypes = {
    classes: PropTypes.object.isRequired,
    theme: PropTypes.object.isRequired,
};

const Tags = connect(
    mapStateToProps,
    mapDispatchToProps
)(withTheme(withStyles(styles)(withRouter(TagsComponent))));

export default Tags;
