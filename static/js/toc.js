(function() {
    'use strict';

    var treeBar = window.treeBar = {
        site: {
            tagName: '#wiki-content',
            mountName: '#toc-container'
        },
        className: 'treeBar',
        /* 创建 TOC DOM */
        createDom: function() {
            var tocDom = document.createElement('div');
            tocDom.className = this.className;
            tocDom.innerHTML = this.ul;
            document.querySelector(this.site.mountName).appendChild(tocDom)
        },
        init: function() {
            console.time('TreeBar');

            // 1. 获取DOM
            var hList = document.querySelector(treeBar.site.tagName)
                .querySelectorAll('h1, h2, h3, h4, h5, h6');
            // 2. 构建数据
            var tree = transformTree(Array.from(hList));
            // 3. 构建DOM
            treeBar.ul = compileList(tree);
            treeBar.createDom();
            console.timeEnd('TreeBar');
        }
    };
    treeBar.init();

    /* 解析DOM，构建树形数据 */
    function transformTree(list) {
        var result = [];
        list.reduce(function(res, cur, index, arr) {
            var prev = res[res.length - 1];
            if (compare(prev, cur)) {
                if (!prev.sub) prev.sub = [];
                prev.sub.push(cur);
                if (index === arr.length - 1) prev.sub = transformTree(prev.sub);
            } else {
                construct(res, cur);
                if (prev && prev.sub) prev.sub = transformTree(prev.sub);
            }
            return res;
        }, result);

        // 转为树形结构的条件依据
        function compare(prev, cur) {
            return prev && cur.tagName.replace(/h/i, '') > prev.tagName.replace(/h/i, '');
        }

        // 转为树形结构后的数据改造
        function construct(arr, obj) {
            arr.push({
                name: obj.innerText,
                id: obj.id = obj.innerText,
                tagName: obj.tagName
            });
        }

        return result;
    }

    /* 根据数据构建目录 */
    function compileList(tree) {
        var list = '';
        tree.forEach(function(item) {
            var ul = item.sub ? compileList(item.sub) : '';
            list +=
                `<li>
					<a href="#${item.id}" title="${item.name}">${item.name}</a>${ul}
				</li>`;
        });
        return `<ul>${list}</ul>`;
    }
})();