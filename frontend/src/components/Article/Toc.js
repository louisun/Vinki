import React, { Component } from 'react';

import * as tocbot from 'tocbot';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import { isEmptyObject } from '../../utils';

const styles = (theme) => ({
  toc: {
    marginTop: "30px",
    padding: "0 20px 0 20px",
    maxWidth: "250px",
    color: "#4C566A",
  },
});


class TocComponent extends Component {
  constructor(props) {
    super(props)
    this.state = {
      menuData: [], // 用来存储目录结构
      menuState: '' // 用来存储当前命中的标题
    }
  }

  // 因为要计算高度，预留 1 秒等待文档加载
  componentDidMount() {
    setTimeout(() => {
      this.getAPs(['H1', 'H2', 'H3', 'H4']);
    }, 1000);
  }

  // 获取标题锚点：nodeArr (H1, H2, H3, H4)
  getAPs = (nodeArr) => {
    let nodeInfo = [];
    document.querySelector(".markdown-content").childNodes.forEach((item) => {
      if (nodeArr.includes(item.nodeName)) {
        nodeInfo.push({
          type: item.nodeName,          // 存储该标题的类型 H1, H2, H3, H4
          txt: item.getAttribute('id'), // 存储该标题的文本 [ps：marked解析出来的h1-h6标题会在id里填上对应的标题文本]
          offsetTop: item.offsetTop     // 存储该标题离页面顶部的距离
        });
      }
    });

    console.log(nodeInfo)
    this.setState({
      menuData: nodeInfo,
      menuState: nodeInfo[0].txt
    }, () => {
      // 检测滚动
      this.checkMenuScroll();
    });
  }

  // 检测页面滚动函数
  checkMenuScroll = () => {
    // this.scroll 为整个页面的根节点，用来监听滚动
    this.scroll.addEventListener('scroll', () => {
      let scrollTop = this.scroll.scrollTop; // 获取当前页面的滚动距离
      let menuState = this.state.menuData[0].txt; // 设置menuState对象默认值为第一个标题的文字

      // 对menuData循环检测，
      // 如果当前页面滚动距离 大于 一个标题离页面顶部 的距离，则将该标题的文字赋值给menuState，循环继续
      // 如果当前页面滚动距离 小于 一个标题离页面顶部 的距离，说明页面还没滚动到该标题的位置，当前标题尚未命中，之后的标题也不可能命中。 循环结束
      for (let item of this.state.menuData) {
        if (scrollTop >= item.offsetTop) {
          menuState = item.txt;
        } else {
          break;
        }
      }

      // 如果滑动到了页面的底部，则命中最后一个标题
      if (this.scroll.clientHeight + scrollTop === this.scroll.scrollHeight) {
        menuState = this.state.menuData[this.state.menuData.length - 1].txt;
      }

      // 如果当前命中标题和前一个命中标题的文本不一样，说明当前页面处于其他标题下的内容，切换menuState
      if (menuState !== this.state.menuState) {
        this.setState({ menuState });
      }
    });
  }

  scrollPage = (item) => {
    // 创建一个setInterval，每16ms执行一次，接近60fps
    let scrollToTop = window.setInterval(() => {
      let currentScroll = this.scroll.scrollTop;


      if (currentScroll > item.offsetTop) {
        // 当页面向上滚动时操作
        this.scroll.scrollTo(0, currentScroll - Math.ceil((currentScroll - item.offsetTop) / 5));
      } else if (currentScroll < item.offsetTop) {
        // 页面向下滚动时的操作
        if (this.scroll.clientHeight + currentScroll === this.scroll.scrollHeight) {
          // 如果已经滚动到了底部，则直接跳出
          this.setState({ menuState: item.txt });
          window.clearInterval(scrollToTop);
        } else {
          this.scroll.scrollTo(0, currentScroll + Math.ceil((item.offsetTop - currentScroll) / 5));
        }
      } else {
        window.clearInterval(scrollToTop);
      }
    }, 16);
  }

  // componentDidUpdate(prevProps, prevState) {
  //   let mountToc = false
  //   if (!isEmptyObject(this.props.article)) {
  //     if (isEmptyObject(prevProps.article)) {
  //       mountToc = true
  //     } else if (this.props.article.ID != prevProps.article.ID) {
  //       mountToc = true
  //     }
  //     if (mountToc) {
  //       vinkiToc.init()
  //     }
  //   }
  // }


  render() {
    const { classes } = this.props
    // 点击目录切换


    return (
      <div ref={(ref) => { this.scroll = ref; }}>
        <div className={classes.toc} class="toc-body">
          <ul>
            {
              this.state.menuData.map((item) => {
                return (
                  <li
                    className={`${item.type}type`}
                    key={item.txt}
                    onClick={() => { this.scrollPage(item); }}
                  >
                    <a className={this.state.menuState === item.txt ? 'on' : ''} >{item.txt}</a>
                  </li>
                );
              })
            }
          </ul>
        </div>
      </div>
    )
  }
}

const Toc = withStyles(styles)(TocComponent)

export default Toc

