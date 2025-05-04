package controller

import (
	"context"
	"homework/pkg/log"
	"net/http"
	"time"
)

type handler struct {
	service paymentService
	log     log.Logger
}

type paymentService interface {
	RegisterInvoice(ctx context.Context, timeRegistration time.Time) (err error)
	RegisterPayment(ctx context.Context, payment float64) (err error)
}

func (h *handler) WithRouter(mux *http.ServeMux) {
	// Invoice
	mux.Handle("POST /api/invoice", http.HandlerFunc(h.RegisterInvoice))
	// Payment
	mux.Handle("POST /api/payments", http.HandlerFunc(h.RegisterPayment))
}

// 1. Если идёт выбор билета, то из тикет сервиса получаем список доступных билетов (статус тиража: запланирован, билет без user_id в статусе: готов)
// 2. Если пользователь добавляет билет сам, то тикет сервис просто добавляет билет в лотерею (статус тиража: запланирован)
// 	в обоих случаях возвращается его номер
// 3. Создаём инвойс с номером билета (сообщаем тикет сервису о резервировании)
// 4. Ждём оплаты
// 5. Если оплата инвойса успешна, то тикет сервису сообщается об этом и билет переводится в статус КУПЛЕН
// Если не успешно или истекло время оплаты, то вызывается сценарий отмены в тикет сервисе, и билет:
// 1. Стирается user_id и билет снова готов к покупке
// 2. Билет удаляется из базы

// /ticket/allow_list/{draw_id} - получение списка доступных билетов
// CreateReservedTicket - резервируешь билет, возвращается номер
// /ticket/add/{draw_id} - добавление билета, возвращается номер

// генерируешь инвойс, ждёшь оплату

// Успешно:
// /ticket/success/{ticket_id} - билет куплен, возвращается ОК

// Неуспешно:
// /ticket/failure/{ticket_id} - билет не куплен (ошибка оплаты), билет или удаляется из БД, или возвращается в общий пул доступных билетов (стирается user_id)
