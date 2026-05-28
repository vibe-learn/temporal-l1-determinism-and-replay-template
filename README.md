        # temporal — Детерминизм и replay

        Homework-шаблон для урока **l1_determinism_and_replay** (Детерминизм и replay) на платформе Vibe Learn.

        ## Что делать

        На go.temporal.io/sdk возьми готовый workflow с намеренными нарушениями детерминизма
(time.Now, rand, итерация по map). Задача: (1) почини их через workflow.Now, workflow.SideEffect
и сортировку ключей; (2) добавь новый шаг через workflow.GetVersion, не ломая старые истории;
(3) напиши replay-тест на worker.NewWorkflowReplayer против приложенного JSON-дампа истории.
Тесты проверят: replayer проходит на починенном коде и падает на исходном; TestWorkflowEnvironment
подтверждает, что GetVersion-ветки выбираются корректно. В шаблоне лежат экспортированные истории.

## Контекст (из transfer-задачи урока)

В проде крутятся тысячи живых SubscriptionWorkflow. Тебе прилетела задача: добавить в начало
процесса новый шаг — антифрод-проверку (активность FraudCheck), и заодно ты заметил, что
кто-то написал в workflow `if time.Now().Hour() < 9 { ... }` и `for k := range configMap { ... }`.
Ты собираешься это поправить и выкатить.

**Вопрос:** разбери ситуацию. Опиши:
(a) почему time.Now() и итерация по map в workflow — баги, и как их корректно заменить;
(b) как добавить шаг FraudCheck, не сломав уже запущенные истории (что произойдёт без мер
    и какой инструмент SDK решает это);
(c) как ты убедишься ещё ДО деплоя, что новый код совместим с тем, что крутится в проде.

## Recap из урока

- Детерминизм обязателен, потому что **replay повторно исполняет workflow-функцию**: при той же истории код должен принять те же решения в том же порядке.
- Запрещено напрямую: **time.Now, rand/uuid, time.Sleep, нативные goroutine/каналы, итерация по map, чтение env/файлов/сети, глобальное мутабельное состояние**.
- Замены из SDK: **workflow.Now, workflow.Sleep, workflow.Go, workflow.SideEffect / MutableSideEffect**; map — через отсортированные ключи.
- Изменение логики у живых workflow — через **workflow.GetVersion** (старые истории → DefaultVersion, новые → новая ветка); иначе non-determinism error.
- Лови недетерминизм **ДО прода**: replay-тесты (NewWorkflowReplayer против записанных историй) в CI. Non-determinism error не теряет данные, но workflow застревает.

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go` (workflow + активности), прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - SDK: `go.temporal.io/sdk`
        - Docker + docker-compose — `docker compose up` поднимает Temporal dev server на `:7233` + Web UI на `:8233`. Адрес переопределяется через env `TEMPORAL_ADDRESS` (дефолт `localhost:7233`).
        - Юнит-тесты на `testsuite.TestWorkflowEnvironment` (активности замоканы) бегут в CI БЕЗ сервера; интеграционный тест включается через `TEMPORAL_INTEGRATION=1`.

        ## Запуск

        ```bash
        # Поднять локальный Temporal dev server + UI
        docker compose up -d
        # Web UI: http://localhost:8233

        # Прогнать тесты (юнит на TestWorkflowEnvironment — без сервера;
        # интеграционный включается через TEMPORAL_INTEGRATION=1)
        go test ./...
        TEMPORAL_INTEGRATION=1 go test ./...

        # Запустить воркер (регистрирует workflow + активности, слушает task queue)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
