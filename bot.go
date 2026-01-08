package dingtalk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Drelf2018/req"
)

// Bot 钉钉机器人
type Bot struct {
	Name     string        // 机器人名称，用于打印错误信息
	Token    string        // 调用接口的凭证
	Secret   string        // 机器人安全密钥
	Timeout  time.Duration // 请求超时时间，非零值时生效
	Keywords []string      // 自定义关键词，不为空且文本中不包含任意一个关键词时会自动在文本末尾添加第一个关键词
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
	api := &Send{Secret: b.Secret, AccessToken: b.Token, Msg: msg}
	for _, handler := range handlers {
		if err := handler(api); err != nil {
			return err
		}
	}
	if b.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.Timeout)
		defer cancel()
	}
	_, err := req.ResultWithContext[SendResponse](ctx, api)
	if err != nil {
		return fmt.Errorf("dingtalk: Bot %s failed to send: %w", b.Name, err)
	}
	return nil
}

// Send 发送消息
func (b *Bot) Send(msg Msg, handlers ...SendHandler) error {
	return b.SendWithContext(context.Background(), msg, handlers...)
}

// SendTextWithContext 携带上下文发送文本类型消息
func (b *Bot) SendTextWithContext(ctx context.Context, content string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(content) {
		content = fmt.Sprintf("%s【%s】", content, b.Keywords[0])
	}
	return b.SendWithContext(ctx, Text{Content: content}, handlers...)
}

// SendText 发送文本类型消息
func (b *Bot) SendText(content string, handlers ...SendHandler) error {
	return b.SendTextWithContext(context.Background(), content, handlers...)
}

// SendLinkWithContext 携带上下文发送链接类型消息
func (b *Bot) SendLinkWithContext(ctx context.Context, title, text, picURL, msgURL string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text = fmt.Sprintf("%s【%s】", text, b.Keywords[0])
	}
	return b.SendWithContext(ctx, Link{MessageURL: msgURL, Title: title, PicURL: picURL, Text: text}, handlers...)
}

// SendLink 发送链接类型消息
func (b *Bot) SendLink(title, text, picURL, msgURL string, handlers ...SendHandler) error {
	return b.SendLinkWithContext(context.Background(), title, text, picURL, msgURL, handlers...)
}

// SendMarkdownWithContext 携带上下文发送 markdown 类型消息
func (b *Bot) SendMarkdownWithContext(ctx context.Context, title, text string, handlers ...SendHandler) error {
	if !b.ContainsAnyKeyword(title) && !b.ContainsAnyKeyword(text) {
		text = fmt.Sprintf("%s【%s】", text, b.Keywords[0])
	}
	return b.SendWithContext(ctx, Markdown{Text: text, Title: title}, handlers...)
}

// SendMarkdown 发送 markdown 类型消息
func (b *Bot) SendMarkdown(title, text string, handlers ...SendHandler) error {
	return b.SendMarkdownWithContext(context.Background(), title, text, handlers...)
}
