package dingtalk

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/Drelf2018/req"
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

	// 自定义模板，使用“消息结构体名.字段名”的名称创建模板后，用于自动填充字符串类型字段值
	Template *template.Template
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

// SendError 发送消息错误
type SendError struct {
	API     *Send
	ErrMsg  string
	ErrCode int
}

func (s SendError) Error() string {
	return fmt.Sprintf("dingtalk: failed to send %T: %s (%d)", s.API.Msg, s.ErrMsg, s.ErrCode)
}

// SendWithContext 携带上下文发送消息
func (b *Bot) SendWithContext(ctx context.Context, msg Msg, handlers ...SendHandler) error {
	api := &Send{Secret: b.Secret, AccessToken: b.Token, Msg: msg}
	for _, handler := range handlers {
		if err := handler(api); err != nil {
			return err
		}
	}
	if b.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.Timeout)
		defer cancel()
	}
	r, err := req.ResultWithContext[SendResponse](ctx, api)
	if err != nil {
		return err
	}
	if r.ErrCode != 0 {
		return SendError{API: api, ErrMsg: r.ErrMsg, ErrCode: r.ErrCode}
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

// Funcs 设置函数
func (b *Bot) Funcs(funcMap template.FuncMap) *Bot {
	if b.Template == nil {
		b.Template = template.New("")
	}
	b.Template.Funcs(funcMap)
	return b
}

// NewTemplate 为机器人创建新模板
func (b *Bot) NewTemplate(name, text string) error {
	if b.Template == nil {
		b.Template = template.New("")
	}
	_, err := b.Template.New(name).Parse(text)
	return err
}

// ErrNilMsg 消息为空
var ErrNilMsg = errors.New("dingtalk: msg cannot be nil")

// Parse 为机器人创建模板，传入消息的字段导出、类型是字符串以及值不为空时，才会创建模板
func (b *Bot) Parse(msg Msg) error {
	if msg == nil {
		return ErrNilMsg
	}
	if b.Template == nil {
		b.Template = template.New("")
	}
	// 获取结构体对象
	elem := reflect.ValueOf(msg)
	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("dingtalk: invalid msg: expected (a pointer to) a struct, got: %T(%v)", msg, elem.Kind())
	}
	// 遍历字段，字段导出、类型是字符串以及值不为空时，才会创建模板
	structType := elem.Type()
	structName := structType.Name()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		if !fieldType.IsExported() || fieldType.Type.Kind() != reflect.String {
			continue
		}
		field := elem.Field(i)
		if field.IsZero() {
			continue
		}
		_, err := b.Template.New(structName + "." + fieldType.Name).Parse(field.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// ErrNilTemplate 机器人自定义模板为空
var ErrNilTemplate = errors.New("dingtalk: template cannot be nil")

// Fill 填充消息字段值，传入的消息为结构体，会新建一个结构体对象并返回，如果传入的消息为结构体指针，则不会额外返回值
func Fill(tmpl *template.Template, data any, msg Msg) (Msg, error) {
	if tmpl == nil {
		return nil, ErrNilTemplate
	}
	if msg == nil {
		return nil, ErrNilMsg
	}
	// 如果传入的是不可设置的结构体对象，则创建对应的指针并设置为当前对象
	var elem reflect.Value
	val := reflect.ValueOf(msg)
	if val.Kind() != reflect.Pointer {
		newValue := reflect.New(val.Type())
		elem = newValue.Elem()
		elem.Set(val)
	} else {
		elem = val.Elem()
	}
	// 再获取指针里的值，就可以设置字段值了
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("dingtalk: invalid msg type: expected (a pointer to) a struct, got: %T", msg)
	}
	// 遍历字段，字段导出、类型是字符串、值为空以及对应的模板存在时，才解析模板并设置值
	var b strings.Builder
	structType := elem.Type()
	structName := structType.Name()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		if !fieldType.IsExported() || fieldType.Type.Kind() != reflect.String {
			continue
		}
		field := elem.Field(i)
		if !field.IsZero() || !field.CanSet() {
			continue
		}
		templateName := structName + "." + fieldType.Name
		t := tmpl.Lookup(templateName)
		if t == nil {
			continue
		}
		if err := t.Execute(&b, data); err != nil {
			return nil, fmt.Errorf("dingtalk: failed to execute template %q: %w", templateName, err)
		}
		field.SetString(b.String())
		b.Reset()
	}
	if val.Kind() != reflect.Pointer {
		return elem.Interface().(Msg), nil
	}
	return nil, nil
}

// Fill 填充消息字段值，传入的消息为结构体，会新建一个结构体对象并返回，如果传入的消息为结构体指针，则不会额外返回值
func (b *Bot) Fill(data any, msg Msg) (Msg, error) {
	return Fill(b.Template, data, msg)
}

// SendTemplateMsgWithContext 携带上下文发送模板消息
func (b *Bot) SendTemplateMsgWithContext(ctx context.Context, data any, msg Msg, handlers ...SendHandler) (err error) {
	msg, err = b.Fill(data, msg)
	if err != nil {
		return
	}
	return b.SendWithContext(ctx, msg, handlers...)
}

// SendTemplateMsg 发送模板消息
func (b *Bot) SendTemplateMsg(data any, msg Msg, handlers ...SendHandler) error {
	return b.SendTemplateMsgWithContext(context.Background(), data, msg, handlers...)
}
