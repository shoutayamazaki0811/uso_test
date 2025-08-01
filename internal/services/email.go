package services

import (
    "bytes"
    "fmt"
    "html/template"
    "github.com/resendlabs/resend-go"
)

type EmailService struct {
    client *resend.Client
    from   string
}

func NewEmailService(apiKey, fromEmail string) *EmailService {
    client := resend.NewClient(apiKey)
    return &EmailService{
        client: client,
        from:   fromEmail,
    }
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(to, name string) error {
    subject := "usoへようこそ！"
    html := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
            <div style="background-color: #0a0a0a; padding: 20px; text-align: center;">
                <h1 style="color: #d4af37; margin: 0;">uso</h1>
            </div>
            <div style="background-color: #1a1a1a; color: #ffffff; padding: 30px;">
                <h2>ようこそ、%sさん！</h2>
                <p>usoへのご登録ありがとうございます。</p>
                <p>高級エンターテイメントマッチングプラットフォームusoで、素敵な出会いをお楽しみください。</p>
                <div style="text-align: center; margin: 30px 0;">
                    <a href="https://uso.app/login" style="background-color: #d4af37; color: #0a0a0a; padding: 15px 30px; text-decoration: none; border-radius: 8px; font-weight: bold;">ログインする</a>
                </div>
                <p style="color: #a0a0a0; font-size: 14px;">このメールに心当たりがない場合は、無視してください。</p>
            </div>
        </div>
    `, name)

    params := &resend.SendEmailRequest{
        From:    s.from,
        To:      []string{to},
        Subject: subject,
        Html:    html,
    }

    _, err := s.client.Emails.Send(params)
    return err
}

// SendBookingConfirmation sends booking confirmation email
func (s *EmailService) SendBookingConfirmation(to, guestName, castName string, bookingDetails BookingDetails) error {
    subject := "予約確認 - uso"
    html := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
            <div style="background-color: #0a0a0a; padding: 20px; text-align: center;">
                <h1 style="color: #d4af37; margin: 0;">uso</h1>
            </div>
            <div style="background-color: #1a1a1a; color: #ffffff; padding: 30px;">
                <h2>予約が確定しました</h2>
                <p>%sさん、以下の内容で予約が確定しました。</p>
                
                <div style="background-color: #2a2a2a; padding: 20px; border-radius: 8px; margin: 20px 0;">
                    <h3 style="color: #d4af37;">予約詳細</h3>
                    <p><strong>キャスト:</strong> %sさん</p>
                    <p><strong>日時:</strong> %s</p>
                    <p><strong>場所:</strong> %s</p>
                    <p><strong>料金:</strong> ¥%s</p>
                </div>
                
                <p>当日は時間に余裕を持ってお越しください。</p>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="https://uso.app/bookings/%s" style="background-color: #d4af37; color: #0a0a0a; padding: 15px 30px; text-decoration: none; border-radius: 8px; font-weight: bold;">予約詳細を見る</a>
                </div>
            </div>
        </div>
    `, guestName, castName, bookingDetails.DateTime, bookingDetails.Location, bookingDetails.Amount, bookingDetails.ID)

    params := &resend.SendEmailRequest{
        From:    s.from,
        To:      []string{to},
        Subject: subject,
        Html:    html,
    }

    _, err := s.client.Emails.Send(params)
    return err
}

// SendMatchNotification sends notification when users match
func (s *EmailService) SendMatchNotification(to, userName, matchName string) error {
    subject := "マッチしました！ - uso"
    html := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
            <div style="background-color: #0a0a0a; padding: 20px; text-align: center;">
                <h1 style="color: #d4af37; margin: 0;">uso</h1>
            </div>
            <div style="background-color: #1a1a1a; color: #ffffff; padding: 30px;">
                <h2 style="text-align: center;">🎉 マッチしました！</h2>
                <p>%sさん、おめでとうございます！</p>
                <p>%sさんとマッチしました。早速メッセージを送ってみましょう。</p>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="https://uso.app/messages" style="background-color: #d4af37; color: #0a0a0a; padding: 15px 30px; text-decoration: none; border-radius: 8px; font-weight: bold;">メッセージを送る</a>
                </div>
                
                <p style="color: #a0a0a0; font-size: 14px; text-align: center;">素敵な出会いになりますように！</p>
            </div>
        </div>
    `, userName, matchName)

    params := &resend.SendEmailRequest{
        From:    s.from,
        To:      []string{to},
        Subject: subject,
        Html:    html,
    }

    _, err := s.client.Emails.Send(params)
    return err
}

// SendPasswordReset sends password reset email
func (s *EmailService) SendPasswordReset(to, name, resetToken string) error {
    subject := "パスワードリセット - uso"
    resetURL := fmt.Sprintf("https://uso.app/reset-password?token=%s", resetToken)
    
    html := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
            <div style="background-color: #0a0a0a; padding: 20px; text-align: center;">
                <h1 style="color: #d4af37; margin: 0;">uso</h1>
            </div>
            <div style="background-color: #1a1a1a; color: #ffffff; padding: 30px;">
                <h2>パスワードリセット</h2>
                <p>%sさん、</p>
                <p>パスワードリセットのリクエストを受け付けました。以下のボタンをクリックして新しいパスワードを設定してください。</p>
                
                <div style="text-align: center; margin: 30px 0;">
                    <a href="%s" style="background-color: #d4af37; color: #0a0a0a; padding: 15px 30px; text-decoration: none; border-radius: 8px; font-weight: bold;">パスワードをリセット</a>
                </div>
                
                <p style="color: #a0a0a0; font-size: 14px;">このリンクは24時間有効です。心当たりがない場合は、このメールを無視してください。</p>
            </div>
        </div>
    `, name, resetURL)

    params := &resend.SendEmailRequest{
        From:    s.from,
        To:      []string{to},
        Subject: subject,
        Html:    html,
    }

    _, err := s.client.Emails.Send(params)
    return err
}

// BookingDetails contains booking information for emails
type BookingDetails struct {
    ID       string
    DateTime string
    Location string
    Amount   string
}

// SendEmailFromTemplate sends email using a template
func (s *EmailService) SendEmailFromTemplate(to, subject, templateName string, data interface{}) error {
    tmpl, err := template.ParseFiles(fmt.Sprintf("templates/email/%s.html", templateName))
    if err != nil {
        return err
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return err
    }

    params := &resend.SendEmailRequest{
        From:    s.from,
        To:      []string{to},
        Subject: subject,
        Html:    buf.String(),
    }

    _, err = s.client.Emails.Send(params)
    return err
}