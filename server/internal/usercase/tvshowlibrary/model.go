package tvshowlibrary

import "time"

type Image struct {
	// ID изображения
	ID string
	// Урезанная версия изображения ширина 92px
	W92 *string
	// Урезанная версия изображения ширина 185px
	W185 *string
	// Урезанная версия изображения ширина 342px
	W342 *string
	// Оригинальный путь до изображения
	Original string
}

// TVShowShort базовая информация о сериале (отображение при поиске сериала например)
type TVShowShort struct {
	// TheMovieDb TVShow ID
	ID uint64
	// Наименование сериала (зависит от языка выбранного при запросе к api)
	Name string
	// Имя сериала на оригинальном языке
	OriginalName string
	// Описание сериала (зависит от языка выбранного при запросе к api)
	Overview string
	// Изображение постера сериала
	Poster *Image
	// Дата выхода первой серии (например 2019-07-08)
	FirstAirDate time.Time
	// Средняя оценка от 0.0 до 10.0
	VoteAverage float32
	// Количество оценок
	VoteCount uint32
	// Популярность сериала (вроде как популярность от 0 до 500)
	Popularity float32
}

// TVShow расширенная информация о сериале (отображении в карточке сериала)
type TVShow struct {
	TVShowShort
	// Фоновое изображение
	Backdrop *Image
	// Список жанров (зависит от языка выбранного при запросе к api)
	// Например мультфильм, драма, боевик и тд
	Genres []string
	// Дата выхода последней серии на данный момент (например 2023-06-20)
	LastAirDate time.Time
	// Количество серий на данный момент (считается только основные сезоны, без доп материалов и тд)
	NumberOfEpisodes uint32
	// Количество сезонов на данный момент (считается только основные сезоны, без доп материалов и тд)
	NumberOfSeasons uint32
	// Список стран участвующих в производстве (jp, ru и тд)
	OriginCountry []string
	// Статус
	// Returning Series - Выходит
	// Ended - Завершен
	// Могут быть еще
	Status *string
	// Слоган сериала (зависит от языка выбранного при запросе к api)
	Tagline string
	// Тип - ???
	// Scripted
	Type *string
	// Список сезонов сериала
	Seasons []Season
}

// Season базовая информация о сезоне сериала
type Season struct {
	// Дата выхода первой серии
	AirDate time.Time
	// Количество эпизодов
	EpisodeCount uint32
	// Наименование сезона (зависит от языка выбранного при запросе к api)
	Name string
	// Описание сезона (зависит от языка выбранного при запросе к api)
	Overview string
	// Изображение постера сериала
	Poster *Image
	// Порядковый номер сезона (0 - спец материалы)
	SeasonNumber uint8
	// Средняя оценка от 0.0 до 10.0
	VoteAverage float32
}

type SeasonWithEpisodes struct {
	Season
	Episodes []Episode
}

// Episode информация о эпизоде сезона
type Episode struct {
	// Дата выхода (пример 2013-04-07)
	AirDate time.Time
	// Номер эпизода в сезоне
	EpisodeNumber int
	// Тип эпизода - ???
	// standard
	EpisodeType *string
	// Наименование эпизода (зависит от языка выбранного при запросе к api)
	Name string
	// Описание эпизода (зависит от языка выбранного при запросе к api)
	Overview string
	// Продолжительность эпизода (минуты)
	Runtime uint32
	// Изображения превью эпизода
	Still *Image
	// Средняя оценка от 0.0 до 10.0
	VoteAverage float32
	// Количество оценок
	VoteCount uint32
}
