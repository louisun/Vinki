# Vinki

Vinki 是一款面向个人的轻量级 wiki 服务，用于快速预览和查询知识文档，有以下特点：

- **安全和便捷**。所有文档都由本地 Markdown 文件生成，文件都在本地，在安全性、可迁移性上都比「在线的第三方服务」更好。

- **高效地预览本地文档**。传统文件缺乏快速、便捷的查询和浏览方式，Vinki 旨在提供一种更优雅的方式利用自己的知识库。
- **无侵入**。Vinki 只负责文档的浏览、查询，不负责文档的编辑与管理，不对原始文件的组织形式做任何更改，用户只需要配置本地 Markdown 目录树的根路径，即可生成 wiki 仓库。

## 功能展示

- 灵活选择多级标签
- 文档预览：同标签文档列表、TOC 跳转

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-04-25-15-33-18_r100.png)

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-04-25-15-33-32_r45.png)

## Philosophy

> Vinki 源自 Wiki（维基），结合了漫画「冰海战记」[Vinland Saga](https://en.wikipedia.org/wiki/Vinland_Saga_(manga)) 的名称。
>
> UI 也是来自北欧的 [Nord](https://www.nordtheme.com/) 配色，旨在提供简洁愉快的阅读体验。

## Usage

```bash
./build.sh
```

### Docker

**一、制作镜像**：

```bash
docker build -t renzo/vinki .
```

在 Docker 环境下，需要映射目录到容器中，推荐的方法是将所有仓库目录放到同一个目录下。比如下面的有 3 个仓库 `REPO_1~3`，统一放在宿主机 `HOST_ROOT_PATH` 路径下，即 `HOST_ROOT_PATH/REPO_1~3`，接着映射 `HOST_ROOT_PATH` 至容器中的目录，推荐为 `/vinki/repository`。因此在容器中，各仓库的路径为 `/vinki/repository/REPO_1~3`。

**二、创建配置**：在宿主机的 `CONF_PATH` 目录下创建 `conf.yml` 文件， 编写内容如下：

```yaml
system:
  debug: false
  port: 6166

repositories:
  - root: "/vinki/repository/{REPO_1}"
    exclude:
      - "{YOUR_EXCLUE_DIR_1}"
      - "{YOUR_EXCLUE_DIR_2}"
  - root: "/vinki/repository/{REPO_2}"
  - root: "/vinki/repository/{REPO_3}"
```

**三、启动容器**：

```bash
docker run -d --name vinki -p 6166:6166 \
	-v {HOST_ROOT_PATH}:/vinki/repository \
	-v {CONF_PATH}:/vinki/conf \
	renzo/vinki
```