# Archivist

## 🇺🇸 English | [🇷🇺 Russian](#русский)

An automated tool for writing terms, dates, and chronologies

**Features:**

- Simple and flexible YAML configuration without code changes
- Support for multiple LLMs: OpenAI (and compatible ones like Gemini, Grok), Gigachat, Github Copilot (no authorization required)
- Flexible and customizable Typst templating language
- Full automation: no manual actions required after configuration

### Installation

1. [**Install Typst**](https://github.com/typst/typst#installation)

2. **Download Archivist**

   Download the appropriate binary for your platform from [here](https://github.com/kuzgoga/archivist/releases) and make it executable.
   Make the binary executable (Unix-like systems):

   ```bash
   chmod +x archivist-*
   ```

   We recommend to rename it to `archivist`.

3. **Create Assets Folder**

   Create an `assets` folder and place your data sources there. You can copy the assets folder from this repository as a starting point.

4. **Configuration File**

   Place an `archivist.yaml` configuration file next to the binary. Use one of the example configurations as a base:
   - [archivist.minimal-example.yaml](https://github.com/kuzgoga/archivist/blob/main/archivist.minimal-example.yaml) - for basic setup
   - [archivist.full-example.yaml](https://github.com/kuzgoga/archivist/blob/main/archivist.full-example.yaml) - for advanced configuration

5. **Run Archivist**

   Simply run the archivist binary and it will handle everything automatically:

   ```bash
   ./archivist
   ```

### Contribute

Before submitting a pull request, please run:

```bash
make lint
make build
```

**Prerequisites:**

- [Go compiler](https://go.dev/doc/install)
- [golangci-lint (install 2.4.0+ version)](https://golangci-lint.run/welcome/install/)

### Report Issue

When reporting an issue, please include:

- Archivist version
- Your configuration file (without API keys)
- Platform and architecture
- If you're modifying the code: Go version

---

## Русский

Инструмент для автоматического написания хронологий, терминов, персоналий.

**Возможности:**

- простая и гибкая YAML-конфигурация без изменения кода
- поддержка различных LLM: OpenAI (и совместимых Gemini, Grok), Gigachat, Github Copilot (не требует авторизации)
- гибкий и настраиваемый язык шаблонов Typst
- полная автоматизация: не требует ручных действий после настройки

### Установка

1. [**Установите Typst**](https://github.com/typst/typst#installation)

2. **Скачайте Archivist**

   Скачайте подходящий бинарный файл для вашей платформы [отсюда](https://github.com/kuzgoga/archivist/releases) и сделайте его исполняемым.

   Сделайте бинарный файл исполняемым (Unix-подобные системы):

   ```bash
   chmod +x archivist-*
   ```

   Рекомендуем переименовать его в `archivist`.

3. **Создайте папку Assets**

   Создайте папку `assets` и разместите там ваши источники данных. Вы можете скопировать папку assets из этого репозитория в качестве отправной точки.

4. **Файл конфигурации**

   Поместите файл конфигурации `archivist.yaml` рядом с бинарным файлом. Используйте один из примеров конфигурации в качестве основы:
   - [archivist.minimal-example.yaml](https://github.com/kuzgoga/archivist/blob/main/archivist.minimal-example.yaml) - для базовой настройки
   - [archivist.full-example.yaml](https://github.com/kuzgoga/archivist/blob/main/archivist.full-example.yaml) - для расширенной конфигурации

5. **Запустите Archivist**

   Просто запустите бинарный файл archivist, и он автоматически выполнит все необходимые действия:

   ```bash
   ./archivist
   ```

### Участие в разработке

Перед отправкой pull request, пожалуйста, выполните:

```bash
make lint
make build
```

**Необходимые компоненты:**

- [Компилятор Go](https://go.dev/doc/install)
- [golangci-lint (установите версию 2.4.0+)](https://golangci-lint.run/welcome/install/)

### Сообщить о проблеме

При сообщении о проблеме, пожалуйста, укажите:

- Версию Archivist
- Ваш файл конфигурации (без API ключей)
- Платформу и архитектуру
- Если вы модифицируете код: версию Go

---

### Contacts / Контакты

- [Telegram](https://t.me/kuzgoga)
- [Email](mailto:me@gogacoder.ru)
