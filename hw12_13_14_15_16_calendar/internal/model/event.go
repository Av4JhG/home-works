package domain

import "time"

type Event struct {
	ID          int64         `db:"id"`          // уникальный идентификатор события (можно воспользоваться UUID)
	Title       string        `db:"title"`       // Заголовок - короткий текст
	Date        time.Time     `db:"date"`        // Дата и время события
	Duration    time.Duration `db:"duration"`    // Длительность события (или дата и время окончания);
	Description string        `db:"description"` // Описание события - длинный текст, опционально;
	UserID      int64         `db:"user_id"`     // ID пользователя, владельца события;
	TimeShift   int64         `db:"timeshift"`   // За сколько времени высылать уведомление, опционально.
}
