import time
import api

server_url = "http://localhost:8083"
tv_show_id=215720 # Королева слез
season_number=1
torrent_href = "https://rutracker.org/forum/viewtopic.php?t=6523513"

### Тест созданный для создания видеоконтента и доставки файлов
# В результате должен быть загружен сезон сериала и добавлен на медиасервер 

def waiting_delivery_step(status):
    count = 0
    while count < 100000:
            result = api.get_delivery_data(
                tv_show_id=tv_show_id,
                season_number=season_number,
                url=server_url
            )
            print(result['step'])
            if result['step'] == status:
                return result
            time.sleep(1)
            count += 1
    raise Exception("waiting timeout")
    
# Пример использования функции
if __name__ == "__main__":
    try:
        ### 0. Копируем (hardlink) файлы раздачи из backup каталога
        api.create_hardlinks(
            src_dir_name="Королева слёз  Queen of Tears (озвучка DubLikTV)",
            dest_base_path="/nfs/media/downloads",
            src_base_path="/nfs/media/backup"
        )
        
        ### 1. Создаем контент для TV шоу
        result = api.create_video_content(
            tv_show_id=tv_show_id,
            season_number=season_number,
            url=server_url
        )
        
        if result['delivery_status'] != "DeliveryStatusNew":
            raise Exception('Fail create_video_content')
        print("Видео контент успешно создан")

        ### 2. Создаем доставку контента
        result = api.create_content_delivery(
            tv_show_id=tv_show_id,
            season_number=season_number,
            url=server_url
        )

        if result['status'] != "NewStatus":
            raise Exception('Fail create_content_delivery')
        print("Доставка контента успешно создана")

        ### 3. Ожидание статуса выбора торрент файла "WaitingUserChoseTorrent"
        print("Ожидание выбора торрента")
        waiting_delivery_step("WaitingUserChoseTorrent")

        ### TODO: проверка списка найденных раздача

        ### 4. Клиент выбирает торрент
        result = api.choose_torrent(
            tv_show_id=tv_show_id,
            season_number=season_number,
            href=torrent_href,
            url=server_url
        )
        print("Торрент успешно выбран")

        ### 5. Ожидание статуса выбора WaitingChoseFileMatches
        print("Ожидание выбора выбора метча файлов")
        result = waiting_delivery_step("WaitingChoseFileMatches")

        ### TODO: проверка списка файлов торрент раздачи

        ### 6.Подверждаем метчинг файлов
        content_matches = result['data']['content_matches']
        result = api.approve_file_matches(
            tv_show_id=tv_show_id,
            season_number=season_number,
            content_matches=content_matches,
            url=server_url
        )
        print("Подтверждение метча файлов")
        
        print(result["step"])

    except Exception as e:
        print(f"Произошла ошибка: {e}")