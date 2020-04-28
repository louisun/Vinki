import '../../assets/css/markdown.css';
import '../../assets/css/nord.css';

import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import { faHashtag } from '@fortawesome/free-solid-svg-icons';
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

import {
  setArticleList,
  setCurrentTag,
} from '../../actions';
import cardImage from '../../assets/img/card.png';
import API from '../../middleware/Api';
import { isEmptyObject } from '../../utils';
import Hilight from './Highlight';

const mapStateToProps = state => {
    return {
        currentTag: state.tag.currentTag,
        articleList: state.tag.articleList,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        setCurrentTag: currentTag => {
            dispatch(setCurrentTag(currentTag))
        },
        setArticleList: articleList => {
            dispatch(setArticleList(articleList))
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
        overflowY: "scroll",
        '&::-webkit-scrollbar': {
            width: '0.4em'
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
    },
    toc: {
        marginTop: "30px",
        padding: "0 20px 0 20px",
        maxWidth: "250px",
        color: "#4C566A",
    },
    tocList: {
        maxHeight: "calc(100vh - 120px)",
        overflowY: "scroll",
        '&::-webkit-scrollbar': {
            width: '0.4em'
        },
        '&::-webkit-scrollbar-track': {
            boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
        },
        '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'rgba(0,0,0,.1)',
            outline: '1px solid slategrey',
        }
    },
    tocItem: {
        paddingTop: "4px",
        paddingBottom: "2px",
        borderLeft: "3px #F2F4F8 solid",
        "&:hover": {
            borderLeft: "3px #88C0D0 solid",
        },
    },
    tocH1Item: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#E9EAEE",
        borderLeft: "3px #88C0D0 solid",
        "&:hover": {
            borderLeft: "3px #88C0D0 solid",
        },
    },
    tocOnItemH2: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#dce4ee",
        borderLeft: "3px #5c7fa9 solid",
        "&:hover": {
            borderLeft: "3px #5c7fa9 solid",
        },
        fontWeight: "bold !important",
    },
    tocOnItemH3: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#edf2e8",
        borderLeft: "3px #96bb78 solid",
        "&:hover": {
            borderLeft: "3px #96bb78 solid",
        },
        fontWeight: "bold !important",
    },
    tocOnItemH4: {
        paddingTop: "4px",
        paddingBottom: "2px",
        backgroundColor: "#ECEFF4",
        borderLeft: "3px #d8bb7d solid",
        "&:hover": {
            borderLeft: "3px #d3a748 solid",
        },
        fontWeight: "bold !important",
    },
    tocH1: {
        color: "#2E3640 !important",
        fontSize: "1.1em",
    },
    tocH2: {
        // textIndent: "1em",
        color: "#2E3540 !important",
        fontWeight: "bold !important",
        fontSize: "1.1em",
    },
    tocH3: {
        // textIndent: "2em",
        color: "#4E5668 !important",
        fontSize: "1.1em",
    },
    tocH4: {
        textIndent: "1em",
        color: "#4E5668 !important",
        fontStyle: "italic",
    },
});


class ArticleComponent extends Component {
    constructor(props) {
        super(props)

        this.state = {
            article: {},
            headings: [], // 用来存储目录结构
            currentHeadingID: "",
        }
        this.handleScroll = this.handleScroll.bind(this)
    }

    // 第一次挂载组件
    componentDidMount() {
        API.get(`/articles/${this.props.match.params.articleID}`).then(response => {
            this.setState({ article: response.data })
            document.title = this.state.article.Title
        })
        if (isEmptyObject(this.props.currentTag)) {
            API.get(`/tags/${this.props.match.params.tagID}/articles`).then(response => {
                this.props.setCurrentTag({
                    ID: response.data.ID,
                    Name: response.data.Name,
                })
                this.props.setArticleList(response.data.ArticleInfos)
            })
        }
        setTimeout(() => {
            this.updateTocHeading(['H1', 'H2', 'H3', 'H4']);
        }, 300);
        window.addEventListener('scroll', this.handleScroll)
    }

    // 设置新文档后，组件更新
    componentDidUpdate(prevProps, prevState) {
        let updateArticle = false
        if (!isEmptyObject(this.state.article)) {
            if (isEmptyObject(prevState.article)) {
                updateArticle = true
            } else if (this.state.article.ID !== prevState.article.ID) {
                updateArticle = true
            }
            if (updateArticle) {
                document.title = this.state.article.Title
                this.updateTocHeading(['H1', 'H2', 'H3', 'H4']);
            }
        }
    }
    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScroll);
    }


    handleTagClick = (event) => {
        this.props.history.push("/tags")
    }

    handleArticleClick = (event, id) => {
        API.get(`/articles/${id}`).then(response => {
            this.setState({
                article: response.data,
                headings: [],
            })
        })
        window.scrollTo(0, 0)
    }

    updateTocHeading = (headings) => {
        if (document.querySelector(".markdown-content") === null) {
            return
        }
        // 找出 headings，给标签添加 id 属性，每个 TocItem.headingID 与其对应
        let headingList = [];
        let headingCountMap = {
            "H1": 0,
            "H2": 0,
            "H3": 0,
            "H4": 0,
        }
        let headingID = ""
        document.querySelector(".markdown-content").childNodes.forEach((item) => {
            if (headings.includes(item.nodeName)) {
                headingCountMap[item.nodeName]++
                headingID = `vinki-${item.nodeName}-${headingCountMap[item.nodeName]}`
                item.setAttribute("id", headingID)
                headingList.push({
                    type: item.nodeName,
                    headingID: headingID,
                    text: item.innerText,
                    offset: 0,
                });
            }
        });

        this.setState({
            headings: headingList,
            currentHeadingID: (headingList.length !== 0 ? headingList[0].headingID : ""),
        });
        // 2s 之后再设置偏移量
        setTimeout(() => {
            for (let item of headingList) {
                item.offset = document.getElementById(item.headingID).offsetTop
            }
            this.setState({
                headings: headingList,
            })
        }, 3000)
    }

    handleScroll = () => {
        let scrollTop = document.documentElement.scrollTop + 200 // 获取当前页面的滚动距离
        let currentHeadingID = this.state.headings.length !== 0 ? this.state.headings[0].headingID : ""

        for (let item of this.state.headings) {
            if (scrollTop >= item.offset) {
                currentHeadingID = item.headingID;
            } else {
                break;
            }
        }


        // 如果当前命中标题和前一个命中标题的文本不一样，说明当前页面处于其他标题下的内容，切换menuState
        if (currentHeadingID !== this.state.currentHeadingID) {
            this.setState({ currentHeadingID });
        }; // 获取当前页面的滚动距离
    }


    scrollToHeading = (item) => {
        let anchorEl = document.getElementById(item.headingID);
        if (anchorEl) {
            // bdoy 当前相对于「视口」的距离（向下滚动为负）
            const bodyRect = document.body.getBoundingClientRect().top;
            // 获取标签对于「视口」的距离
            const elementRect = anchorEl.getBoundingClientRect().top;
            // body 和 元素的实际距离，可以看做元素的 top
            const elPosition = elementRect - bodyRect;
            // 由于 header 占用了一些位置
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
                l.push(
                    <ListItem button={true} onClick={(event) => {
                        this.handleArticleClick(event, articleList[i].ID)
                    }}>
                        <ListItemText
                            primary={articleList[i].Title}
                            style={{ color: "#4C566A", textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1", }}
                            primaryTypographyProps={{ variant: "subtitle1" }}
                        />
                    </ListItem>,
                )
            }
            return <List className={classes.articleList}> {l} </List>
        }
        const generateTOC = (itemList) => {
            let l = []
            for (let i = 0; i < itemList.length; i++) {
                // 不同级别的 toc heading 样式：ListItem、Text
                let className1 = ""
                let className2 = ""
                if (itemList[i].type === "H1") {
                    className1 = classes.tocH1Item
                    className2 = classes.tocH1
                } else if (itemList[i].type === "H2") {
                    className1 = classes.tocItem
                    className2 = classes.tocH2
                    if (itemList[i].headingID === this.state.currentHeadingID) {
                        className1 = classes.tocOnItemH2
                    }
                } else if (itemList[i].type === "H3") {
                    className1 = classes.tocItem
                    className2 = classes.tocH3
                    if (itemList[i].headingID === this.state.currentHeadingID) {
                        className1 = classes.tocOnItemH3
                    }
                } else if (itemList[i].type === "H4") {
                    className1 = classes.tocItem
                    className2 = classes.tocH4
                    if (itemList[i].headingID === this.state.currentHeadingID) {
                        className1 = classes.tocOnItemH4
                    }
                }
                l.push(
                    <ListItem button={true} onClick={(event) => {
                        this.scrollToHeading(itemList[i])
                    }} className={className1}>
                        <span style={{ color: "#4C566A", textShadow: "0 0 .9px #E5E9F1, 0 0 .9px #E5E9F1", }} className={className2}>
                            {itemList[i].text}
                        </span>
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
                                            {this.props.currentTag.Name}
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
                        <div class="markdown-body" ref={(ref) => { this.scroll = ref }}>
                            <Hilight content={this.state.article.HTML}></Hilight>
                        </div>
                    </div>
                </div>

                <Hidden smDown>
                    <div className={classes.right}>
                        <div className={classes.rightWrapper}>
                            <div className={classes.toc}>
                                {generateTOC(this.state.headings)}
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
