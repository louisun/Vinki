# Docker 技术入门与实战

## 1. 初识 Docker 与容器

### 介绍

如果说主机时代比拼的是单个服务器物理性能（如 CPU 主频和内存）的强弱，那么**在云时代，最为看重的则是凭借「虚拟化技术」所构建的「集群处理能力」** 。

Docker 的构想：「**Build, Ship and Run Any App, Anywhere.**」

Docker 通过对应用的「封装 Packaging」、「分发 Distribution」、「部署 Deployment」、「运行 Runtime」生命周期进行管理，达到应用组件级别的「**一次封装，到处运行**」。这里的应用组件，既可以是一个 Web 应用、一个编译环境，也可以是一套数据库平台服务，甚至是一个操作系统或集群。

**Linux 容器**（Linux Containers, `LXC`）技术是「巨人的肩膀」。容器可以在核心 CPU 本地运行指令。

早期 Docker 代码实现是直接基于 LXC，`0.9` 版本后用 Golang 开发了[libcontainer](https://github.com/docker/libcontainer)，作为更广泛的容器驱动实现，替代掉了 LXC 的实现。目前 Docker 推动了 [runc](https://github.com/opencontainers/runc) 项目，它是一个 CLI 工具，用于根据 OCI 规范生成和运行容器。

### 容器化的好处

之前构建应用、部署调试不方便的根源：无法保证同一份应用在不同的运行环境中行为一致。

Docker 通过容器来打包应用、解耦应用和运行平台。

Docker 在开发和运维的优势：

- 更快的交付和部署（开发、测试、运维可以用同一套环境）
- 更高效的资源利用（Docker 是内核级虚拟化）
- 更轻松的迁移和扩展
- 更简单的更新管理（Dockerfile）

## 2. Docker 仓库

### 搭建私有仓库

搭建本地私有仓库，监听端口为 5000：

```bash
docker run -d -p 5000:5000 registry:2
```

默认仓库目录为容器 `/var/lib/registry`，可通过容器卷映射在宿主机上保存镜像目录：

```bash
docker run -d -p 5000:5000 -v /opt/data/registry:/var/lib/registry registry:2
```

### 管理私有仓库

上传镜像到私有仓库，必须标记此镜像为私有仓库的镜像。

若当前仓库的 IP 和端口为 `10.0.2.2:5000`，用 `docker tag` 标记镜像：

```bash
# docker tag IMAGE[:TAG] [REGISTRYHOST/][USERNAME/]NAME[:TAG]
docker tag ubuntu:18.04 10.0.2.2:5000/test
```

**使用 `docker push` 上传标记的镜像：**

```
docker push 10.0.2.2:5000/test
```

**使用 `curl` 查看仓库的镜像**：

```bash
curl http://10.0.2.2:5000/v2/search
```

现在任意一台能访问 10.2.2.2 的机器都能下载这个镜像

新版本 Docker 为了安全性会要求仓库支持 SSL/TLS 证书，私有仓库可自行配置证书或关闭对仓库的安全性检查：

```bash
DOCKER_OPTS="--insecure-registry 10.0.2.2:5000"
```

重启服务，即可拉取镜像：

```bash
sudo service docker restart

docker pull 10.0.2.2:5000/test
```

下载后，可以添加一个更通用的标签 `ubuntu:18.04`，方便后续使用：

```bash
docker tag 10.0.2.2:5000/test ubuntu:18.04
```

## 3. 数据管理

在生产环境中用 Docker，往往需要对数据进行持久化，或者多个容器之间进行数据共享。

容器中管理数据主要有 2 种方式：

- 数据卷 Data Volumes：容器内数据直接映射至宿主机本地环境
- 数据卷容器 Data Volumes Containers：使用特定容器维持数据卷

### 数据卷

数据卷是一个可供容器使用的「**特殊目录**」，将主机操作系统目录直接映射进容器，类似于 Linux 中的 mount 行为。

```bash
# 本地创建一个数据卷
docker volume create -d local test

# 查看 /var/lib/docker/volumes 路径，会找到 test 目录
ls -l /var/lib/docker/volumes
```

其余命令：

- `docker volumes inpect`：查看详细信息
- `docker volumes ls`：列出已有数据卷
- `docker volumes prune`：清除无效数据卷
- `docker rm`：删除数据卷

除了使用 volume 命令来管理数据卷，还可以在创建容器时，将本地任意的路径挂载到容器内作为数据卷，这种形式叫「**绑定数据卷**」。

说白了就是 `-v` 参数，或者 `--mount` 参数

Docker 挂载数据卷的默认权限是「**读写** `rw`」，可以自己指定如 `ro` 为只读：

```bash
# -P 是随机映射端口
docker run -d -P --name web -v /webapp:/opt/webapp:ro traning/webapp python app.py
```

> 如果直接用宿主机中用 mount 的方式挂载一个文件到容器，使用文件编辑工具如 vim 或 sed 的时候，可能会造成文件 inode 的改变，从 Docker 1.1.0 起这会报告错误信息。

### 数据卷容器

如果用户需要在多个容器之间共享一些持续更新的数据，最简单的方式是使用「**数据卷容器**」，它也是一个容器，但目的是专门提供数据给其他容器挂载。

现在创建一个「**数据卷容器**」`dbdata`，并在其中创建一个数据卷挂载到 `/dbdata`：

```bash
# 宿主机会在 /var/lib/docker/volumes 下生成一个随机目录，映射到容器的 /dbdata
docker run -it -v /dbdata --name dbdata ubuntu
```

在其他容器中，使用 `--volumes-from` 可以挂载另一个容器中的数据卷，比如：

```bash
docker run -it --volumes-from dbdata --name db1 ubuntu
docker run -it --volumes-from dbdata --name db2 ubuntu
docker run -d --volumes-from db1 --name db3 training/postgres
```

> `--volumes-from` 参数所挂载的数据卷容器，并不需要保持在运行状态
>
> 如果删除了挂载的容器，数据卷并不会自动删除，要删除一个数据卷，必须在删除最后一个挂载它的容器时，显式用 `docker rm -v` 指定同时删除关联的容器。

使用「数据卷容器」可以让用户在容器间自由地升级和移动数据卷。

**备份数据卷容器内的数据卷**：

```bash
docker run --volumes-from dbdata 
           -v ${pwd}:/backup \
           --name worker \
           ubuntu \
           tar cvf /backup/backup.tar /dbdata
```

运行一个名为 worker 的容器（执行 tar 打包命令到指定的目录）

**恢复数据到一个容器**：

```bash
# 创建一个带有数据卷的容器 dbdata2
docker run -v /dbdata --name dbdata2 ubuntu /bin/bash

# 创建一个新的容器，挂载 dbdata2 的容器
# 将 /backup 目录（宿主机当前目录）下的 tar 备份文件恢复到容器中
docker run --volumes-from dbdata2
           -v ${pwd}:/backup \
           busybox tar xvf /backup/backup.tar
```

## 4. 端口映射与容器互联

除了通过网络访问外，Docker 还提供了 2 各方便的功能满足服务之间相互访问的基本需求：

- 映射容器内应用的服务端口到宿主机
- 互联机制：试下多个容器间通过容器名来快速访问

### 端口映射

`-P` 参数 Docker 会随机映射一个「49000~49900」的端口到「内部容器开放的网络端口」:

```bash
docker run -d -P training/webapp python app.py
# 0.0.0.0:491155->5000/tcp
```

`-p` 就是普通的端口映射：

- 映射所有接口地址：`5000:5000` 
- 映射到指定地址的指定端口：`127.0.0.1:5000/5000`
- 映射到指定地址的任意端口：`127.0.0.1::5000`（本地主机会自动分配一个端口）
- 使用 tcp/udp 标志来标记映射 tcp 或 udp 端口：`127.0.0.1:5000:5000/udp`

查看映射端口配置：

```bash
docker port CONTAINER 5000
127.0.0.1:49155.
```

### 容器互联 linking

linking 是一种让多个容器的应用进行快速交互的方式，它会在源和接收容器之间创立连接关系，接收容器可以通过容器名快速访问源容器，而不用指定具体的 IP 地址。

**一、自定义容器命名**

连接系统根据容器的「**名称**」来执行，为了方便，最好用 `--name` 命名容器：

```bash
docker run -d -P --name web training/webapp python app.py
```

**二、容器互联**

创建一个新的数据库容器：

```bash
docker run -d --name db training/postgres
```

删除之前的 web 容器，新创建一个 web 容器，连接到 db 容器：

```bash
docker rm -f web

# 连接方式 --link name:alias
# name 是要链接的容器的名称
# alias 是别名
docker run -d -P --name web --link db:db \
	training/webapp python app.pydocker run 
```

```bash
docker ps
IMAGE                      ...  NAMES
training/postgres:latest        db, web/db
training/webapp:latest          web
```

`NAMES` 列有「db」和「web/db」，表示 web 容器链接到 db 容器，web 容器被允许访问 db 容器的信息。Docker 相当于在两个互联的容器之间创建了一个「**虚拟通道**」，而且**不用映射它们的端口到宿主机上**。在启动 db 的时候并没有用 `-p` 或 `-P` 标记，从而避免了暴露数据库服务端口到外部网络上。

Docker 通过 2 种方式为容器公开连接信息：

- 更新环境变量
- 更新 `/etc/hosts` 文件

```bash
# 查看 web2 容器内的环境变量
docker run --rm --name web2 --link db:db training/webapp env
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-22-38-51_r87.png)

其中 `DB_` 开头的环境变量是供 web 容器连接 db 容器使用，前缀采用大写的连接别名。

```bash
# 查看 web 容器内的 /etc/hosts
docker run -t -i --rm --link db:db training/webapp /bin/bash
$ cat /etc/hosts 
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-22-40-38_r94.png)

这两个 hosts 信息，分别表示 web 容器自己和 db 容器的 IP 和主机名（web 容器默认主机名是自己的 ID），在 web 容器中安装 ping 命令来测试跟 db 容器的连通：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-22-42-21_r2.png)

## 5. 操作系统

### BusyBox

> `BusyBox` 是一个集成了一百多个最常用 Linux 命令（如 cat、echo、grep、mount、telnet 等）的精简工具箱，它只有不到 `2MB` 大小，被誉为「Linux 系统的瑞士军刀」。 BusyBox 可运行于多款 POSIX 环境的操作系统中，如 Linux（包括 Android）、Hurd 、FreeBSD 等。

```bash
docker pull busybox:latest

docker run -it busybox
# 进入容器就可以运行大量的命令
```

### Alpine

> `Alpine` 是一个「**面向安全**」的「**轻型**」Linux 发行版，关注安全、性能、资源效能。
>
> 口号是：Small. Simple. Secure.
>
> ![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-01-09-00-26-00_r57.png)
>
> 不同于其他发行版，Alpine 采用了 `musl libc` 和 `BusyBox` 以减小系统的体积和运行时小号，比 `BusyBox` 功能更完善。在保持瘦身的同时，还提供了包管理工具 `apk`，方便查询和安装软件包。
>
> `Alpine` Docker 镜像继承了该发行版的这些优势，镜像大小只有 `5MB`（Ubuntu 系列镜像有 `200MB`）。目前 Docker 官方推荐使用 `Alpine` 作为默认基础镜像环境。
>
> ![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-01-09-00-31-30_r12.png)

**TODO**：Alpine 制作 Go 镜像

### Debian / Ubuntu

 Debian 和 Ubuntu 比较适合研发场景，各大容器云服务都提供了完整的支持。

### CentOS / Fedora

CentOS 和 Fedora 都是基于 Redhat 的 Linux 发行版。前者是企业级服务器常用的操作系统，后者主要面向个人桌面用户。

## 6. 为镜像添加 SSH 服务

当需要远程登录到容器内进行操作，需要支持 SSH 的支持，有两种方法：

- 基于 `commit` 命令创建
- 基于 `Dockerfile` 创建

### 基于 commit 命令创建

`docker commit CONTAINER [REPOSITORY[:TAG]]` 可以对容器进行定制的修改，并生成新的镜像。

```bash
docker pull ubuntu:18.04

# 进入容器 root 用户
docker -it ubuntu:18.04 bash

# 配置软件源
apt-get update

vi /etc/apt/sources/list.d/aliyun.list

# 内容略... 添加源

apt-get update

# 安装 ssh 服务
apt-get install openssh-server

# 启动 ssh 服务要存在 /var/run/sshd
mkdir -p /var/run/sshd
/usr/sbin/sshd -D &

# 查看 22 端口是否启用
netstat -tunlp 

# 修改 ssh 服务安全登录配置，取消 pam 登录限制
# 把 /etc/pam.d/sshd 的 session required pam_loginuid.so 注释掉

mkdir root/.ssh
# 添加本地 ssh 公钥到 /root/.ssh/authorized_keys
vi /root/.ssh/authorized_keys

# 创建一个脚本，用于自启动 sshd
vi /run.sh
chmod +x run.sh

# 内容为
#!/bin/bash
/usr/sbin/sshd -D

# 退出容器
exit
```

保存镜像，名为 ` sshd:ubuntu `：

```shell
docker commit <容器名> sshd:ubuntu 
```

使用镜像：

```shell
docker run -p 10022:22 -d sshd:ubuntu /run.sh
```

其他主机 ssh 登录容器：

```shell
ssh 192.168.1.200 -p 10022
```

### 使用 Dockerfile 创建

```shell
mkdir sshd_ubuntu

cd sshd_ubuntu

# 创建 Dockerfile 和 run.sh 脚本
touch Dockerfile run.sh

# run.sh 内容与上面一样
#!/bin/bash
/usr/sbin/sshd -D
```

宿主机上生成 SSH 秘钥对，创建 `authorized_keys` 文件：

```shell
ssh-keygen -t rsa
cat ~/.ssh/id_rsa.pub > authorized_keys
```

编写 Dockerfile：

```dockerfile
FROM ubuntu:18.04

MAINTAINER docker_user (user@docker.com)

# 换源命令
RUN echo "deb http://mirrors.163.com/ubuntu/ bionic main restricted universe multiverse" > /etc/apt/source.list
... # 略
RUN apt-get update

# 安装 ssh 服务
RUN apt-get install -y openssh-server
RUN mkdir -p /var/run/sshd
RUN mkdir -p /root/.ssh
RUN sed -ri 's/session   required    pam_loginuid.so/#session   required    pam_loginuid.so/g' /etc/pam.d/sshd

# 复制配置文件到相应位置，赋予脚本可执行权限
ADD authorized_keys /root/.ssh/authorized_keys
ADD run.sh /run.sh
RUN chmod 755 /run.sh

# 开放端口
EXPOSE 22

# 设置自启动命令
CMD ["/run.sh"]
```

创建镜像：

```shell
cd sshd_ubuntu

docker build -t sshd:dockerfile .
```

使用镜像：

```shell
docker run -d -p 10122:22 sshd:dockerfile
```

其他主机 ssh 登录容器：

```shell
ssh 192.168.1.200 -p 10022
```



## 7. Docker 三剑客介绍

### Machine

Machine 负责使用使用 Docker 容器的第一步：在多种平台上快速安装和维护 Docker 运行环境，在 短时间在本地货云环境中搭建一套 Docker 主机集群。

Machine 连接不同类型的操作平台是通过对应驱动来实现的，目标集成了包括 AWS、IBM、Google、OpenStack、VirtualBox、vSphere 等多种云平台支持，此外还支持 Microsoft Hyper-V 虚拟化平台。

**1. 虚拟机**

```bash
# 启动一个全新的虚拟机，并安装 Docker 引擎
docker-machine create --driver=virtualbox test、
# 查看访问所创建 Docker 环境所需要的配置信息
docker-machine env
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-05-11_r96.png)

```bash
# 停止 Docker 主机
docker stop test
```

**2. 本地主机**

这种驱动（`generic`）适合主机可以通过 user 账号的 key 直接 ssh 到目标主机

```bash
# 要确保本地主机可以通过 user 账号的 key 直接 ssh 到目标主机
docker-machine create -d generic --generic-ip-address=10.0.100.102 --generic-ssh-user=user test
```

Machine 通过 SSH 连接到指定节点，并在上面安装 Docker 引擎。

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-09-06_r44.png)

**3. 云平台驱动**

```bash
docker-machine create -d --driver amazonec2 ... # 略
```

**4. 客户端配置**

默认情况下，所有配置数据都会自动存放在 `~/.docker/machine/machines` 路径下，用户应当定期备份。

> machine 具体命令略

### Compose

编排（Orchestration）功能，是复杂系统是否具有灵活可操作性的关键。在 Docker 应用场景中，编排意味着用户可以灵活地对各种容器资源实现定义和管理。Docker Compose 是 Docker 官方的编排工具。

Docker Compose 是 Python 实现的，可 pip 安装或下二进制文件，安装之前要先装 Docker，

其实 `docker-compose` 文件是个 `run.sh` Shell 脚本，实际执行的时候又是通过运行「`docker/compose`」镜像起作用的：

```bash
# 一系列变量配置
exec docker run --rm $DOCKER_RUN_OPTIONS $DOCKER_ADDR $COMPOSE_OPTIONS $VOLUMES -w "$(pwd)" $IMAGE "$@"
```

模板文件是 Compose 的核心，默认名称是 `docker-compose.yml`，具体见 [Docker Compose.md](./Docker Compose.md)

### Swarm

Swarm 提供 Docker **容器集群服务**，是 Docker 官方对容器云生态进行支持的核心方案。

使用 Docker Swarm ，用户可将多个 Docker 主机抽象为大规模的虚拟 Docker 服务，快速打造一套容器云平台。

Swarm 最大的优势之一是**支持原生 Docker API**，给用户带来极大的便利，轻松集成各种标准 API 工具如 Compose 等等，**方便用户将原先单节点的系统移植到 Swarm 上**。同时 Swarm 内置了对 Docker 网络插件的支持，用户可以**轻松地部署跨主机的容器集群服务**。

Swarm 采用了主从结构，通过 Raft 协议在多个管理节点中实现公式。工作节点上运行 agent 接收管理节点的统一管理和任务分配。用户提交服务请求只要发给管理节点即可，管理节点会按照调度策略在集群中分配节点来运行服务相关的任务。

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-22-18_r4.png)

Swarm V2 中，集群会通过 Raft 协议选出 Manager 节点，无需额外的发现服务支持，避免了单点瓶颈。V2 中内置了基于 DNS 的负载均衡和对外负载均衡机制的集成支持。



Swarm 集群的主要操作：

- `swarm init`：在管理节点上创建一个集群
- `node list`：列出集群中节点信息
- `swarm join`：加入一个新的节点到已有集群中
- `swarm update`：更新一个 Swarm 集群
- `swarm leave`：离开一个 Swarm 集群

**1. 创建集群**：

```bash
docker swarm init --advertise-addr <manager ip>
# 返回 token 串：集群唯一 ID，此外还有 ip:port 
```

**2. 查看集群信息**：

```bash
docker node ls
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-28-29_r8.png)

**3. 加入集群**

在工作节点上，执行加入到集群：

```bash
docker swarm join --token <token> <manager ip>:<port>
```

再到管理机上查看集群节点情况：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-30-02_r38.png)

**4. 使用集群服务**

在管理节点上执行命令，创建一个应用服务，并制定服务的复制份数为 2：

```bash
docker service create --replicas 2 --name ping_app debian:jessie ping docker.com
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-33-08_r6.png)

管理节点查看服务：

```bash
docker service ls
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-33-30_r39.png)

```bash
docker service inspect --pretty ping_app
# 详细信息

docker service ps ping_app
```

发现管理节点和工作节点上，都运行了一个容器：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-34-39_r57.png)

扩展服务：

```bash
docker service scale <service-id>=<number-of-tasks>

# 副本数从 2 变为 1
docker service scale ping_app=1

# 删除服务
docker service rm <service-id>
```

使用外部服务地址：

Swarm 通过路由机制支持服务对外映射指定端口，该端口可以在集群中「**任意节点**」上机械能访问，即使该节点上没有运行服务实例。

```bash
docker service create \
	--name <service-name>
	--publish published=<pub port>,target=<container port> \
	<IMAGE>
```

用户访问集群中任意节点，都会被 Swarm 的负载均衡器代理到对应的服务实例。用户也可以配置独立的负载均衡服务，后端指向集群中各个节点对应的外部端口，获取高可用性。

**5. 更新集群**

```bash
docker swarm update [OPTIONS]
```

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-39-48_r24.png)

**6. 离开集群**

```bash
docker swarm leave [OPTIONS]
# -f --force 强制离开
```

## 8. Mesos

Mesos 是 UC Berkeley 对集群资源进行抽象和管理的开源项目，类似于操作系统内核，使用它可以实现分布式应用的自动化调度。

Mesos 主要由 C++ 编写，它要做的事情，其实就是今天操作系统内核的职责：**抽象资源 + 调度任务**。将整个数据中心的资源（CPU、内存、存储、网络等）进行抽象和调度，使得多个应用同时运行在集群中分享资源，无需关心资源的物理分布情况。

Mesos 自身知识一个「**资源抽象的平台**」，需要结合运行其上的分布式应用才能发挥作用。

Mesos 基本架构：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-08-23-45-03_r32.png)

> 一般说来，如果只是用于容器集群管理，Kubernetes 更加合适，如果定制需求比较多，或者要搭建大数据平台，架构相对松耦合的 Mesos 显然更加合适。
>
> 从软件设计初衷来看，Kubernetes 希望成为容器管理领域的领导者，而 AWS、Azure 加入 CNCF、Docker 官方表态原生支持 Kubernetes ，说明 Kubernetes 凭借源自 Google 的优秀设计，在容器领域的地位已经不可动摇，社区和生态越来越繁荣。
>
> Mesos 的目标则是资源共享，可以让企业把已经存在的业务负载，比如 Hadoop、Spark，放到一个共同管理的环境。

## 9. Kubernetes

略啦