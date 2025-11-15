package tvshowdelete

import "context"

func (s *Service) DeleteTorrentFromTorrentClient(ctx context.Context, magnetHash string) error {
	panic("implement me")
}

func (s *Service) DeleteTorrentFiles(ctx context.Context, torrentPath string) error {
	panic("implement me")
}
