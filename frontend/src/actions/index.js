// Actions

export const setRepos = repos => {
    return {
        type: "SET_REPOS",
        repos: repos,
    }
}

export const setCurrentRepo = currentRepo => {
    return {
        type: "SET_CURRENT_REPO",
        currentRepo: currentRepo,
    }
}

export const setTopTags = topTags => {
    return {
        type: "SET_TOP_TAGS",
        topTags: topTags,
    }
}


export const setSecondTags = secondTags => {
    return {
        type: "SET_SECOND_TAGS",
        secondTags: secondTags,
    }
}

export const setSubTags = subTags => {
    return {
        type: "SET_SUB_TAGS",
        subTags: subTags,
    }
}

export const setCurrentTopTag = currentTopTag => {
    return {
        type: "SET_CURRENT_TOP_TAG",
        currentTopTag: currentTopTag,
    }
}

export const setCurrentTag = currentTag => {
    return {
        type: "SET_CURRENT_TAG",
        currentTag: currentTag,
    }
}

export const setArticleList = articleList => {
    return {
        type: "SET_ARTICLE_LIST",
        articleList: articleList,
    }
}
