package dependencyinjections

import (
	"downloader/internal/domain"
	memoria "downloader/internal/infra/db/mem_db"
)

var db domain.Database[domain.Video]

func init() {
	db = memoria.NewMemoriaDatabase[domain.Video]()
}

func GetVideoDatabase() *domain.Database[domain.Video] {
	return &db
}
