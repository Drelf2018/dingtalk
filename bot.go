package dingtalk

import (
	"context"
	"strings"
	"time"
)

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

// ContainsAnyKeyword 检测字符串是否包含任意一个关键词，关键词切片为空也返回真
func (b *Bot) ContainsAnyKeyword(text string) bool {
	if len(b.Keywords) == 0 {
		return true
	}
	for _, keyword := range b.Keywords {
		if keyword == "" {
			continue
		}
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// SendWithContext 携带上下文发送消息
func (b *Bot) SendWithContext(ctx context.Context, msg Msg, handlers ...SendHandler) error {
	if b.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.Timeout)
		defer cancel()
	}
	if b.Secret != "" {
		handlers = append(handlers, Secret(b.Secret))
	}
	_, err := PostSendWithContext(ctx, b.Token, msg, handlers...)
	return err
}

// Send 发送消息
func (b *Bot) Send(msg Msg, handlers ...SendHandler) error {
	return b.SendWithContext(context.Background(), msg, handlers...)
}

// SendTextWithContext 携带上下文发送文本类型消息
func (b *Bot) SendTextWithContext(ctx context.Context, content string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(content) {
		content += b.Keywords[0]
	}
	return b.SendWithContext(ctx, Text{Content: content}, handlers...)
}

// SendText 发送文本类型消息
func (b *Bot) SendText(content string, handlers ...SendHandler) error {
	return b.SendTextWithContext(context.Background(), content, handlers...)
}

// SendLinkWithContext 携带上下文发送链接类型消息
func (b *Bot) SendLinkWithContext(ctx context.Context, title, text, msgURL, picURL string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text += b.Keywords[0]
	}
	return b.SendWithContext(ctx, Link{Title: title, Text: text, MessageURL: msgURL, PicURL: picURL}, handlers...)
}

// SendLink 发送链接类型消息
func (b *Bot) SendLink(title, text, msgURL, picURL string, handlers ...SendHandler) error {
	return b.SendLinkWithContext(context.Background(), title, text, msgURL, picURL, handlers...)
}

// SendMarkdownWithContext 携带上下文发送 markdown 类型消息
func (b *Bot) SendMarkdownWithContext(ctx context.Context, title, text string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text += b.Keywords[0]
	}
	return b.SendWithContext(ctx, Markdown{Title: title, Text: text}, handlers...)
}

// SendMarkdown 发送 markdown 类型消息
func (b *Bot) SendMarkdown(title, text string, handlers ...SendHandler) error {
	return b.SendMarkdownWithContext(context.Background(), title, text, handlers...)
}

// SendActionCardWithContext 携带上下文发送整体跳转 actionCard 类型消息
func (b *Bot) SendActionCardWithContext(ctx context.Context, title, text, singleTitle, singleURL string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text += b.Keywords[0]
	}
	return b.SendWithContext(ctx, ActionCard{Title: title, Text: text, SingleTitle: singleTitle, SingleURL: singleURL}, handlers...)
}

// SendActionCard 发送整体跳转 actionCard 类型消息
func (b *Bot) SendActionCard(title, text, singleTitle, singleURL string, handlers ...SendHandler) error {
	return b.SendActionCardWithContext(context.Background(), title, text, singleTitle, singleURL, handlers...)
}

// SendActionsCardWithContext 携带上下文发送独立跳转 actionCard 类型消息
func (b *Bot) SendActionsCardWithContext(ctx context.Context, title, text string, btns []ActionCardBtn, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text += b.Keywords[0]
	}
	return b.SendWithContext(ctx, ActionsCard{Title: title, Text: text, Btns: btns}, handlers...)
}

// SendActionsCard 发送独立跳转 actionCard 类型消息
func (b *Bot) SendActionsCard(title, text string, btns []ActionCardBtn, handlers ...SendHandler) error {
	return b.SendActionsCardWithContext(context.Background(), title, text, btns, handlers...)
}

// SendFeedCardWithContext 携带上下文发送 feedCard 类型消息
func (b *Bot) SendFeedCardWithContext(ctx context.Context, links []FeedCardLink, handlers ...SendHandler) error {
	if len(b.Keywords) != 0 {
		var hasKeyword bool
		for i := range links {
			if b.ContainsAnyKeyword(links[i].Title) {
				hasKeyword = true
				break
			}
		}
		if !hasKeyword {
			links[len(links)-1].Title += b.Keywords[0]
		}
	}
	return b.SendWithContext(ctx, FeedCard{Links: links}, handlers...)
}

// SendFeedCard 发送 feedCard 类型消息
func (b *Bot) SendFeedCard(links []FeedCardLink, handlers ...SendHandler) error {
	return b.SendFeedCardWithContext(context.Background(), links, handlers...)
}
