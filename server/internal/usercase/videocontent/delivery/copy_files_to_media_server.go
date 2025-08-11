package delivery

import "context"

type CreateHardLinkCopyParams struct {
	ContentMatches []ContentMatches
}

// CreateHardLinkCopyToMediaServer шаг копирования файлов на медиа сервер
func (s *Service) CreateHardLinkCopyToMediaServer(ctx context.Context, params CreateHardLinkCopyParams) error {
	// TODO: формирование каталогов сериала на медиасервер
	// TODO: создание симлинков видеофайлов с торрент раздачи в каталогах медиасервера
	// TODO: Переход на следующий шаг - установки методаных
	panic("implement me")
}
