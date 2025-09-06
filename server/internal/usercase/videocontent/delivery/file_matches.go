package delivery

func (s *Service) NeedPrepareFileMatches(contentMatches []ContentMatch) bool {
	needToMerge := false
	for _, m := range contentMatches {
		// Если есть субтитры или аудиодорожки то нужно мержить
		if len(m.AudioFiles) > 0 || len(m.Subtitles) > 0 {
			needToMerge = true
			break
		}
	}
	return needToMerge
}
