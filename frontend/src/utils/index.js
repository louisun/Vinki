export const isEmptyObject = obj => {
    for (var n in obj) {
        return false
    }
    return true;
}

export const lastOfArray = arr => {
    if (arr.length > 0) {
        return arr[arr.length - 1]
    } else {
        return ""
    }
}