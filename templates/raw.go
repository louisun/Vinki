package templates

var Html = `<h1 id="docker-网络连接方式">Docker 网络连接方式</h1>

<p>Docker 默认提供了 5 种网络驱动模式：</p>

<ul>
<li><code>bridge</code>: <strong>默认的网络驱动模式</strong>。如果不指定驱动程序，bridge 便会作为默认的网络驱动模式。<strong>当应用程序运行在需要通信的独立容器 (standalone containers) 中时，通常会选择 bridge 模式</strong>。</li>
<li><code>host</code>：移除容器和 Docker 宿主机之间的网络隔离，并<strong>直接使用主机的网络</strong>。host 模式仅适用于 Docker 17.06+。</li>
<li><code>overlay</code>：overlay 网络<strong>将多个 Docker 守护进程连接在一起，并使集群服务能够相互通信</strong>。您还可以使用 overlay 网络来实现 swarm 集群和独立容器之间的通信，或者不同 Docker 守护进程上的两个独立容器之间的通信。该策略<strong>实现了在这些容器之间进行操作系统级别路由的需求</strong>。</li>
<li><code>macvlan</code>：Macvlan 网络<strong>允许为容器分配 MAC 地址，使其显示为网络上的物理设备</strong>。 Docker 守护进程通过其 MAC 地址将流量路由到容器。对于<strong>希望直连到物理网络的传统应用程序</strong>而言，使用 macvlan 模式一般是最佳选择，而不应该通过 Docker 宿主机的网络进行路由。</li>
<li><code>none</code>：对于此容器，<strong>禁用所有联网</strong>。通常与自定义网络驱动程序一起使用。none 模式不适用于集群服务。</li>
</ul>

<h2 id="bridge-桥接介绍">bridge 桥接介绍</h2>

<p>早期的二层网络中，bridge 可以连接不同的 LAN 网（局域网），如下图所示。当主机 1 发出一个数据包时，LAN 1 的其他主机和网桥 br0 都会收到该数据包。网桥再将数据包从入口端复制到其他端口上（本例中就是另外一个端口）。因此，LAN 2 上的主机也会接收到主机 A 发出的数据包，从而实现不同 LAN 网上所有主机的通信。</p>

<blockquote>
<p>此时网桥可以看做是一台「物理的交换机」，连接两个不同的网段。</p>
</blockquote>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-16-46-59_r54.png" alt="" /></p>

<h3 id="linux-bridge-模式">Linux Bridge 模式</h3>

<p>随着网络技术的发展，传统 bridge 衍生出适用不同应用场景的模式，其中最典型要属 <code>Linux Bridge</code> 模式，它是 Linux Kernel 网络模块的一个重要组成部分，用以<strong>保障不同虚拟机之间的通信，或是虚拟机与宿主机之间的通信</strong>，如下图所示 ：</p>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-16-47-42_r79.png" alt="" /></p>

<p>Linux Bridge 模式下，Linux Kernel 会创建出一个虚拟网桥 ，用以实现主机网络接口与虚拟网络接口间的通信。从功能上来看，<strong>Linux Bridge 像一台虚拟交换机</strong>，所有桥接设置的虚拟机分别连接到这个交换机的一个接口上，接口之间可以相互访问且互不干扰，这种连接方式对物理主机而言也是如此。</p>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-16-52-45_r29.png" alt="" /></p>

<p>在桥接的作用下，「<strong>虚拟网桥</strong>」会把「主机网络接口接收到的网络流量」转发给「<strong>虚拟网络接口</strong>」，于是虚拟网络连接的*虚拟机能够接收到路由器发出的 DHCP*（动态主机设定协议，用于获取局域网 IP）信息及路由更新。这样的工作流程，同样适用于不同虚拟网络接口间的通信。具体的实现方式如下所示：</p>

<ul>
<li><p><strong>虚拟机与宿主机通信</strong>： 用户可以手动为虚拟机配置 IP 地址、子网掩码，*该 IP 需要和宿主机 IP 处于同一网段*，这样虚拟机才能和宿主机进行通信。</p></li>

<li><p><strong>虚拟机与外界通信</strong>： 如果虚拟机需要联网，还需为它手动配置网关，该*网关也要和宿主机网关保持一致*。</p></li>

<li><p>除此之外，还有一种较为简单的方法，那就是虚拟机通过 DHCP 自动获取 IP，实现与宿主机或宿主机以外的世界通信。</p></li>
</ul>

<h3 id="docker-bridge-模式">Docker Bridge 模式</h3>

<h4 id="veth-pair">Veth Pair</h4>

<p>在该模式下，Docker Daemon 会创建一个名为「<code>docker0</code>」的虚拟网桥，用来连接宿主机和容器，或者连接不同的容器。</p>

<p>Docker 利用 <code>veth pair</code> 技术，在宿主机上创建了两个虚拟网络接口 <code>veth0</code> 和 <code>veth1</code>。</p>

<blockquote>
<p>veth pair 是用于不同 network namespace 间进行通信的方式，veth pair 将一个 network namespace 数据发往另一个 network namespace 的 veth。</p>

<p>veth pair 技术的特性可以保证无论哪一个 veth 接收到网络报文，都会无条件地传输给另一方。</p>
</blockquote>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-16-57-16_r90.png" alt="" /></p>

<p><strong>容器与宿主机通信</strong> : 在桥接模式下，Docker Daemon 将 <code>veth0</code> 附加到 <code>docker0</code> 网桥上，保证宿主机的报文有能力发往 <code>veth0</code>。再将 <code>veth1</code> 添加到「<strong>Docker 容器所属的网络命名空间</strong>」，保证宿主机的网络报文若发往 <code>veth0</code> 可以立即被 <code>veth1</code> 收到。</p>

<blockquote>
<p>网络命名空间是用于隔离网络资源（/proc/net、IP 地址、网卡、路由等）。由于一个物理的网络设备最多存放在一个网络命名空间中，所以通过 veth pair 在不同的网络命名空间中创建通道，才能达到通信的目的。</p>
</blockquote>

<p><strong>容器与外界通信</strong> : 容器如果需要联网，则需要采用 <code>NAT</code> （<strong>网络地址转换</strong>）方式。准确的说，是 NATP (<strong>网络地址端口转换</strong>) 方式。NATP 包含两种转换方式：源 NAT 和 目的 NAT 。</p>

<h4 id="dnat">DNAT</h4>

<p><strong>目的 NAT (DNAT):</strong> 修改数据包的目的地址。<strong>当宿主机以外的世界需要访问容器时</strong>，数据包的流向如下图所示：
<img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-17-01-41_r65.png" alt="" /></p>

<p>由于容器的 IP 与端口对外都是不可见的，所以数据包的*目的地址为宿主机的 ip 和端口*，为 <code>192.168.1.10:24</code> 。数据包经过路由器发给宿主机 <code>eth0</code>，再经 <code>eth0</code> 转发给 <code>docker0</code> 网桥。由于存在「<strong>DNAT 规则</strong>」，会将数据包的*目的地址转换为容器的 ip 和端口*，为 <code>172.17.0.n:24</code> 。</p>

<p>宿主机上的 <code>docker0</code> 网桥识别到*容器 ip 和端口*，于是将数据包发送附加到 <code>docker0</code> 网桥上的 <code>veth0</code> 接口，<code>veth0</code> 接口再将数据包发送给容器内部的 <code>veth1</code> 接口，容器接收数据包并作出响应。</p>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-17-04-10_r7.png" alt="" /></p>

<h4 id="snat">SNAT</h4>

<p><strong>源 NAT (SNAT</strong>): 修改数据包的源地址。<strong>当容器需要访问宿主机以外的世界时</strong>，数据包的流向为下图所示：</p>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-17-06-28_r39.png" alt="" /></p>

<p>此时数据包的*源地址为容器的 ip 和端口*，为 <code>172.17.0.n:24</code>，容器内部的 <code>veth1</code> 接口将数据包发往 <code>veth0</code> 接口，到达 <code>docker0</code> 网桥。</p>

<p>宿主机上的 <code>docker0</code> 网桥发现数据包的*目的地址为外界的 IP 和端口*，便会将数据包转发给 <code>eth0</code> ，并从 <code>eth0</code> 发出去。由于存在「<strong>SNAT 规则</strong>」，会将数据包的*源地址转换为宿主机的 ip 和端口*，为 <code>192.168.1.10:24</code> 。</p>

<p>由于路由器可以识别到*宿主机的 ip 地址*，所以再将数据包转发给外界，外界接受数据包并作出响应。这时候，在外界看来，这个数据包就是从 <code>192.168.1.10:24</code> 上发出来的，Docker 容器对外是不可见的。</p>

<p><img src="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-02-09-17-06-33_r95.png" alt="" /></p>

<h2 id="docker-bridge-实战">Docker Bridge 实战</h2>

<pre><code class="language-bash"># docker 默认有 3 种网络：bridge host none
NETWORK ID          NAME                DRIVER              SCOPE
1efbc6e45580        bridge              bridge              local
2346272d23c1        host                host                local
a366ee8fdccb        none                null                local

# 运行容器时，可以使用 –network 参数指定在哪种网络模式下运行该容器
docker run --it --network bridge ubuntu /bin/bash
</code></pre>

<p>所有 Docker 安装后都存在的 <code>docker0</code> 网络。除非使用 <code>docker run –network=</code> 选项另行指定，否则 Docker 守护进程默认情况下会将容器连接到 <code>docker0</code> 这个网络。</p>

<h3 id="创建自定义网络">创建自定义网络</h3>

<pre><code class="language-bash"># 网络驱动模式默认为 bridge
$ docker network create my-net

$ docker network ls
NETWORK ID          NAME                DRIVER              SCOPE
1efbc6e45580        bridge              bridge              local
2346272d23c1        host                host                local
b0508c2abaad        my-net              bridge              local
a366ee8fdccb        none                null                local
</code></pre>

<h3 id="用-busybox-测试容器连通性">用 busybox 测试容器连通性</h3>

<blockquote>
<p>BusyBox 是一个集成了一百多个最常用 Linux 命令和工具（如 cat、echo、grep、mount、telnet 、ping、ifconfig 等）的精简工具箱，它只需要几 MB 的大小，很方便进行各种快速验证，被誉为“Linux 系统的瑞士军刀”。</p>
</blockquote>

<h4 id="使用默认网桥-docke0">使用默认网桥 docke0</h4>

<pre><code class="language-bash">$ docker run -it --rm --name box1 busybox sh
# 容器中操作
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 02:42:AC:11:00:07
          inet addr:172.17.0.7  Bcast:172.17.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:6 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:508 (508.0 B)  TX bytes:0 (0.0 B)

lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
# 得知 box1 的 eth0 IP 为 172.17.0.7

</code></pre>

<p>重新操作一遍，运行 box2 容器：</p>

<pre><code class="language-bash">$ docker run -it --rm --name box2 busybox sh
# 容器中操作
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 02:42:AC:11:00:08
          inet addr:172.17.0.8  Bcast:172.17.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:6 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:508 (508.0 B)  TX bytes:0 (0.0 B)

lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
# 得知 box2 的 eth0 IP 为 172.17.0.8
</code></pre>

<p>在 box2 中 ping box1：</p>

<pre><code class="language-bash">$ ping 172.17.0.7
PING 172.17.0.7 (172.17.0.7): 56 data bytes
64 bytes from 172.17.0.7: seq=0 ttl=64 time=9.670 ms
64 bytes from 172.17.0.7: seq=1 ttl=64 time=0.210 ms
64 bytes from 172.17.0.7: seq=2 ttl=64 time=0.148 ms
64 bytes from 172.17.0.7: seq=3 ttl=64 time=0.109 ms
64 bytes from 172.17.0.7: seq=4 ttl=64 time=0.301 ms

$ ping box1
# 无响应
</code></pre>

<p>发现使用默认网桥 <code>docker0</code> 的桥接模式下，ip 是通的，但是无法使用容器名作为通信的 host。</p>

<h4 id="使用自定义网桥-my-net">使用自定义网桥 my-net</h4>

<pre><code class="language-bash">$ docker run -it --rm --name box3 --network my-net busybox sh
# 容器中操作
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 02:42:AC:12:00:02
          inet addr:172.18.0.2  Bcast:172.18.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:14 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:1156 (1.1 KiB)  TX bytes:0 (0.0 B)

lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
# 得知 box3 的 eth0 IP 为 172.18.0.2
</code></pre>

<pre><code class="language-bash">$ docker run -it --rm --name box4 --network my-net busybox sh
# 容器中操作
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 02:42:AC:12:00:03
          inet addr:172.18.0.3  Bcast:172.18.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:5 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:418 (418.0 B)  TX bytes:0 (0.0 B)

lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
# 得知 box4 的 eth0 IP 为 172.18.0.3
</code></pre>

<p>在 box4 中 ping box3：</p>

<pre><code class="language-bash">$ ping 172.18.0.2
PING 172.18.0.2 (172.18.0.2): 56 data bytes
64 bytes from 172.18.0.2: seq=0 ttl=64 time=0.234 ms
64 bytes from 172.18.0.2: seq=1 ttl=64 time=0.118 ms
64 bytes from 172.18.0.2: seq=2 ttl=64 time=0.097 ms
64 bytes from 172.18.0.2: seq=3 ttl=64 time=0.167 ms
64 bytes from 172.18.0.2: seq=4 ttl=64 time=0.118 ms
64 bytes from 172.18.0.2: seq=5 ttl=64 time=0.429 ms

$ ping box3
PING box3 (172.18.0.2): 56 data bytes
64 bytes from 172.18.0.2: seq=0 ttl=64 time=0.162 ms
64 bytes from 172.18.0.2: seq=1 ttl=64 time=0.133 ms
64 bytes from 172.18.0.2: seq=2 ttl=64 time=0.203 ms
64 bytes from 172.18.0.2: seq=3 ttl=64 time=0.130 ms
# box4 和 box3 容器都用了自定义的 netwrok：my-network，可以通过容器名相互通信
</code></pre>

<p>与默认的网络 docker0 不同的是，<strong>指定了自定义 network 的容器可以使用容器名称相互通信</strong>，实际上这也是 docker 官方推荐使用 <code>—network</code> 参数运行容器的原因之一。</p>

<h3 id="使用自定义-bridge-的好处">使用自定义 bridge 的好处</h3>

<h4 id="更好的隔离性和容器间的互操作性">更好的隔离性和容器间的互操作性</h4>

<p><strong>连接到同一个自定义 bridge 网络的容器会自动将所有端口相互暴露，并且无法连接到容器之外的网络</strong>。这使得容器化的应用能轻松地相互通信，并且<strong>与外部环境产生了良好的隔离性</strong>。</p>

<p>例如一个包含了 web 应用，数据库，redis 等组件的应用程序。很有可能*只希望对外界暴露 80 端口*，而不允许外界访问数据库端口和 redis 端口，而又不至于让 web 应用本身无法访问数据库和 redis， 便可以使用自定义 bridge 网络轻松实现。</p>

<p>如果在默认 bridge 网络上运行相同的应用程序，则需要使用 <code>-p</code> 或 <code>--publish</code> 标志打开 web 端口，数据库端口，redis 端口。这意味着 Docker 宿主机需要通过其他方式阻止对数据库端口，redis 端口的访问，无意增大了工作量。</p>

<h4 id="容器间的自动-dns-解析">容器间的自动 DNS 解析</h4>

<p>这一点在上一节的实验中已经验证过了。<strong>默认 bridge 网络上的容器只能通过 IP 地址互相访问，除非使用在 <code>docker run</code> 时添加 <code>—link</code> 参数</strong>。这么做个人认为有两点不好的地方：</p>

<ol>
<li><p>容器关系只要稍微复杂一些，便会对管理产生不便。</p></li>

<li><p><code>—link</code> 参数在官方文档中已经被标记为过期的参数，不被建议使用。</p></li>
</ol>

<p>在用户定义的桥接网络上，容器可以通过容器名称 (<code>--name</code> 指定的名称) 或别名来解析对方。可能有人说，在默认 bridge 模式下我可以去修改 <code>/etc/hosts</code> 文件呀，但这显然不是合理的做法。</p>

<h4 id="容器可以在运行中与自定义-bridge-网络连接和分离">容器可以在运行中与自定义 bridge 网络连接和分离</h4>

<p>在容器的生命周期中，<strong>可以在运行中将其与自定义网络连接或断开连接</strong>。 而要从默认 bridge 网络中移除容器，则需要停止容器并使用不同的网络选项重新创建容器。</p>

<h4 id="每个自定义的-bridge-网络都会创建一个可配置的网桥">每个自定义的 bridge 网络都会创建一个可配置的网桥</h4>

<p>如果容器使用默认 bridge 网络，虽然可以对其进行配置，但所有容器都使用相同的默认设置，例如 MTU 和防火墙规则。另外，配置默认 bridge 网络隔离于 Docker 本身之外，并且需要重新启动 Docker 才可以生效。</p>

<p>自定义的 bridge 是使用 <code>docker network create</code> 创建和配置的。如果<strong>不同的应用程序组具有不同的网络要求，则可以在创建时分别配置每个用户定义的 bridge 网络，这无疑增加了灵活性和可控性</strong>。</p>

<h4 id="使用默认-bridge-容器共享所有的环境变量">使用默认 bridge 容器共享所有的环境变量</h4>

<p>在 Docker 的旧版本中，两个容器之间<strong>共享环境变量的唯一方法是使用 <code>—link</code> 标志来进行链接</strong>。</p>

<p><strong>这种类型的变量共享对于自定义的网络是不存在的</strong>。但是，自定义网络有更好方式来实现共享环境变量：</p>

<ul>
<li>多个容器可以<strong>使用 Docker 卷来挂载包含共享信息的文件或目录</strong>。</li>
<li>多个容器可以使用 <code>docker-compose</code> 一起启动，并且 <code>docker-compose.yml</code> 文件可以定义共享变量。</li>
<li>使用集群服务而不是独立容器，并利用<strong>共享密钥和配置。</strong></li>
</ul>

<p>结合上述这些论述和官方文档的建议，使用 bridge 网络驱动模式时，最好添加使用 <code>—network</code> 来指定自定义的网络。</p> `
