package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/spf13/viper"
)

type Message struct {
	Host     string
	Port     int
	Username string
	Password string
	header   map[string]string
	body     string
}

func NewMessage() *Message {
	m := &Message{
		header: make(map[string]string),
	}

	return m
}

func (m *Message) SetHeader(field string, value string) {
	m.header[field] = value
}

func (m *Message) SetBody(value string) {
	m.body = value
}

// 生成邮件内容
func (m *Message) GenerateMessage() (message string) {
	for k, v := range m.header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + m.body
	return
}

func (m *Message) SetDialer(host string, port int, username, password string) {
	m.Host = host
	m.Port = port
	m.Username = username
	m.Password = password
}

func (m *Message) DialAndSend(email string) error {
	msg := []byte(m.GenerateMessage())
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}

	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("AUTH"); ok {
		if err = c.Auth(auth); err != nil {
			return err
		}
	}

	if err = c.Mail(m.Username); err != nil {
		return err
	}

	if err = c.Rcpt(email); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write(msg); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

/**
 * 发送邮箱验证码
 * param: email 目标邮箱
 * param: code 邮箱验证码
 * return: 发送失败时的错误信息
 */
func SendCaptcha(email string, code string) error {
	// 定义收件人
	mailTo := email
	// 邮件主题
	subject := "验证码"
	// 邮件正文
	body := "<h3>尊敬的用户：</h3><p>您好! 您的验证码是 <span style='color:red'> " + code + "</span>，五分钟内有效，祝您生活愉快！</p>"
	return Send(mailTo, subject, body)
}

/**
 * 发送电子邮件
 * param: emailList 目标邮箱数组
 * param: subject 邮件主题
 * param: body 邮件内容
 * return: 发送失败时的错误信息
 */
func Send(email string, subject string, body string) error {
	user, pass, host, port, addresser :=
		viper.GetString("mail.user"),
		viper.GetString("mail.pass"),
		viper.GetString("mail.host"),
		viper.GetInt("mail.port"),
		viper.GetString("mail.addresser")

	m := NewMessage()
	m.SetHeader("From", addresser+"<"+user+">") //添加别名
	m.SetHeader("To", email)                    //发送给多个用户
	m.SetHeader("Subject", subject)             //设置邮件主题
	m.SetBody(body)                             //设置邮件正文

	m.SetDialer(host, port, user, pass)

	err := m.DialAndSend(email)
	return err
}
