# cloudflare 零信任网络


### VPN基本原理
```
// 在需要通讯的两个端点之间构建虚拟的tunnel通道
内网A   <-----> proxyA | 公网  | proxyB <------> 内网B
```

### cloudflare 实现方式
```
用户设备 ---------> cloudflare服务器 ---------> 公司内网代理节点的tunnel代理客户端(负责内网出站入站流量) ---------> 实际访问节点
```

### 操作步骤
#### 步骤1: 注册 Cloudflare 账户并激活 Zero Trust Dashboard
1. 访问 Cloudflare 网站：注册一个新的 Cloudflare 账户。
2. 进入 Zero Trust 仪表板：登录后，在左侧导航栏找到 "Zero Trust" 或 "Cloudflare One" 选项并点击。
3. 设置团队名称：系统会要求您创建一个“团队名称”（例如：my-company-team）。这个名称将用于您的 Zero Trust 登录页面 URL（例如：https://my-company-team.cloudflareaccess.com）。

#### 步骤2: 添加域名解析(这里可以是子域名，比如jcd.com,内网的域名定位vpn.jcd.com),这里的域名等同于指向公司内网代理节点的tunnel代理客户端(可以通过浏览器访问内部的应用)

#### 添加身份管理,一般在access部分或者部分
1. 选择github的组织，最好根据组织 + 团队名称来授权，简单

#### 步骤3: 添加和配置隧道
1. 网络部分添加和创建隧道(在主页的部分创建隧道)
2. 内网节点上下载和安装好cloudflared（企业内网代理客户端）
3. 本机下载和安装好warp client(用户设备代理客户端), 在客户端页面选择zero trust，看看是否能够登录进去(需要输入管理员邮箱来验证设备身份,之后设备通过以后会自动跳转到github登录页面授权)
#### 步骤4
1. 在隧道处配置内网需要暴漏的服务和域名
```
https://developers.cloudflar[VPN.md](VPN.md)e.com/cloudflare-one/networks/connectors/cloudflare-tunnel/private-net/cloudflared/connect-private-hostname/
```
2. 验证暴漏的服务在非公司内网是否能访问


### 建议
1. 先大致浏览一下cloudflare的控制面板,看看控制面板那一部分的功能在什么地方
2. 整体上最重要的就是隧道的配置(网络部分)以及访问权限控制(access部分)
3. 可以在cloudflare或者AI上逐步咨询,整个流程下来还是有些复杂











               





