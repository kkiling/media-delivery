/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/** @default "StatusUnknown" */
export enum MediadeliveryStatus {
  StatusUnknown = "StatusUnknown",
  NewStatus = "NewStatus",
  InProgressStatus = "InProgressStatus",
  CompletedStatus = "CompletedStatus",
  FailedStatus = "FailedStatus",
}

/** @default "TORRENT_STATE_UNKNOWN" */
export enum TorrentState {
  TORRENT_STATE_UNKNOWN = "TORRENT_STATE_UNKNOWN",
  TORRENT_STATE_ERROR = "TORRENT_STATE_ERROR",
  TORRENT_STATE_UPLOADING = "TORRENT_STATE_UPLOADING",
  TORRENT_STATE_DOWNLOADING = "TORRENT_STATE_DOWNLOADING",
  TORRENT_STATE_STOPPED = "TORRENT_STATE_STOPPED",
  TORRENT_STATE_QUEUED = "TORRENT_STATE_QUEUED",
}

/**
 * - TVShowDeliveryStatusUnknown: Неизвестный статус доставки
 *  - GenerateSearchQuery: Генерация запроса к трекеру
 *  - SearchTorrents: Поиск раздач сезона сериала/фильма
 *  - WaitingUserChoseTorrent: Ожидание выбора раздачи пользователем
 *  - GetMagnetLink: Получение магнет ссылки
 *  - AddTorrentToTorrentClient: Добавление раздачи для скачивания торрент клиентом
 *  - PrepareFileMatches: Получение информации о файлах раздачи
 *  - WaitingChoseFileMatches: Ожидание подтверждения пользователем соответствий выбора файлов
 *  - WaitingTorrentDownloadComplete: Ожидание завершения окончания скачивания раздачи
 *  - CreateVideoContentCatalogs: Формирование каталогов и иерархии файлов
 *  - DeterminingNeedConvertFiles: Определение необходимости конвертации файлов
 *  - StartMergeVideoFiles: Запуск конвертирования файлов
 *  - WaitingMergeVideoFiles: Ожидание завершения конвертации файлов
 *  - CreateHardLinkCopy: Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
 *  - GetCatalogsSize: GetCatalogsSize получение размеров каталогов сериала
 *  - SetMediaMetaData: Установка методаных серий сезона сериала/фильма в медиасервере
 *  - SendDeliveryNotification: Отправка уведомления в telegramm о успешной доставки
 *  - WaitingTorrentFiles: Ожидание когда появится информация о файлах в раздаче
 *  - GetEpisodesData: получение информации о эпизодах и каталоге сезона
 * @default "TVShowDeliveryStatusUnknown"
 */
export enum TVShowDeliveryStatus {
  TVShowDeliveryStatusUnknown = "TVShowDeliveryStatusUnknown",
  GenerateSearchQuery = "GenerateSearchQuery",
  SearchTorrents = "SearchTorrents",
  WaitingUserChoseTorrent = "WaitingUserChoseTorrent",
  GetMagnetLink = "GetMagnetLink",
  AddTorrentToTorrentClient = "AddTorrentToTorrentClient",
  PrepareFileMatches = "PrepareFileMatches",
  WaitingChoseFileMatches = "WaitingChoseFileMatches",
  WaitingTorrentDownloadComplete = "WaitingTorrentDownloadComplete",
  CreateVideoContentCatalogs = "CreateVideoContentCatalogs",
  DeterminingNeedConvertFiles = "DeterminingNeedConvertFiles",
  StartMergeVideoFiles = "StartMergeVideoFiles",
  WaitingMergeVideoFiles = "WaitingMergeVideoFiles",
  CreateHardLinkCopy = "CreateHardLinkCopy",
  GetCatalogsSize = "GetCatalogsSize",
  SetMediaMetaData = "SetMediaMetaData",
  SendDeliveryNotification = "SendDeliveryNotification",
  WaitingTorrentFiles = "WaitingTorrentFiles",
  GetEpisodesData = "GetEpisodesData",
}

/**
 * - TorrentSiteForbidden: Торрент трекер не доступен
 *  - FilesAlreadyExist: Файлы на медиасервере уже существуют
 * @default "TVShowDeliveryError_Unknown"
 */
export enum ErrorType {
  TVShowDeliveryErrorUnknown = "TVShowDeliveryError_Unknown",
  TorrentSiteForbidden = "TorrentSiteForbidden",
  FilesAlreadyExist = "FilesAlreadyExist",
}

/** @default "DeliveryStatusUnknown" */
export enum DeliveryStatus {
  DeliveryStatusUnknown = "DeliveryStatusUnknown",
  DeliveryStatusFailed = "DeliveryStatusFailed",
  DeliveryStatusInProgress = "DeliveryStatusInProgress",
  DeliveryStatusDelivered = "DeliveryStatusDelivered",
}

export interface Any {
  "@type"?: string;
  [key: string]: any;
}

export interface ChoseFileMatchesOptionsRequest {
  content_id?: ContentID;
  /** Пользователь подтверждает сметченные файлы */
  approve?: boolean;
}

export interface ChoseFileMatchesOptionsResponse {
  result?: TVShowDeliveryState;
}

export interface ChoseTorrentOptionsRequest {
  content_id?: ContentID;
  /** Пользователь выбрал конкретный торрента файл */
  href?: string;
  /** Пользователь поменял поисковый запрос */
  new_search_query?: string;
}

export interface ChoseTorrentOptionsResponse {
  result?: TVShowDeliveryState;
}

export interface ContentID {
  /** @format uint64 */
  movie_id?: number;
  tv_show?: TVShowID;
}

export interface ContentMatches {
  episode?: EpisodeInfo;
  video?: VideoFile;
  audio_files?: Track[];
  subtitles?: Track[];
}

export interface CreateVideoContentRequest {
  content_id?: ContentID;
}

export interface CreateVideoContentResponse {
  result?: VideoContent;
}

export interface Episode {
  /** @format uint64 */
  id?: number;
  /** @format date-time */
  air_date?: string;
  /** @format int64 */
  episode_number?: number;
  episode_type?: string;
  name?: string;
  overview?: string;
  /** @format int64 */
  runtime?: number;
  still?: Image;
  /** @format float */
  vote_average?: number;
  /** @format int64 */
  vote_count?: number;
}

export interface EpisodeInfo {
  /** @format int64 */
  season_number?: number;
  episode_name?: string;
  /** @format int64 */
  episode_number?: number;
  file_name?: string;
  relative_path?: string;
}

export interface FileInfo {
  relative_path?: string;
  full_path?: string;
  /** @format int64 */
  size?: string;
  extension?: string;
}

export interface GetSeasonInfoResponse {
  season?: Season;
  episodes?: Episode[];
}

export interface GetTVShowDeliveryDataResponse {
  result?: TVShowDeliveryState;
}

export interface GetTVShowInfoResponse {
  result?: TVShow;
}

export interface GetTVShowsFromLibraryResponse {
  items?: TVShowShort[];
}

export interface GetVideoContentResponse {
  items?: VideoContent[];
}

export interface Image {
  id?: string;
  w92?: string;
  w154?: string;
  w185?: string;
  w342?: string;
  w500?: string;
  w780?: string;
  original?: string;
}

export interface MergeVideoStatus {
  /** @format float */
  progress?: number;
  is_complete?: boolean;
}

export interface SearchQuery {
  Query?: string;
}

export interface SearchTVShowResponse {
  items?: TVShowShort[];
}

export interface Season {
  /** @format uint64 */
  id?: number;
  /** @format date-time */
  air_date?: string;
  /** @format int64 */
  episode_count?: number;
  name?: string;
  overview?: string;
  poster?: Image;
  /** @format int64 */
  season_number?: number;
  /** @format float */
  vote_average?: number;
}

export interface TVShow {
  /** @format uint64 */
  id?: number;
  name?: string;
  original_name?: string;
  overview?: string;
  poster?: Image;
  /** @format date-time */
  first_air_date?: string;
  /** @format float */
  vote_average?: number;
  /** @format int64 */
  vote_count?: number;
  /** @format float */
  popularity?: number;
  backdrop?: Image;
  genres?: string[];
  /** @format date-time */
  last_air_date?: string;
  /** @format int64 */
  number_of_episodes?: number;
  /** @format int64 */
  number_of_seasons?: number;
  origin_country?: string[];
  status?: string;
  tagline?: string;
  type?: string;
  seasons?: Season[];
}

export interface TVShowCatalog {
  /** Путь до раздачи сезона сериала */
  torrent_path?: string;
  /** Размер файлов раздачи сезона сериала */
  torrent_size_pretty?: string;
  /** Путь до сезона сериала на медиасервере */
  media_server_path?: TVShowCatalogPath;
  /** Размер файлов раздачи сезона сериала (байты) */
  media_server_size_pretty?: string;
  /**
   * Файлы скопированы с раздачи или созданы ссылочная связь
   * True - файлы скопированы
   * False - файлы созданы через линки
   */
  is_copy_files_in_media_server?: boolean;
}

export interface TVShowCatalogPath {
  /** Путь до каталога сериала */
  tv_show_path?: string;
  /** Путь до каталога сезона (относительно каталога сериала) */
  season_path?: string;
}

export interface TVShowDeliveryData {
  /** Поисковый запрос поиска торрент файла */
  search_query?: SearchQuery;
  /** Результат поиска торрент раздач */
  torrent_search?: TorrentSearch[];
  /** Результат метча файлов */
  content_matches?: ContentMatches[];
  /** статус скачивания раздачи */
  torrent_download_status?: TorrentDownloadStatus;
  /** статус сшивания файлов */
  merge_video_status?: MergeVideoStatus;
  /** TVShowCatalogInfo информация о каталогах сериала */
  tv_show_catalog_info?: TVShowCatalog;
}

export interface TVShowDeliveryError {
  raw_error?: string;
  error_type?: ErrorType;
}

export interface TVShowDeliveryState {
  data?: TVShowDeliveryData;
  step?: TVShowDeliveryStatus;
  status?: MediadeliveryStatus;
  error?: TVShowDeliveryError;
}

export interface TVShowID {
  /** @format uint64 */
  id?: number;
  /** @format int64 */
  season_number?: number;
}

export interface TVShowShort {
  /** @format uint64 */
  id?: number;
  name?: string;
  original_name?: string;
  overview?: string;
  poster?: Image;
  /** @format date-time */
  first_air_date?: string;
  /** @format float */
  vote_average?: number;
  /** @format int64 */
  vote_count?: number;
  /** @format float */
  popularity?: number;
}

export interface TorrentDownloadStatus {
  state?: TorrentState;
  /** @format float */
  progress?: number;
  is_complete?: boolean;
}

export interface TorrentSearch {
  title?: string;
  href?: string;
  category?: string;
  size?: string;
  /** @format int64 */
  seeds?: string;
  /** @format int64 */
  leeches?: string;
  /** @format int64 */
  downloads?: string;
  added_date?: string;
}

export interface Track {
  file?: FileInfo;
  name?: string;
  language?: string;
}

export interface VideoContent {
  id?: string;
  /** @format date-time */
  created_at?: string;
  delivery_status?: DeliveryStatus;
}

export interface VideoFile {
  file?: FileInfo;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: Any[];
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Media delivery API
 * @version 0.1
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  v1 = {
    /**
     * No description
     *
     * @tags VideoContentService
     * @name VideoContentServiceGetVideoContent
     * @summary Получение доставок для кино/тв сериала
     * @request GET:/v1/content
     */
    videoContentServiceGetVideoContent: (
      query?: {
        /** @format uint64 */
        "content_id.movie_id"?: number;
        /** @format uint64 */
        "content_id.tv_show.id"?: number;
        /** @format int64 */
        "content_id.tv_show.season_number"?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<GetVideoContentResponse, RpcStatus>({
        path: `/v1/content`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VideoContentService
     * @name VideoContentServiceCreateVideoContent
     * @summary Создание файловой раздачи
     * @request POST:/v1/content
     */
    videoContentServiceCreateVideoContent: (
      body: CreateVideoContentRequest,
      params: RequestParams = {},
    ) =>
      this.request<CreateVideoContentResponse, RpcStatus>({
        path: `/v1/content`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VideoContentService
     * @name VideoContentServiceChoseFileMatchesOptions
     * @summary Подтверждение метча файлов
     * @request PATCH:/v1/tvshow/delivery/chose-file-matches
     */
    videoContentServiceChoseFileMatchesOptions: (
      body: ChoseFileMatchesOptionsRequest,
      params: RequestParams = {},
    ) =>
      this.request<ChoseFileMatchesOptionsResponse, RpcStatus>({
        path: `/v1/tvshow/delivery/chose-file-matches`,
        method: "PATCH",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VideoContentService
     * @name VideoContentServiceChoseTorrentOptions
     * @summary Выбор раздачи с торрента
     * @request PATCH:/v1/tvshow/delivery/chose-torrent
     */
    videoContentServiceChoseTorrentOptions: (
      body: ChoseTorrentOptionsRequest,
      params: RequestParams = {},
    ) =>
      this.request<ChoseTorrentOptionsResponse, RpcStatus>({
        path: `/v1/tvshow/delivery/chose-torrent`,
        method: "PATCH",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VideoContentService
     * @name VideoContentServiceGetTvShowDeliveryData
     * @summary Получение данных стейта доставки
     * @request GET:/v1/tvshow/delivery/data
     */
    videoContentServiceGetTvShowDeliveryData: (
      query?: {
        /** @format uint64 */
        "content_id.movie_id"?: number;
        /** @format uint64 */
        "content_id.tv_show.id"?: number;
        /** @format int64 */
        "content_id.tv_show.season_number"?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<GetTVShowDeliveryDataResponse, RpcStatus>({
        path: `/v1/tvshow/delivery/data`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TVShowLibraryService
     * @name TvShowLibraryServiceGetTvShowInfo
     * @summary Получение подробной информации о сериале
     * @request GET:/v1/tvshow/info/{tv_show_id}
     */
    tvShowLibraryServiceGetTvShowInfo: (
      tvShowId: string,
      params: RequestParams = {},
    ) =>
      this.request<GetTVShowInfoResponse, RpcStatus>({
        path: `/v1/tvshow/info/${tvShowId}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TVShowLibraryService
     * @name TvShowLibraryServiceGetSeasonInfo
     * @summary Получение информации о сезоне сериала и его сериях
     * @request GET:/v1/tvshow/info/{tv_show_id}/{season_number}
     */
    tvShowLibraryServiceGetSeasonInfo: (
      tvShowId: string,
      seasonNumber: number,
      params: RequestParams = {},
    ) =>
      this.request<GetSeasonInfoResponse, RpcStatus>({
        path: `/v1/tvshow/info/${tvShowId}/${seasonNumber}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TVShowLibraryService
     * @name TvShowLibraryServiceGetTvShowsFromLibrary
     * @summary Получение списка сериалов из библиотеки
     * @request GET:/v1/tvshow/library
     */
    tvShowLibraryServiceGetTvShowsFromLibrary: (params: RequestParams = {}) =>
      this.request<GetTVShowsFromLibraryResponse, RpcStatus>({
        path: `/v1/tvshow/library`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TVShowLibraryService
     * @name TvShowLibraryServiceSearchTvShow
     * @summary Поиск сериалов по названию
     * @request GET:/v1/tvshow/search
     */
    tvShowLibraryServiceSearchTvShow: (
      query?: {
        query?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<SearchTVShowResponse, RpcStatus>({
        path: `/v1/tvshow/search`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),
  };
}
