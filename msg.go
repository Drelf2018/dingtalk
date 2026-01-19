package dingtalk

// MsgType 表示消息类型的字符串，已内置五种类型
//
//	MsgText       // 文本类型
//	MsgLink       // 链接类型，不支持@人
//	MsgMarkdown   // Markdown 类型
//	MsgActionCard // 整体跳转、独立跳转类型
//	MsgFeedCard   // FeedCard 类型，不支持@人
type MsgType string

const (
	MsgText       MsgType = "text"       // 文本类型
	MsgLink       MsgType = "link"       // 链接类型，不支持@人
	MsgMarkdown   MsgType = "markdown"   // Markdown 类型
	MsgActionCard MsgType = "actionCard" // 整体跳转、独立跳转类型
	MsgFeedCard   MsgType = "feedCard"   // FeedCard 类型，不支持@人
)

// Msg 消息接口
type Msg interface {
	Type() MsgType
}

// Text 文本类型消息
type Text struct {
	Content string `json:"content" yaml:"content" toml:"content" long:"content"` // 文本消息的内容
}

func (Text) Type() MsgType {
	return MsgText
}

var _ Msg = Text{}

// Link 链接类型消息
type Link struct {
	// 链接消息标题
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// 链接消息的内容
	Text string `json:"text" yaml:"text" toml:"text" long:"text"`

	// 点击消息跳转的 URL
	MessageURL string `json:"messageUrl" yaml:"messageUrl" toml:"messageUrl" long:"messageUrl"`

	// 链接消息内的图片地址，建议使用上传媒体文件接口获取
	PicURL string `json:"picUrl,omitempty" yaml:"picUrl" toml:"picUrl" long:"picUrl"`
}

func (Link) Type() MsgType {
	return MsgLink
}

var _ Msg = Link{}

// Markdown markdown 类型消息
type Markdown struct {
	// 消息会话列表中展示的标题，非消息体的标题
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// markdown 类型消息的文本内容
	Text string `json:"text" yaml:"text" toml:"text" long:"text"`
}

func (Markdown) Type() MsgType {
	return MsgMarkdown
}

var _ Msg = Markdown{}

// ActionCardBtn actionCard 类型消息的按钮
type ActionCardBtn struct {
	// 按钮上显示的文本
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// 按钮跳转的 URL
	ActionURL string `json:"actionURL" yaml:"actionURL" toml:"actionURL" long:"actionURL"`
}

// ActionCard actionCard 类型消息
type ActionCard struct {
	// 消息会话列表中展示的标题，非消息体的标题
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// actionCard 类型消息的正文内容，支持 markdown 语法
	Text string `json:"text" yaml:"text" toml:"text" long:"text"`

	// 单个按钮的方案，设置此项和 singleURL 后 btns 无效
	SingleTitle string `json:"singleTitle,omitempty" yaml:"singleTitle" toml:"singleTitle" long:"singleTitle"`

	// 点击 singleTitle 按钮触发的 URL
	SingleURL string `json:"singleURL,omitempty" yaml:"singleURL" toml:"singleURL" long:"singleURL"`

	// 按钮的信息列表
	Btns []ActionCardBtn `json:"btns,omitempty" yaml:"btns" toml:"btns" long:"btns"`

	// 消息内按钮排列方式，0：按钮竖直排列，1：按钮横向排列
	BtnOrientation string `json:"btnOrientation,omitempty" yaml:"btnOrientation" toml:"btnOrientation" long:"btnOrientation"`
}

func (ActionCard) Type() MsgType {
	return MsgActionCard
}

var _ Msg = ActionCard{}

// FeedCardLink feedCard 类型消息的内容
type FeedCardLink struct {
	// 每条内容的标题
	Title string `json:"title" yaml:"title" toml:"title" long:"title"`

	// 每条内容上午跳转链接
	MessageURL string `json:"messageURL" yaml:"messageURL" toml:"messageURL" long:"messageURL"`

	// 每条内容的图片 URL ，建议使用上传媒体文件接口获取
	PicURL string `json:"picURL" yaml:"picURL" toml:"picURL" long:"picURL"`
}

// FeedCard feedCard 类型消息
type FeedCard struct {
	// feedCard 类型消息的内容列表
	Links []FeedCardLink `json:"links" yaml:"links" toml:"links" long:"links"`
}

func (FeedCard) Type() MsgType {
	return MsgFeedCard
}

var _ Msg = FeedCard{}
