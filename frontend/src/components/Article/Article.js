import '../../assets/css/markdown.css';
import '../../assets/css/nord.css';

import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import {
  faHashtag,
  faListUl,
} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Hidden from '@material-ui/core/Hidden';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import {
  withStyles,
  withTheme,
} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Skeleton from '@material-ui/lab/Skeleton';

import {
  setArticleList,
  setCurrentRepo,
  setCurrentTag,
  toggleSnackbar,
} from '../../actions';
import cardImage from '../../assets/img/card.png';
import API from '../../middleware/Api';
import {
  isEmptyObject,
  lastOfArray,
} from '../../utils';
import Hilight from './Highlight';
import Latex from './Latex';

const mapStateToProps = state => {
    return {
        currentRepo: state.repo.currentRepo,
        currentTag: state.tag.currentTag,
        articleList: state.tag.articleList,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        setCurrentRepo: currentRepo => {
            dispatch(setCurrentRepo(currentRepo))
        },
        setCurrentTag: currentTag => {
            dispatch(setCurrentTag(currentTag))
        },
        setArticleList: articleList => {
            dispatch(setArticleList(articleList))
        },
        toggleSnackbar: (vertical, horizontal, msg, color) => {
            dispatch(toggleSnackbar(vertical, horizontal, msg, color));
        },
    }
}


const styles = (theme) => ({
    left: {
        flex: "0 0 250px",
    },
    leftWrapper: {
        position: "fixed",
        top: "50px",
        left: "0",
    },
    mid: {
        marginTop: "30px",
        minWidth: "300px",
        flex: "1 1",
    },
    midWrapper: {
        margin: "auto",
        maxWidth: "1000px",
        padding: "20px 60px 20px 60px",
        backgroundColor: "#FFFFFF",
        boxShadow: "0 4px 6px rgba(184,194,215,0.25), 0 5px 7px rgba(184,194,215,0.1)",
        // boxShadow: "0 12px 15px 0 rgba(0,0,0,0.24), 0 17px 50px 0 rgba(0,0,0,0.19)",
        borderRadius: "8px",
        marginBottom: "50px",
    },
    card: {
        marginTop: "30px",
        maxWidth: "220px",
        marginLeft: "20px",
    },
    articleList: {
        maxHeight: "50vh",
        overflow: "scroll",
        wordBreak: "break-all",
        '&::-webkit-scrollbar': {
            width: '0.4em',
            height: '0.4em'
        },
        '&::-webkit-scrollbar-track': {
            boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
        },
        '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'rgba(0,0,0,.1)',
            outline: '1px solid slategrey',
        }
    },
    right: {
        flex: "0 0 250px",
    },
    rightWrapper: {
        position: "fixed",
        top: "50px",
        width: "100%"
    },
    toc: {
        marginTop: "30px",
        padding: "0 20px 0 20px",
        maxWidth: "250px",
        color: "#4C566A",
    },

    tocList: {
        maxHeight: "calc(100vh - 120px)",
        wordBreak: "break-all",
        overflowY: "scroll",
        '&::-webkit-scrollbar': {
            width: '0.4em',
            height: '0.4em'
        },
        '&::-webkit-scrollbar-track': {
            boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
        },
        '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'rgba(0,0,0,.1)',
            outline: '1px solid slategrey',
        }
    },
    tocSelected: {
        backgroundColor: "#E9EAEE !important",
        // borderLeft: "4px #80B2C3 solid !important",
        // "&:hover": {
        //     borderLeft: "4px #80B2C3 solid !important",
        // },
    },
    tocTextSelected: {
        fontWeight: "bold !important",
    },
    tocH1Item: {
        paddingTop: "4px",
        paddingBottom: "2px",
        borderLeft: "4px #88C0D0 solid",
        "&:hover": {
            borderLeft: "4px #88C0D0 solid",
        },
    },
    tocH2Item: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#dce4ee",
        borderLeft: "4px #5c7fa9 solid",
        "&:hover": {
            borderLeft: "4px #4C698C solid",
        },
    },
    tocH3Item: {
        paddingTop: "4px",
        paddingBottom: "2px",
        borderLeft: "4px #96bb78 solid",
        // borderBottom: "1px black dashed",
        "&:hover": {
            borderLeft: "4px #7AA05B solid",
        },
    },
    tocH4Item: {
        paddingTop: "4px",
        paddingBottom: "2px",
        paddingLeft: "30px",
        borderLeft: "4px #d8bb7d solid",
        "&:hover": {
            borderLeft: "4px #d3a748 solid",
        },
    },
    tocTextCommon: {
        color: "#4C566A", textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1",
    },
    tocH1Text: {
        fontFamily: "Inter",
        color: "#2E3640",
        fontSize: "1.2em",
    },
    tocH2Text: {
        fontFamily: "Inter",
        color: "#364C65",
        fontSize: "1.1em",
        fontWeight: "500",
    },
    tocH3Text: {
        fontFamily: "Inter",
        color: "#4E5668",
        fontSize: "1.0em",
        fontWeight: "400",
    },
    tocH4Text: {
        width: "200px",
        // textIndent: "1em",
        color: "#4E5668",
        fontFamily: "Inter",
        fontSize: "0.9em",
        fontWeight: "350",
    },
    bottomBar: {
        position: "fixed",
        width: "100%",
        height: "40px",
        backgroundColor: "#383D4A",
        bottom: 0,
        boxShadow: "0px -2px 4px -1px rgba(0,0,0,0.2), 0px -4px 5px 0px rgba(0,0,0,0.14), 0px -1px 10px 0px rgba(0,0,0,0.12)",
    },
    // 对应 toc
    tocBottom: {
        position: "fixed",
        width: "100%",
        overflowX: "hidden",
        backgroundColor: "#4E5668",
        color: "#FFFFFF",
        bottom: "40px",
        height: "50vh",
        paddingTop: "10px",
        paddingLeft: "10px",
    },
    // 对应 toc list
    tocListBottom: {
        wordBreak: "break-all",
        overflowY: "scroll",
        '&::-webkit-scrollbar': {
            width: '0.4em',
            height: '0.4em'
        },
        '&::-webkit-scrollbar-track': {
            boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
        },
        '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'rgba(0,0,0,.1)',
            outline: '1px solid slategrey',
        }
    },
    tocH1ItemBottom: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#4E5668",
        borderLeft: "3px #88C0D0 solid",
        "&:hover": {
            borderLeft: "3px #88C0D0 solid",
        },
    },
    tocH2ItemBottom: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#dce4ee",
        borderLeft: "3px #5c7fa9 solid",
        "&:hover": {
            borderLeft: "3px #5c7fa9 solid",
        },
        fontWeight: "bold !important",
    },
    tocH3ItemBottom: {
        paddingTop: "4px",
        paddingBottom: "2px",
        // backgroundColor: "#012E46",
        borderLeft: "3px #96bb78 solid",
        "&:hover": {
            borderLeft: "3px #96bb78 solid",
        },
        fontWeight: "bold !important",
    },
    tocH4ItemBottom: {
        paddingTop: "4px",
        paddingBottom: "2px",
        // backgroundColor: "#012E46",
        borderLeft: "3px #d8bb7d solid",
        "&:hover": {
            borderLeft: "3px #d3a748 solid",
        },
        fontWeight: "bold !important",
    },
    tocH1TextBottom: {
        color: "#FFFFFF !important",
        fontWeight: "bold !important",
        fontSize: "1.1em",
        fontFamily: "Inter"
    },
    tocH2TextBottom: {
        // textIndent: "1em",
        color: "#364C65 !important",
        fontWeight: "bold !important",
        fontSize: "1.1em",
        fontFamily: "Inter"
    },
    tocH3TextBottom: {
        // textIndent: "2em",
        color: "#FFFFFF !important",
        fontSize: "1.1em",
        fontFamily: "Inter"
    },
    tocH4TextBottom: {
        textIndent: "1em",
        color: "#FFFFFF !important",
        fontStyle: "italic",
        fontFamily: "Inter"
    },
});


class ArticleComponent extends Component {
    constructor(props) {
        super(props)
        this.state = {
            article: {},
            headings: [], // 用来存储目录结构
            tocMenu: false,
            tocIndex: 0,
        }
        this.tocTimer = null
    }

    // 第一次挂载组件
    componentDidMount() {
        API.get("/article", {
            params: {
                repoName: this.props.match.params.repoName,
                tagName: this.props.match.params.tagName,
                articleName: this.props.match.params.articleName,
            }
        }).then(response => {
            if (response.data) {
                this.setState({ article: response.data })
                document.title = this.state.article.Title
            }
        }).catch(error => {
            this.props.toggleSnackbar(
                "top",
                "center",
                error.message,
                "error"
            );
        });
        if (this.props.currentRepo === "" || this.props.currentRepo !== this.props.match.params.repoName) {
            this.props.setCurrentRepo(this.props.match.params.repoName)
            API.post("/admin/config/repo", {
                currentRepo: this.props.match.params.repoName
            }).catch(error => {
                this.props.toggleSnackbar(
                    "top",
                    "center",
                    error.message,
                    "error"
                );
            });
        }
        if (this.props.currentTag === "") {
            API.get("/tag", {
                params: {
                    repoName: this.props.match.params.repoName,
                    tagName: this.props.match.params.tagName,
                }
            }).then(response => {
                this.props.setCurrentTag(response.data.Name)
                this.props.setArticleList(response.data.ArticleInfos)
            }).catch(error => {
                this.props.toggleSnackbar(
                    "top",
                    "center",
                    error.message,
                    "error"
                );
            });
        }

        this.tocTimer = setTimeout(this.updateTocHeading.bind(this), 300)
    }

    // 设置新文档后，组件更新
    // componentDidUpdate(prevProps, prevState) {
    //     if (!isEmptyObject(this.state.article)) {
    //         let updateArticle = false
    //         if (isEmptyObject(prevState.article)) {
    //             updateArticle = true
    //         } else if (this.state.article.ID !== prevState.article.ID) {
    //             updateArticle = true
    //         }
    //         if (updateArticle) {
    //             document.title = this.state.article.Title
    //             this.updateTocHeading(['H1', 'H2', 'H3', 'H4']);
    //         }
    //     }
    // }

    componentWillUnmount() {
        clearTimeout(this.tocTimer)
    }


    handleTagClick = (event) => {
        this.props.history.push("/tags")
    }

    handleArticleClick = (event, articleName) => {
        this.props.history.push(`/article/${this.props.match.params.repoName}/${this.props.match.params.tagName}/${articleName}`)
    }

    handleTocMenuClick = (event) => {
        this.setState({ tocMenu: !this.state.tocMenu })
    }

    updateTocHeading = () => {
        // 找出 headings，给标签添加 id 属性，每个 TocItem.headingID 与其对应
        let headings = ['H1', 'H2', 'H3', 'H4']
        let headingList = [];
        let headingCountMap = {
            "H1": 0,
            "H2": 0,
            "H3": 0,
            "H4": 0,
        }
        let headingID = ""
        setTimeout(() => {
            if (document.querySelector(".markdown-content") === null) {
                // TDOO 重试
                return
            }
            document.querySelector(".markdown-content").childNodes.forEach((item) => {
                if (headings.includes(item.nodeName)) {
                    headingCountMap[item.nodeName]++
                    headingID = `vinki-${item.nodeName}-${headingCountMap[item.nodeName]}`
                    item.setAttribute("id", headingID)
                    headingList.push({
                        type: item.nodeName,
                        headingID: headingID,
                        text: item.innerText,
                    });
                }
            });
            this.setState({
                headings: headingList,
            });
        }, 200)
    }

    // updateTocOffset = () => {
    //     // 过一会再设置偏移量
    //     setTimeout(() => {
    //         let headingList = this.state.headings
    //         for (let item of headingList) {
    //             let header = document.getElementById(item.headingID)
    //             if (header == null) {
    //                 // heading 还没渲染
    //                 return
    //             } else {
    //                 item.offset = header.offsetTop
    //             }
    //         }
    //         this.setState({
    //             headings: headingList,
    //         })
    //     }, 1000)
    // }

    // handleScroll = () => {
    //     let scrollTop = document.documentElement.scrollTop + 300 // 获取当前页面的滚动距离
    //     let currentHeadingID = this.state.headings.length !== 0 ? this.state.headings[0].headingID : ""

    //     for (let item of this.state.headings) {
    //         if (scrollTop >= item.offset) {
    //             currentHeadingID = item.headingID;
    //         } else {
    //             break;
    //         }
    //     }
    //     if (currentHeadingID !== this.state.currentHeadingID) {
    //         this.setState({ currentHeadingID });
    //     };

    // }


    scrollToHeading = (item, hideToc) => {
        if (hideToc) {
            this.setState({ tocMenu: false })
        }
        let anchorEl = document.getElementById(item.headingID);
        if (anchorEl) {
            // bdoy 当前相对于「视口」的距离（向下滚动为负）
            const bodyRect = document.body.getBoundingClientRect().top;
            // 获取标签对于「视口」的距离
            const elementRect = anchorEl.getBoundingClientRect().top;
            // body 和 元素的实际距离，可以看做元素的 top
            const elPosition = elementRect - bodyRect;
            // 由于 header 占用了一些位置，-60 为最上方
            const offPostion = elPosition - 60;
            window.scrollTo({
                top: offPostion,
                behavior: "instant"
            });
        }

    }

    render() {
        const { classes } = this.props
        const generateArticles = (articleList) => {
            let l = []
            for (let i = 0; i < articleList.length; i++) {
                if (articleList[i] === this.props.match.params.articleName) {
                    l.push(
                        <ListItem button={true} style={{ backgroundColor: "#F5F5F5" }} key={"a_al" + i}
                            onClick={(event) => {
                                this.handleArticleClick(event, articleList[i])
                            }}>
                            <ListItemText
                                primary={articleList[i]}
                                style={{ color: "#4C566A", textShadow: "0 0 .6px #2E3440, 0 0 .6px #2E3440", }}
                                primaryTypographyProps={{ variant: "subtitle1" }}
                            />
                        </ListItem>,
                    )

                } else {
                    l.push(
                        <ListItem button={true} key={"a_article_list" + i} onClick={(event) => {
                            this.handleArticleClick(event, articleList[i])
                        }}>
                            <ListItemText
                                primary={articleList[i]}
                                style={{ color: "#4C566A", textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1", }}
                                primaryTypographyProps={{ variant: "subtitle1" }}
                            />
                        </ListItem>,
                    )
                }
            }
            return <List className={classes.articleList}> {l} </List>
        }
        const generateTOC = (itemList) => {
            let l = []
            for (let i = 0; i < itemList.length; i++) {
                // 不同级别的 toc heading 样式：ListItem、Text
                let itemClass = ""
                let textClass = ""
                if (itemList[i].type === "H1") {
                    itemClass = classes.tocH1Item
                    textClass = classes.tocH1Text
                } else if (itemList[i].type === "H2") {
                    itemClass = classes.tocH2Item
                    textClass = classes.tocH2Text
                } else if (itemList[i].type === "H3") {
                    itemClass = classes.tocH3Item
                    textClass = classes.tocH3Text
                } else if (itemList[i].type === "H4") {
                    itemClass = classes.tocH4Item
                    textClass = classes.tocH4Text
                }
                textClass = textClass + ' ' + classes.tocTextCommon
                if (this.state.tocIndex === i) {
                    itemClass = itemClass + ' ' + classes.tocSelected
                    textClass = textClass + ' ' + classes.tocTextSelected
                }
                l.push(
                    <ListItem button={true} disableRipple key={"a_toc" + i} onClick={(event) => {
                        this.setState({ tocIndex: i })
                        this.scrollToHeading(itemList[i], false)
                    }} className={itemClass}>
                        <span className={textClass}> {itemList[i].text} </span>
                    </ListItem>,
                )
            }
            return <List className={classes.tocList}> {l} </List>
        }
        const generateTOCBottom = (itemList) => {
            let l = []
            for (let i = 0; i < itemList.length; i++) {
                // 不同级别的 toc heading 样式：ListItem、Text
                let itemClass = ""
                let textClass = ""
                if (itemList[i].type === "H1") {
                    itemClass = classes.tocH1ItemBottom
                    textClass = classes.tocH1TextBottom
                } else if (itemList[i].type === "H2") {
                    itemClass = classes.tocH2ItemBottom
                    textClass = classes.tocH2TextBottom
                } else if (itemList[i].type === "H3") {
                    itemClass = classes.tocH3ItemBottom
                    textClass = classes.tocH3TextBottom
                } else if (itemList[i].type === "H4") {
                    itemClass = classes.tocH4ItemBottom
                    textClass = classes.tocH4TextBottom
                }
                l.push(
                    <ListItem button={true} key={"a_toc_bottom" + i} onClick={(event) => {
                        this.scrollToHeading(itemList[i], true)
                    }} className={itemClass}>
                        <span className={textClass}> {itemList[i].text} </span>
                    </ListItem>,
                )
            }
            return <List className={classes.tocList}> {l} </List>
        }
        return (
            <React.Fragment >
                <Hidden mdDown>
                    <div className={classes.left}>
                        <div className={classes.leftWrapper}>
                            <Card className={classes.card}>
                                <CardActionArea onClick={this.handleTagClick}>
                                    <CardMedia
                                        component="img"
                                        alt="Contemplative Reptile"
                                        height="140"
                                        image={cardImage}
                                        title="Contemplative Reptile"
                                    />
                                    <CardContent>
                                        <Typography gutterBottom variant="h5" component="h2">
                                            {lastOfArray(this.props.currentTag.split("|"))}
                                        </Typography>
                                        <Typography variant="body2" color="textSecondary" component="p">
                                            <FontAwesomeIcon icon={faHashtag} style={{ fontSize: "0.9rem", marginRight: "0.4em" }} />
                                            标签
                                        </Typography>
                                    </CardContent>
                                </CardActionArea>
                                {generateArticles(this.props.articleList)}
                            </Card>
                        </div>
                    </div>
                </Hidden>

                <div className={classes.mid}  >
                    <div className={classes.midWrapper} >
                        <div className="markdown-body" ref={(ref) => { this.scroll = ref }}>
                            {isEmptyObject(this.state.article) ? (
                                <div>
                                    <Skeleton animation="wave" />
                                    <Skeleton animation="wave" />
                                    <Skeleton animation="wave" />
                                    <Skeleton animation="wave" />
                                </div>
                            ) : (
                                <Latex>
                                    <Hilight content={this.state.article.HTML}></Hilight>
                                </Latex>
                            )}
                        </div>
                    </div>
                    <Hidden smUp>
                        <div>
                            <div className={classes.tocBottom} style={{ display: this.state.tocMenu ? "block" : "none" }}>
                                {isEmptyObject(this.state.headings) ? (
                                    <div style={{ width: "200px" }}>
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                    </div>
                                ) : generateTOCBottom(this.state.headings)
                                }
                            </div>
                            <div className={classes.bottomBar}>
                                <Typography variant="h6" style={{ color: "#81A1C1", fontWeight: "bold", textAlign: "center", paddingTop: "5px" }} onClick={this.handleTocMenuClick}  >
                                    <FontAwesomeIcon icon={faListUl} style={{ fontSize: "1.1em", marginRight: "0.4em" }} />
                                    TOC
                                </Typography>
                            </div>
                        </div>
                    </Hidden>
                </div >

                <Hidden smDown>
                    <div className={classes.right}>
                        <div className={classes.rightWrapper}>
                            <div className={classes.toc}>
                                {isEmptyObject(this.state.headings) ? (
                                    <div style={{ width: "200px" }}>
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                        <Skeleton animation="wave" />
                                    </div>
                                ) : generateTOC(this.state.headings)
                                }
                            </div>
                        </div>
                    </div>
                </Hidden>

            </React.Fragment >
        )
    }
}

ArticleComponent.propTypes = {
    classes: PropTypes.object.isRequired,
    theme: PropTypes.object.isRequired
};

const Article = connect(
    mapStateToProps,
    mapDispatchToProps
)(withTheme(withStyles(styles)(withRouter(ArticleComponent))));

export default Article;
