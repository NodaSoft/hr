package task

import (
	"testing"
	"time"
)

func TestTask_Execute(t *testing.T) {
	type fields struct {
		ID         int64
		CreatedAt  time.Time
		FinishedAt time.Time
		Result     string
		Successful bool
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "creation_time_error",
			fields: fields{
				CreatedAt: time.Date(0, 0, 0, 0, 0, 0, 1, time.UTC),
				Result: taskCreationTimeErrorMsg,
			},
		},
		{
			name: "execution_timeout",
			fields: fields{
				CreatedAt: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
				Result: taskExecutionTimeoutMsg,
			},
		},
		{
			name: "normal_execution",
			fields: fields{
				Result: taskSucceededMsg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				ID:         tt.fields.ID,
				CreatedAt:  tt.fields.CreatedAt,
				FinishedAt: tt.fields.FinishedAt,
				Result:     tt.fields.Result,
				Successful: tt.fields.Successful,
			}

			// для кейсов без необходимости манипулирования временем создания задачи
			if task.CreatedAt.IsZero() {
				task.CreatedAt = time.Now()
			}

			task.Execute()

			if task.Result != tt.fields.Result {
				t.Errorf("got unexpected result: %s, expected: %s\n", task.Result, tt.fields.Result)
			}
		})
	}
}
