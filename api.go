package dingtalk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/Drelf2018/req"
	"github.com/Drelf2018/req/method"
)

// GenerateSign 生成加密时间戳和签名，加签的方式是将时间戳和密钥当做签名字符串，
// 开发者服务内当前系统时间戳，单位是毫秒，与请求调用时间误差不能超过 1 小时，
// 使用 HmacSHA256 算法计算签名，然后进行 Base64 编码，得到最终的签名
func GenerateSign(secret string) (int64, string, error) {
	hmacSHA256 := hmac.New(sha256.New, []byte(secret))
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	_, err := fmt.Fprintf(hmacSHA256, "%d\n%s", timestamp, secret)
	if err != nil {
		return 0, "", fmt.Errorf("dingtalk: failed to generate signature: %w", err)
	}
	return timestamp, base64.StdEncoding.EncodeToString(hmacSHA256.Sum(nil)), nil
}

// At 被@的群成员信息
type At struct {
	IsAtAll   bool     `json:"isAtAll,omitempty"`   // 是否@所有人
	AtMobiles []string `json:"atMobiles,omitempty"` // 被@的群成员手机号
	AtUserIDs []string `json:"atUserIds,omitempty"` // 被@的群成员 userId
}

// Send 自定义机器人发送群消息
type Send struct {
	// 要发送的消息
	Msg Msg

	// 自定义机器人调用接口的凭证
	AccessToken string `req:"query"`

	// 使用时间戳和密钥生成的加密签名
	Sign string `req:"query,omitempty"`

	// 开发者服务内当前系统时间戳，单位是毫秒，与请求调用时间误差不能超过 1 小时
	Timestamp int64 `req:"query,omitempty"`

	// 消息类型
	MsgType MsgType `req:"body:msgtype"`

	// 消息幂等，发消息时接口调用超时或未知错误等报错，开发者可使用同一个消息幂等重试，避免重复发出消息
	MsgUUID string `req:"body:msgUuid,omitempty"`

	// 被@的群成员信息
	At At `req:"body,omitempty"`

	// 请求头
	ContentType string `req:"header" default:"application/json"`
}

func (*Send) Method() string {
	return http.MethodPost
}

func (*Send) RawURL() string {
	return "https://oapi.dingtalk.com/robot/send"
}

var _ req.API = (*Send)(nil)

func (s *Send) Body(r *http.Request, value reflect.Value, body []reflect.StructField) (io.Reader, error) {
	if s.Msg != nil {
		s.MsgType = s.Msg.Type()
	}
	m := method.MakeJSONMap(r.Context(), value, body)
	m[string(s.MsgType)] = s.Msg
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return nil, err
	}
	buf.Truncate(buf.Len() - 1) // Seeing the source code of (*json.Encoder).Encode
	return buf, nil
}

var _ req.APIBody = (*Send)(nil)

// 发送消息接口的前处理器，可以用来设置@、指定消息UUID和修改请求参数
type SendHandler func(*Send) error

// Secret 会自动设置生成的加密签名，密钥参数为机器人安全设置页面，加签一栏下面显示的 SEC 开头的字符串
func Secret(secret string) SendHandler {
	return func(s *Send) (err error) {
		s.Timestamp, s.Sign, err = GenerateSign(secret)
		return
	}
}

// UUID 设置消息幂等
func UUID(uuid string) SendHandler {
	return func(s *Send) error {
		s.MsgUUID = uuid
		return nil
	}
}

// AtAll @所有人
func AtAll(s *Send) error {
	s.At.IsAtAll = true
	return nil
}

// AtMobile @指定群成员手机号
func AtMobile(mobiles ...string) SendHandler {
	return func(s *Send) error {
		s.At.AtMobiles = mobiles
		return nil
	}
}

// AtUserID @指定群成员 userId
func AtUserID(ids ...string) SendHandler {
	return func(s *Send) error {
		s.At.AtUserIDs = ids
		return nil
	}
}

var _ = []SendHandler{Secret(""), UUID(""), AtAll, AtMobile(""), AtUserID("")}

// SendResponse 发送消息响应体
type SendResponse struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
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

// PostSendWithContext 携带上下文发送消息
func PostSendWithContext(ctx context.Context, token string, msg Msg, handlers ...SendHandler) (r SendResponse, err error) {
	api := &Send{Msg: msg, AccessToken: token}
	for _, handler := range handlers {
		if err = handler(api); err != nil {
			return
		}
	}
	r, err = req.ResultWithContext[SendResponse](ctx, api)
	if err == nil && r.ErrCode != 0 {
		err = SendError{API: api, ErrMsg: r.ErrMsg, ErrCode: r.ErrCode}
	}
	return
}

// PostSendWithContext 发送消息
func PostSend(token string, msg Msg, handlers ...SendHandler) (SendResponse, error) {
	return PostSendWithContext(context.Background(), token, msg, handlers...)
}
