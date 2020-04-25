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
        minWidth: "500px",
        flex: "1 1",
    },
    midWrapper: {
        margin: "auto",
        maxWidth: "1000px",
        padding: "20px 60px 20px 60px",
        backgroundColor: "#FFFFFF",
        boxShadow: "0 4px 6px rgba(184,194,215,0.25), 0 5px 7px rgba(184,194,215,0.1)",
        borderRadius: "8px",
        marginBottom: "50px"
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
        maxHeight: "70vh",
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
            borderLeft: "3px #4E5668 solid",
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
            refresh: false,
            headings: [], // 用来存储目录结构
        }
    }

    componentDidMount() {
        API.get(`/articles/${this.props.match.params.articleID}`).then(response => {
            this.setState({ article: response.data })
        })
        this.setState({ refresh: true })
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
            this.getHeadings(['H1', 'H2', 'H3', 'H4']);
        }, 300);
    }

    componentDidUpdate(prevProps, prevState) {
        let mountToc = false
        if (!isEmptyObject(this.state.article)) {
            if (isEmptyObject(prevState.article)) {
                mountToc = true
            } else if (this.state.article.ID != prevState.article.ID) {
                mountToc = true
            }
            if (mountToc) {
                setTimeout(() => {
                    this.getHeadings(['H1', 'H2', 'H3', 'H4']);
                }, 300);
            }
        }
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

    getHeadings = (headings) => {
        if (document.querySelector(".markdown-content") === null) {
            return
        }
        let headingList = [];
        document.querySelector(".markdown-content").childNodes.forEach((item) => {
            if (headings.includes(item.nodeName)) {
                headingList.push({
                    type: item.nodeName,
                    text: item.getAttribute('id'),
                    offsetTop: item.offsetTop - 60
                });
            }
        });

        this.setState({
            headings: headingList,
        });
    }


    scrollPage = (item) => {
        let anchorEl = document.getElementById(item.text);
        if (anchorEl) {
            const bodyRect = document.body.getBoundingClientRect().top;
            const elementRect = anchorEl.getBoundingClientRect().top;
            const elPosition = elementRect - bodyRect;
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
                let className1 = ""
                let className2 = ""
                if (itemList[i].type === "H1") {
                    className1 = classes.tocH1Item
                    className2 = classes.tocH1
                } else if (itemList[i].type === "H2") {
                    className1 = classes.tocItem
                    className2 = classes.tocH2
                } else if (itemList[i].type === "H3") {
                    className1 = classes.tocItem
                    className2 = classes.tocH3
                } else if (itemList[i].type === "H4") {
                    className1 = classes.tocItem
                    className2 = classes.tocH4
                }
                l.push(
                    <ListItem button={true} onClick={(event) => {
                        this.scrollPage(itemList[i])
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
            <React.Fragment>
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

                <div className={classes.mid}>
                    <div className={classes.midWrapper}>
                        <div class="markdown-body" ref={(ref) => { this.scroll = ref }} >
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
