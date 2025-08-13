package dependencyinjections

import (
	"downloader/internal/domain"
	memoria "downloader/internal/infra/db/mem_db"
)

func NewVideoDatabase() domain.Database[domain.Video] {
	return memoria.NewMemoriaDatabase[domain.Video]()
}
