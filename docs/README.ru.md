# Gitea Mail Templates — Шаблоны писем Gitea

Коллекция шаблонов писем для [Gitea](https://about.gitea.com).

> **110 файлов — 10 стилей, по 11 типов писем**

---

## Галерея стилей

| Превью | Стиль | Аудитория | Характер |
|---|---|---|---|
| ![Horizon](images/horizon.png) | **Horizon** | Корпоративный | Синий акцент, slate, centred cards |
| ![Terminal](images/terminal.png) | **Terminal** | Разработчики | Тёмная тема, monospace, зелёный CLI |
| ![Ember](images/ember.png) | **Ember** | Open Source | Тёплый янтарь, rounded, инклюзивный |
| ![Bloom](images/bloom.png) | **Bloom** | Креатив / Стартапы | Матовое стекло, мягкий синий свет |
| ![Heritage](images/heritage.png) | **Heritage** | Образование | Navy и золото, serif, классика |
| ![Neon](images/neon.png) | **Neon** | Игры / Web3 / Креатив | Киберпанк-неон, розовый и циан |
| ![Mono](images/mono.png) | **Mono** | Дизайн / Редакция | Швейцарский брутализм, чёрный-белый-красный |
| ![Terra](images/terra.png) | **Terra** | Устойчивость / Здоровье | Тёплые земляные тона, органические текстуры |
| ![Ink](images/ink.png) | **Ink** | Издательство / Новости | Редакционная печать, navy и золото |
| ![Aurora](images/aurora.png) | **Aurora** | Премиум SaaS | Эфирные градиенты, пурпур и бирюза |

> Изображения шириной 600px из [предпросмотра](../preview/index.html). Инструкции по созданию скриншотов: [images/README.md](images/README.md).

[**Галерея предпросмотра**](../preview/index.html)

---

## Установка

```bash
cp -r themes/horizon/mail/* /var/lib/gitea/custom/templates/mail/
systemctl restart gitea
```

### Проверка работы

Тестовое письмо администратора не использует пользовательские шаблоны. Чтобы
убедиться, что шаблоны активны, вызовите реальное уведомление. Самый быстрый
способ — сброс пароля: выйдите из системы, нажмите **"Forgot password"** на
странице входа и проверьте письмо для сброса.

## Предпросмотр

**Статический режим:**
```bash
cd tools && go run . preview all
```
Затем откройте `preview/index.html`.

**Dev-сервер (live reload):**
```bash
cd tools && go run . dev 
# → http://localhost:3456
```

## Совместимость

- **Gitea 1.25+** — структура директорий из v1.25
- **Последняя проверка:** Gitea 1.26.4
- 100% совместимость с официальными шаблонами Gitea — см. [COMPATIBILITY.md](COMPATIBILITY.md)

## Лицензия

MIT — [LICENSE](../LICENSE).
