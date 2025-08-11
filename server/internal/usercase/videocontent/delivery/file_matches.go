package delivery

func (s *Service) NeedPrepareFileMatches(ContentMatches []ContentMatches) bool {
	needToMerge := false
	for _, m := range ContentMatches {
		// Если есть субтитры или аудиодорожки то нужно мержить
		if len(m.AudioFiles) > 0 || len(m.Subtitles) > 0 {
			needToMerge = true
			break
		}
	}
	return needToMerge
}
