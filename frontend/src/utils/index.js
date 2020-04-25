export function isEmptyObject(obj) {
    for (var n in obj) {
        return false
    }
    return true;
}