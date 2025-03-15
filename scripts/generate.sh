#!/bin/bash

set -e

# Проверяем наличие api.yml
if [ ! -f "api.yml" ]; then
  echo "Ошибка: файл api.yml не найден в текущей директории"
  exit 1
fi

# Создаем директорию для сгенерированного кода
mkdir -p generated

# Устанавливаем Ogen, если его еще нет
go install github.com/ogen-go/ogen/cmd/ogen@latest

# Запускаем генерацию кода из файла api.yml
echo "Генерация кода из api.yml..."
ogen --target generated --package api --clean api.yml

echo "Код успешно сгенерирован в директории 'generated'"
