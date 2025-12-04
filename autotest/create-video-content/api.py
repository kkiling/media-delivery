import os
import json
import requests
from typing import Dict, Any

def create_hardlinks(src_dir_name, dest_base_path, src_base_path):
    """
    Создает директорию и хардлинки для всех файлов из исходной директории
    
    Args:
        src_dir_name (str): Название исходной директории
        dest_base_path (str): Базовый путь для целевой директории
        src_base_path (str): Базовый путь для исходной директории
    """
    # Формируем полные пути
    src_dir = os.path.join(src_base_path, src_dir_name)
    dest_dir = os.path.join(dest_base_path, src_dir_name)
    
    # Создаем целевую директорию (если не существует)
    os.makedirs(dest_dir, exist_ok=True)
    
    # Создаем хардлинки для всех файлов в исходной директории
    try:
        for item in os.listdir(src_dir):
            src_path = os.path.join(src_dir, item)
            dest_path = os.path.join(dest_dir, item)
            
            # Создаем хардлинк только для файлов (не директорий)
            if os.path.isfile(src_path):
                os.link(src_path, dest_path)
                print(f"Создан хардлинк: {dest_path}")
    
    except FileNotFoundError:
        print(f"Ошибка: Исходная директория не найдена: {src_dir}")
    except PermissionError:
        print(f"Ошибка: Недостаточно прав для создания хардлинков")
    except Exception as e:
        print(f"Произошла ошибка: {e}")

def create_video_content(
    tv_show_id: int,
    season_number: int,
    url: str
) -> Dict[str, Any]:
    """
    Создает видео контент для указанного TV шоу и сезона.
    
    Args:
        tv_show_id (int): ID TV шоу
        season_number (int): Номер сезона
        url (str): URL API сервера
        
    Returns:
        Dict[str, Any]: Ответ от сервера с информацией о созданном контенте
        
    Raises:
        requests.exceptions.RequestException: При ошибках сетевого запроса
    """
    
    # Формируем полный URL
    url = f"{url}/v1/content"
    
    # Подготавливаем данные для запроса
    payload = {
        "content_id": {
            "tv_show": {
                "id": tv_show_id,
                "season_number": season_number
            }
        }
    }
    
    # Устанавливаем заголовки
    headers = {
        'accept': 'application/json',
        'Content-Type': 'application/json'
    }
    
    try:
        # Выполняем POST запрос
        response = requests.post(
            url,
            headers=headers,
            data=json.dumps(payload),
            timeout=30  # Таймаут 30 секунд
        )
        
        # Проверяем статус ответа
        response.raise_for_status()
        
        # Возвращаем JSON ответ
        return response.json()["result"]
        
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при создании видео контента: {e}")
        raise

def create_content_delivery(
    tv_show_id: int,
    season_number: int,
    url: str
) -> Dict[str, Any]:
    """
    Создает доставку контента для указанного TV шоу и сезона.
    
    Args:
        tv_show_id (int): ID TV шоу
        season_number (int): Номер сезона
        url (str): Базовый URL API сервера
        
    Returns:
        Dict[str, Any]: Ответ от сервера с информацией о доставке
        
    Raises:
        requests.exceptions.RequestException: При ошибках сетевого запроса
    """
    
    url = f"{url}/v1/content/state/delivery"
    
    payload = {
        "content_id": {
            "tv_show": {
                "id": tv_show_id,
                "season_number": season_number
            }
        }
    }
    
    headers = {
        'accept': 'application/json',
        'Content-Type': 'application/json'
    }
    
    try:
        response = requests.post(
            url,
            headers=headers,
            data=json.dumps(payload),
            timeout=30
        )
        response.raise_for_status()
        return response.json()["result"]
        
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при создании доставки контента: {e}")
        raise

def choose_torrent(
    tv_show_id: int,
    season_number: int,
    href: str,
    url: str
) -> Dict[str, Any]:
    """
    Выбирает торрент для доставки контента.
    
    Args:
        tv_show_id (int): ID TV шоу
        season_number (int): Номер сезона
        href (str): Ссылка на торрент
        url (str): Базовый URL API сервера
        
    Returns:
        Dict[str, Any]: Ответ от сервера
        
    Raises:
        requests.exceptions.RequestException: При ошибках сетевого запроса
    """
    
    url = f"{url}/v1/content/state/delivery/chose-torrent"
    
    payload = {
        "content_id": {
            "tv_show": {
                "id": tv_show_id,
                "season_number": season_number
            }
        },
        "href": href
    }
    
    headers = {
        'accept': 'application/json',
        'Content-Type': 'application/json'
    }
    
    try:
        response = requests.patch(
            url,
            headers=headers,
            data=json.dumps(payload),
            timeout=30
        )
        response.raise_for_status()
        return response.json()["result"]
        
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при выборе торрента: {e}")
        raise

def get_delivery_data(
    tv_show_id: int,
    season_number: int,
    url: str
) -> Dict[str, Any]:
    """
    Получает данные доставки контента для указанного TV шоу и сезона.
    
    Args:
        tv_show_id (int): ID TV шоу
        season_number (int): Номер сезона
        url (str): Базовый URL API сервера
        
    Returns:
        Dict[str, Any]: Ответ от сервера с данными доставки
        
    Raises:
        requests.exceptions.RequestException: При ошибках сетевого запроса
    """
    
    url = f"{url}/v1/content/state/delivery"
    
    # Параметры запроса
    params = {
        'content_id.tv_show.id': tv_show_id,
        'content_id.tv_show.season_number': season_number
    }
    
    headers = {
        'accept': 'application/json'
    }
    
    try:
        response = requests.get(
            url,
            headers=headers,
            params=params,
            timeout=30
        )
        response.raise_for_status()
        return response.json()["result"]
        
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при получении данных доставки: {e}")
        raise

def approve_file_matches(
    tv_show_id: int,
    season_number: int,
    content_matches: Dict[str, Any],
    url: str,
) -> Dict[str, Any]:
    """
    Подтверждает совпадения файлов для доставки контента.
    
    Args:
        tv_show_id (int): ID TV шоу
        season_number (int): Номер сезона
        content_matches (Dict[str, Any]): Данные совпадений контента
        url (str): Базовый URL API сервера
        
    Returns:
        Dict[str, Any]: Ответ от сервера
        
    Raises:
        requests.exceptions.RequestException: При ошибках сетевого запроса
    """
    
    url = f"{url}/v1/content/state/delivery/chose-file-matches"
    
    payload = {
        "content_id": {
            "tv_show": {
                "id": tv_show_id,
                "season_number": season_number
            }
        },
        "approve": True,
        "content_matches": content_matches
    }
    
    headers = {
        'accept': 'application/json',
        'Content-Type': 'application/json'
    }
    
    try:
        response = requests.patch(
            url,
            headers=headers,
            data=json.dumps(payload),
            timeout=30
        )
        response.raise_for_status()
        return response.json()["result"]
        
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при подтверждении совпадений файлов: {e}")
        raise