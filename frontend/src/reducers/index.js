const vinkiApp = (state = [], action) => {
    switch(action.type) {
        case 'SET_REPOS':
            return Object.assign({}, state, {
                repo: Object.assign({}, state.repo, {
                    repos: action.repos,
                })
            })
        case 'SET_CURRENT_REPO':
            return Object.assign({}, state, {
                repo: Object.assign({}, state.repo, {
                    currentRepo: action.currentRepo,
                })
            })
        case 'SET_TOP_TAGS':    
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    topTags: action.topTags,
                })
            })
        case 'SET_SECOND_TAGS':    
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    secondTags: action.secondTags,
                })
            })
        case 'SET_SUB_TAGS':    
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    subTags: action.subTags,
                })
            })
        case 'SET_CURRENT_TOP_TAG':
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    currentTopTag: action.currentTopTag,
                })
            })
        case 'SET_CURRENT_TAG':
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    currentTag: action.currentTag,
                })
            })
        case 'SET_ARTICLE_LIST':
            return Object.assign({}, state, {
                tag: Object.assign({}, state.tag, {
                    articleList: action.articleList,
                })
            })
        default:
            return state    
    }
}

export default vinkiApp