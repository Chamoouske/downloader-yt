package domain

type ProgressBar interface {
    Start(total int64)
    Update(current int64)
    Finish()
}
