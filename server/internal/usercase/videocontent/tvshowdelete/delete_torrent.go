package tvshowdelete

import (
	"context"
	"fmt"
)

func (s *Service) DeleteTorrentFromTorrentClient(ctx context.Context, magnetHash string) error {
	return fmt.Errorf("not yet implemented")
}

func (s *Service) DeleteTorrentFiles(ctx context.Context, torrentPath string) error {
	panic("implement me")
}
