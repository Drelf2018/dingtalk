# dingtalk

钉钉群聊机器人，基于原项目 [`blinkbean/dingtalk`](https://github.com/blinkbean/dingtalk) 重构。

目前支持发送的消息类型有：

- [文本](#text-文本类型)
- [链接](#link-链接类型)
- [Markdown](#markdown-类型)
- [整体跳转 ActionCard](#整体跳转-actioncard-类型)
- [独立跳转 ActionCard](#独立跳转-actioncard-类型)
- [FeedCard](#feedcard-类型)
- [模板](#template-模板类型)

## 配置

### 创建机器人

1. 选择添加`自定义`机器人。
2. 安全设置
    共有关键词、加签、IP白名单三种设置，需要根据情况进行选择。
    ![Xnip2020-07-05_15-55-24.jpg](https://i.loli.net/2020/07/05/4XqHG2dOwo8StEu.jpg)

### 获取库

```go
go get github.com/Drelf2018/dingtalk
```

### 初始化

参考机器人结构体中字段的注释，根据需要填入需要的值。

```go
// Bot 钉钉机器人
type Bot struct {
	// 名称，可自定义
	Name string `json:"name" yaml:"name" toml:"name" long:"name"`

	// 调用接口的凭证，钉钉提供的 Webhook 链接中 access_token 的值
	Token string `json:"token" yaml:"token" toml:"token" long:"token"`

	// 安全密钥，创建机器人时在安全设置项选择了加签后，钉钉提供的 SEC 开头的字符串
	Secret string `json:"secret" yaml:"secret" toml:"secret" long:"secret"`

	// 自定义关键词，创建机器人时在安全设置项填入的所有关键词。当消息文本中不包含任何一个关键词时，会自动在文本末尾添加第一个关键词
	Keywords []string `json:"keywords" yaml:"keywords" toml:"keywords" long:"keywords"`

	// 全局请求超时时间，值为正时生效
	Timeout time.Duration `json:"timeout" yaml:"timeout" toml:"timeout" long:"timeout"`
}
```

## 使用

本库的方法参数与官方文档 [自定义机器人发送消息的消息类型](https://open.dingtalk.com/document/dingstart/custom-bot-send-message-type) 保持一致，你可以在文档中查看消息类型区别、参数含义、消息样式。

### SendHandler 请求处理器

```go
// 发送消息接口的前处理器，可以用来生成的加密签名、设置消息幂等、设置@等
type SendHandler func(*Send) error

// 内置了五个常用的处理器，可自行在代码中查看使用方法
var _ = []SendHandler{Secret(""), UUID(""), AtAll, AtMobile(""), AtUserID("")}
```

### Text 文本类型

```go
// SendText 发送文本类型消息
func (b *Bot) SendText(content string, handlers ...SendHandler) error
```

### Link 链接类型

```go
// SendLink 发送链接类型消息
func (b *Bot) SendLink(title, text, msgURL, picURL string, handlers ...SendHandler) error
```

### Markdown 类型

```go
// SendMarkdown 发送 markdown 类型消息
func (b *Bot) SendMarkdown(title, text string, handlers ...SendHandler) error
```

### 整体跳转 ActionCard 类型

```go
// SendActionCard 发送整体跳转 actionCard 类型消息
func (b *Bot) SendActionCard(title, text, singleTitle, singleURL string, handlers ...SendHandler) error
```

### 独立跳转 ActionCard 类型

```go
// ActionCardBtn actionCard 类型消息的按钮
type ActionCardBtn struct {
	// 按钮上显示的文本
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// 按钮跳转的 URL
	ActionURL string `json:"actionURL" yaml:"actionURL" toml:"actionURL" long:"actionURL"`
}

// SendActionsCard 发送独立跳转 actionCard 类型消息
func (b *Bot) SendActionsCard(title, text string, btns []ActionCardBtn, handlers ...SendHandler) error
```

### FeedCard 类型

```go
// FeedCardLink feedCard 类型消息的内容
type FeedCardLink struct {
	// 每条内容的标题
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// 每条内容上午跳转链接
	MessageURL string `json:"messageURL" yaml:"messageURL" toml:"messageURL" long:"messageURL"`

	// 每条内容的图片 URL ，建议使用上传媒体文件接口获取
	PicURL string `json:"picURL" yaml:"picURL" toml:"picURL" long:"picURL"`
}

// SendFeedCard 发送 feedCard 类型消息
func (b *Bot) SendFeedCard(links []FeedCardLink, handlers ...SendHandler) error
```
