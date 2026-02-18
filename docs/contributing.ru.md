# Участие в разработке

Спасибо за интерес к проекту! Этот документ описывает правила участия.

## Сообщения об ошибках

- Проверьте существующие [issues](https://github.com/Educentr/go-iproto/issues), чтобы избежать дублирования.
- Создайте новый issue с чётким заголовком и описанием.
- Укажите версию Go, ОС и минимальный воспроизводимый пример.

## Предложения по улучшению

- Создайте issue с описанием улучшения и его мотивацией.
- Для значительных изменений обсудите подход до отправки PR.

## Настройка окружения

1. **Go 1.24+** — [https://go.dev/dl/](https://go.dev/dl/)
2. **golangci-lint** — [https://golangci-lint.run/welcome/install/](https://golangci-lint.run/welcome/install/)
3. **make** — для запуска целей

```bash
git clone https://github.com/Educentr/go-iproto.git
cd go-iproto
make test
make lint
```

## Стиль кода

- Форматирование: `gofmt` и `goimports`.
- Запускайте `make lint` перед коммитом.
- Следуйте существующим соглашениям проекта.

## Процесс Pull Request

1. Сделайте форк и создайте ветку для изменений.
2. Коммиты с понятными сообщениями.
3. `make test` и `make lint` должны проходить.
4. Отправьте PR в ветку `main`.
5. Опишите, что делает PR, и укажите связанные issues.

## Лицензия

Отправляя изменения, вы соглашаетесь с их лицензированием по [MIT License](https://github.com/Educentr/go-iproto/blob/main/LICENSE).
