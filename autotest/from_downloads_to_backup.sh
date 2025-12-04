# Определяем название каталога
# dir_name="Королева слёз  Queen of Tears (озвучка DubLikTV)"
dir_name="Peacemaker.S02.1080p.AMZN.WEB-DL.H.264-EniaHD"


# Создаем целевую директорию
mkdir -p "/nfs/media/backup/${dir_name}"

# Создаем хардлинки для всех файлов
ln "/nfs/media/downloads/${dir_name}"/* "/nfs/media/backup/${dir_name}/"