package entity

type SubmissionResponseDto struct {
	Id   int64           `json:"id"`
	Task TaskResponseDto `json:"task"`
}
type ListSubmissionsRequestDto struct {
	UserId    string
	TaskId    string
	StartDate string
	EndDate   string
}
type ScoreSubmissionRequestDto struct {
	Score    int32  `json:"score"`
	Feedback string `json:"feedback"`
}
