package notification

import (
	"fmt"

	"petshop/pkg/email"
)

type Service struct {
	client email.Client
	appURL string
}

func NewService(client email.Client, appURL string) *Service {
	return &Service{client: client, appURL: appURL}
}

func (s *Service) SendOrderStatusUpdate(to, customerName, orderID, newStatus string) {
	subject := fmt.Sprintf("Pesanan #%s - %s", orderID[:8], statusLabel(newStatus))
	body := fmt.Sprintf(`
<h2>Halo %s,</h2>
<p>Status pesanan Anda (#%s) telah diperbarui menjadi: <strong>%s</strong></p>
<p>Lihat detail pesanan: <a href="%s/pesanan/%s">%s/pesanan/%s</a></p>
<br><p>Salam,<br>PetShop Team</p>`,
		customerName, orderID[:8], statusLabel(newStatus),
		s.appURL, orderID, s.appURL, orderID)
	_ = s.client.Send(to, subject, body)
}

func (s *Service) SendPasswordReset(to, name, resetLink string) {
	subject := "Reset Password - PetShop"
	body := fmt.Sprintf(`
<h2>Halo %s,</h2>
<p>Kami menerima permintaan reset password untuk akun Anda.</p>
<p>Klik link berikut untuk mengatur password baru:</p>
<p><a href="%s">%s</a></p>
<p>Link berlaku 15 menit.</p>
<br><p>Salam,<br>PetShop Team</p>`, name, resetLink, resetLink)
	_ = s.client.Send(to, subject, body)
}

func statusLabel(s string) string {
	labels := map[string]string{
		"pending":    "Menunggu Pembayaran",
		"confirmed":  "Dikonfirmasi",
		"processing": "Diproses",
		"shipped":    "Dikirim",
		"delivered":  "Terkirim",
		"cancelled":  "Dibatalkan",
	}
	if l, ok := labels[s]; ok {
		return l
	}
	return s
}
