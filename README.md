# dingtalk

钉钉群聊机器人，基于原项目 [`blinkbean/dingtalk`](https://github.com/blinkbean/dingtalk) 重构

目前支持发送的消息类型有：

- [文本](#text-文本类型)
- [链接](#link-链接类型)
- [Markdown](#markdown-类型)
- [整体跳转 ActionCard](#整体跳转-actioncard-类型)
- [独立跳转 ActionCard](#独立跳转-actioncard-类型)
- [FeedCard](#feedcard-类型)

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

```go
// Bot 钉钉机器人
type Bot struct {
	// 机器人名称
	Name string `json:"name" yaml:"name" toml:"name" long:"name"`

	// 调用接口的凭证
	Token string `json:"token" yaml:"token" toml:"token" long:"token"`

	// 机器人安全密钥
	Secret string `json:"secret" yaml:"secret" toml:"secret" long:"secret"`

	// 请求超时时间，非零值时生效
	Timeout time.Duration `json:"timeout" yaml:"timeout" toml:"timeout" long:"timeout"`

	// 自定义关键词，不为空且文本中不包含任意一个关键词时会自动在文本末尾添加第一个关键词
	Keywords []string `json:"keywords" yaml:"keywords" toml:"keywords" long:"keywords"`
}
```

上面是我们的机器人模型，请根据需要填入字段值：

- `Name` 是机器人名称，可自定义，用于请求失败时打印错误提示。
- `Token` 是机器人凭证，也就是钉钉提供的 `Webhook` 链接中 `access_token` 的值。
- `Secret` 是机器人密钥，创建机器人时安全设置项选择了加签后，钉钉提供的 `SEC` 开头的字符串。
- `Timeout` 是机器人超时时间，默认不提供，即不设置超时
- `Keywords` 是机器人自定义关键词，创建机器人时安全设置项选择了此项后的所有关键词，不要求顺序。当消息文本中不包含任何一个关键词时，会自动在文本末尾添加第一个关键词。

## 使用

本库的方法参数与官方文档 [自定义机器人发送消息的消息类型](https://open.dingtalk.com/document/dingstart/custom-bot-send-message-type) 保持一致，你可以在文档中查看消息类型区别、参数含义、消息样式。

### SendHandler 请求处理器

```go
// 发送消息接口的前处理器，可以用来设置@、指定消息UUID和修改请求参数
type SendHandler func(*Send) error

// 内置了四个常用的处理器，可自行在代码中查看使用方法
var _ = []SendHandler{AtAll, AtMobile(""), AtUserID(""), UUID("")}
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
// SendSingleActionCard 发送整体跳转 actionCard 类型消息
func (b *Bot) SendSingleActionCard(title, text, singleTitle, singleURL string, handlers ...SendHandler) error
```

### 独立跳转 ActionCard 类型

```go
// ActionCardBtn actionCard 类型消息的按钮
type ActionCardBtn struct {
	Title     string `json:"title,omitempty"`     // 按钮上显示的文本
	ActionURL string `json:"actionURL,omitempty"` // 按钮跳转的 URL
}

// SendActionCard 发送独立跳转 actionCard 类型消息
func (b *Bot) SendActionCard(title, text string, btns []ActionCardBtn, handlers ...SendHandler) error
```

### FeedCard 类型

```go
// FeedCardLink feedCard 类型消息的内容
type FeedCardLink struct {
	Title      string `json:"title,omitempty"`      // feedCard 消息内每条内容的标题
	MessageURL string `json:"messageURL,omitempty"` // feedCard 消息内每条内容上午跳转链接
	PicURL     string `json:"picURL,omitempty"`     // feedCard 消息内每条内容的图片 URL ，建议使用上传媒体文件接口获取
}

// SendFeedCard 发送 feedCard 类型消息
func (b *Bot) SendFeedCard(links []FeedCardLink, handlers ...SendHandler) error
```