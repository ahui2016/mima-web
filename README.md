# mima-web

- A password manage using a web server.

本软件是一个有特色的密码管理器，需要由用户自行假设到一台 Linux 服务器中。
采用简单有效的 NaCl (libsodium) 加密方式（该方式以 “容易正确处理” 为最大特点），
因此可以确保加密过程得到正确处理。

自己架设密码管理软件的好处是，拥有数据的绝对控制权。

## Install

1. `git clone https://github.com/ahui2016/mima-web.git`
2. mkdir ~/MimaWebDB
3. cd mima-web && go build
4. go run cmd/create-account.go -password yourpassword
5. ./mima-web

Done.

## 全原生代码

- 后端采用 Go 语言，未使用框架，由标准库搭建
- 不使用数据库，既免除了数据库相关设置的麻烦，又能减轻服务器负担
  (不使用 sqlite, 因此编译速度不会被 cgo 拖慢)
- 前端没有 React、Vue，没有 JQuery，没有 Bootstrap，只有原生 HTML, JS, CSS
- 页面非常简陋（其中 CSS 只有约 50 行，可见一斑）
- 前后端分离，前端可根据个人需要用各种框架/库改写，不需要修改后端

## 解决什么问题

我有一台服务器，我想在上面安装一个密码管理工具，我希望它小巧，占用资源少，
随时可以访问，我不需要指纹解锁、自动填写等的高级功能，希望它是简单的，
使用起来的 “手感” 接近命令行。

比如，假设有一个程序叫做 mypass，在命令行执行 `mypass v2` 就能复制我在
v2ex.com 的用户名到剪贴板，执行 `mypass -p v2` 就能复制密码到剪贴板。
（其中 v2 是 v2ex.com 的别名，由用户设定）

但是为了方便随时随地使用，显然做成一个网站比命令行更合适，因此我做了一个简单的网站，
并且让它的使用 “手感” 接近命令行。

## 命令行手感

- 打开网站，只有一个密码框，别的啥都没有，非常简洁。
- 输入密码进去后，只有一个搜索框。
- 搜索框只能搜索别名，只能精确搜索，并且区分大小写。

## 别名

- 别名 (Alias) 是本软件的特色功能
- 建议设定尽量简短的别名，比如：
  - 用 v2 表示 v2ex.com
  - 用 w 或 vx 表示微信
  - 用 jd 表示京东

### 使用别名有三大好处

#### 1. 直接得出想要的精确条目，就像命令行一样精确，很爽

#### 2. 可以代替顶置功能，并且比顶置功能好用很多

- 使用普通密码管理软件，当你记录的条目越来越多，你会希望顶置一些重要或常用的条目
- 但顶置的条目会越来越多，逐渐失去顶置的意义
- 于是只好从一堆顶置条目中忍痛取消一些，这个过程很烦
- 另外，你可能不想让别人知道你上某些网站，因此使用密码管理软件时，你需要躲避旁人的目光（即使密码已隐藏，但你连网站名字都不想让人看到）
- 只要使用别名，上述问题全部解决

#### 3. 轻松支持多密码

- 比如一个网站你有登录密码、支付密码、两部验证密码等
- 一般密码管理软件可能不支持这种情况
- 使用别名，比如一个网站有两个密码，那就新建两个条目，然后设定相同的别名即可
- 比如我自己的真实使用案例，输入 jd 即可得到一个列表，共两个条目，一是登录密码，一是支付密码。

## 历史记录功能

每次修改，不管修改了用户名、密码、标题还是备注，一律保留历史记录。
因此，可以尽管放心修改，绝对不用担心覆盖旧信息。

## Requirement

安装和初始设置本软件，你需要具备：

- 一台 Linux 服务器
- 在服务器上假设网站的基础知识
- Go 语言基础知识（入门级知识即可）
- 最好知道 Nginx 的基本使用方法（非必须）

## 本地版

- 本程序另外有一个本地版，不需要架设服务器，安装过程简单很多，并且支持云备份功能。
  <https://github.com/ahui2016/mima-go>
- 注：与本地版 (mima-go) 相比，这个服务器版本 (mima-web) 的代码精简了很多。而且，
  mima-go 没有做前后端分离，前端用的是 Go 标准库里的 html template，而 mima-web
  则改成了前后端分离。

## 注意事项

由于本软件采用网站的形式，通过浏览器访问，因此，为了提高安全性，避免浏览器缓存，
请使用浏览器的隐身模式 (incognito mode)。

## FAQ

### 如何修改主密码

$ cd mima-web
$ go run cmd/change-password -oldPass 当前密码 -newPass 新密码

### 如何备份数据

本程序提供了下载数据备份的网页，但未在页面中提供连接，请直接在浏览器地址栏手动输入

www.example.com/download

下载回来的是加密数据，用于备份。

### 如何导出全部数据为已解密的 json 文件

即将提供该功能

## 免责声明

我已尽自己的能力最大程度地确保本软件的安全性 (因为我自己就是本软件的深度用户),
用户可免费使用本软件, 可自行审查本软件的源代码，但万一有什么泄密、删除数据、
造成直接或间接损失等, 我一概不负责任。
（使用别的任何密码管理软件, 即使是收费的, 他们也一样不会负责用户的损失.）
