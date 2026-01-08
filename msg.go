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
	Content string `json:"content,omitempty"` // 文本消息的内容
}

func (Text) Type() MsgType {
	return MsgText
}

var _ Msg = Text{}

// Link 链接类型消息
type Link struct {
	MessageURL string `json:"messageUrl,omitempty"` // 点击消息跳转的 URL
	Title      string `json:"title,omitempty"`      // 链接消息标题
	PicURL     string `json:"picUrl,omitempty"`     // 链接消息内的图片地址，建议使用上传媒体文件接口获取
	Text       string `json:"text,omitempty"`       // 链接消息的内容
}

func (Link) Type() MsgType {
	return MsgLink
}

var _ Msg = Link{}

// Markdown markdown 类型消息
type Markdown struct {
	Text  string `json:"text,omitempty"`  // markdown 类型消息的文本内容
	Title string `json:"title,omitempty"` // 消息会话列表中展示的标题，非消息体的标题
}

func (Markdown) Type() MsgType {
	return MsgMarkdown
}

var _ Msg = Markdown{}

// ActionCardBtn actionCard 类型消息的按钮
type ActionCardBtn struct {
	ActionURL string `json:"actionURL,omitempty"` // 按钮跳转的 URL
	Title     string `json:"title,omitempty"`     // 按钮上显示的文本
}

// ActionCard actionCard 类型消息
type ActionCard struct {
	HideAvatar     string          `json:"hideAvatar,omitempty"`     // 是否显示消息发送者头像，0：正常发消息者头像，1：隐藏发消息者头像
	BtnOrientation string          `json:"btnOrientation,omitempty"` // 消息内按钮排列方式，0：按钮竖直排列，1：按钮横向排列
	SingleURL      string          `json:"singleURL,omitempty"`      // 点击 singleTitle 按钮触发的 URL
	SingleTitle    string          `json:"singleTitle,omitempty"`    // 单个按钮的方案，设置此项和 singleURL 后 btns 无效
	Text           string          `json:"text,omitempty"`           // actionCard 类型消息的正文内容，支持 markdown 语法
	Title          string          `json:"title,omitempty"`          // 消息会话列表中展示的标题，非消息体的标题
	Btns           []ActionCardBtn `json:"btns,omitempty"`           // 按钮的信息列表
}

func (ActionCard) Type() MsgType {
	return MsgActionCard
}

var _ Msg = ActionCard{}

// FeedCardLink feedCard 类型消息的内容
type FeedCardLink struct {
	PicURL     string `json:"picURL,omitempty"`     // feedCard 消息内每条内容的图片 URL ，建议使用上传媒体文件接口获取
	MessageURL string `json:"messageURL,omitempty"` // feedCard 消息内每条内容上午跳转链接
	Title      string `json:"title,omitempty"`      // feedCard 消息内每条内容的标题
}

// FeedCard feedCard 类型消息
type FeedCard struct {
	Links []FeedCardLink `json:"links,omitempty"` // feedCard 类型消息的内容列表
}

func (FeedCard) Type() MsgType {
	return MsgFeedCard
}

var _ Msg = FeedCard{}
