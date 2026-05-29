# Helm chart `course-service`

Этот документ — полный гайд по тому, как развернуть микросервис `course-service`
в Kubernetes с помощью Helm. Написан в расчёте на то, что **с k8s ты до этого не работал**,
поэтому тут много пояснений «почему именно так».

PostgreSQL и Kafka **внутри кластера не запускаются**. Чарт принимает их адреса как
параметры (`config.postgres.*`, `config.kafka.brokers`) и подключается к ним по сети.

---

## Содержание

1. [Что такое Helm и зачем он нужен](#1-что-такое-helm-и-зачем-он-нужен)
2. [Структура чарта](#2-структура-чарта)
3. [Что мы создаём в Kubernetes](#3-что-мы-создаём-в-kubernetes)
4. [Подготовка инструментов](#4-подготовка-инструментов)
5. [Шаг 1. Поднимаем кластер Kubernetes](#5-шаг-1-поднимаем-кластер-kubernetes)
6. [Шаг 2. Готовим внешние Postgres и Kafka](#6-шаг-2-готовим-внешние-postgres-и-kafka)
7. [Шаг 3. Собираем и публикуем Docker-образ](#7-шаг-3-собираем-и-публикуем-docker-образ)
8. [Шаг 4. Создаём namespace и (опционально) секрет реестра](#8-шаг-4-создаём-namespace-и-опционально-секрет-реестра)
9. [Шаг 5. Готовим свой values-файл](#9-шаг-5-готовим-свой-values-файл)
10. [Шаг 6. Установка чарта](#10-шаг-6-установка-чарта)
11. [Шаг 7. Проверяем, что всё работает](#11-шаг-7-проверяем-что-всё-работает)
12. [Обновление, откат, удаление](#12-обновление-откат-удаление)
13. [Частые ошибки и как чинить](#13-частые-ошибки-и-как-чинить)
14. [Полная справка по values](#14-полная-справка-по-values)

---

## 1. Что такое Helm и зачем он нужен

В Kubernetes описание сервиса — это набор YAML-файлов: `Deployment`, `Service`,
`Secret`, `ConfigMap`, `Ingress`. Если каждый раз писать их руками — это:

- много копипасты между сервисами;
- сложно параметризовать (dev, stage, prod);
- неудобно обновлять (нужно следить, что нигде не забыл).

**Helm** — пакетный менеджер для Kubernetes. Он:

- Описание ресурсов хранит в **шаблонах** (`templates/`), а параметры — в `values.yaml`.
- Один комплект шаблонов + параметры = **chart**.
- При установке чарт превращается в **release** — конкретно установленный экземпляр чарта в кластере.
- Помнит историю релизов, умеет откатывать, обновлять, удалять.

Аналогии:
- chart ≈ образ Docker (шаблон),
- release ≈ контейнер из этого образа (запущенный экземпляр),
- `values.yaml` ≈ переменные окружения / `--build-arg`.

---

## 2. Структура чарта

```
helm/course-service/
├── Chart.yaml                 # метаданные: имя, версия чарта, версия приложения
├── values.yaml                # значения по умолчанию (с подробными комментариями)
├── values-prod.example.yaml   # пример override для прод-стенда
├── README.md                  # этот файл
├── .helmignore                # что не упаковывать в архив чарта
└── templates/
    ├── _helpers.tpl           # вспомогательные шаблоны (имена, метки)
    ├── deployment.yaml        # как запускать поды
    ├── service.yaml           # стабильный сетевой адрес для подов
    ├── secret.yaml            # содержимое .env (с паролем БД)
    ├── ingress.yaml           # публикация наружу по HTTP-имени (опционально)
    ├── serviceaccount.yaml    # «учётка» подов внутри кластера
    └── NOTES.txt              # что показывается после `helm install`
```

### Важная особенность приложения

Сервис написан так, что обязательно читает файл `.env` через `godotenv.Load()`
(см. `internal/config/config.go`). В Kubernetes мы:

1. Кладём содержимое `.env` в **Secret** (потому что там пароль БД).
2. Монтируем этот Secret внутрь контейнера как файл `/app/.env`.

То есть приложение даже не подозревает, что работает в k8s — для него это просто файл.

---

## 3. Что мы создаём в Kubernetes

При `helm install` Helm создаст в namespace следующие объекты:

| Объект           | Зачем                                                                       |
|------------------|-----------------------------------------------------------------------------|
| `Deployment`     | Управляет подами: следит, чтобы их было ровно `replicaCount`, обновляет образ. |
| `Pod`(ы)         | Сам контейнер с приложением (поды создаёт Deployment, мы их не описываем).  |
| `Service`        | Стабильный DNS/IP, чтобы ходить на поды (балансирует трафик между ними).    |
| `Secret`         | Содержит `.env` (включая пароль). Монтируется в под как файл.               |
| `ServiceAccount` | «Учётка» подов в k8s API (не используется приложением, но полезно иметь).   |
| `Ingress`        | (опционально) Публикация наружу по HTTP-имени.                              |

Связи (упрощённо):

```
       Ingress  ─►  Service  ─►  Pod (контейнер)
                                    ▲
                                    │ монтируется как /app/.env
                                  Secret
```

---

## 4. Подготовка инструментов

На своём компьютере должны стоять:

| Инструмент | Версия           | Команда проверки                  |
|------------|------------------|-----------------------------------|
| `kubectl`  | ≥ 1.24           | `kubectl version --client`        |
| `helm`     | ≥ 3.11           | `helm version`                    |
| `docker`   | любая свежая     | `docker --version`                |

Установка:
- **kubectl**: <https://kubernetes.io/docs/tasks/tools/>
- **helm**:    <https://helm.sh/docs/intro/install/>
- **docker**:  <https://docs.docker.com/get-docker/>

Проверить, что `kubectl` видит кластер:

```bash
kubectl cluster-info
kubectl get nodes
```

Должен показать хотя бы один узел в состоянии `Ready`.

---

## 5. Шаг 1. Поднимаем кластер Kubernetes

Если кластера ещё нет — для учебных целей хватает любого локального:

### Вариант A. Docker Desktop (Windows / macOS)

В настройках Docker Desktop включить «Enable Kubernetes». Внутри уже едет
готовый кластер из одной ноды. `kubectl` будет настроен автоматически.

**Лайфхак для Windows**: чтобы из пода достучаться до сервисов на хосте
(там у нас Postgres и Kafka), используется DNS-имя `host.docker.internal`.
Именно оно по умолчанию прописано в `values.yaml`.

### Вариант B. minikube

```bash
minikube start --driver=docker --cpus=2 --memory=4g
```

Чтобы из подов был доступ к хосту: использовать `host.minikube.internal`
(нужно `minikube addons enable host-resolver` или прописать вручную).

### Вариант C. kind (Linux)

На **Linux** имя `host.docker.internal` **по умолчанию не существует** (в отличие от Docker Desktop на Windows/macOS). Поэтому с дефолтным `values.yaml` Postgres из Docker часто недоступен.

**Рекомендуемый способ** — подключить контейнер Postgres к docker-сети `kind`:

```bash
kind create cluster --name course

cd testenv && docker compose up -d postgres
# имя контейнера из compose — postgres (см. docker ps)

docker network connect kind postgres
# при необходимости то же для kafka:
# docker network connect kind <имя-контейнера-kafka>
```

Дальше в Helm укажи `config.postgres.host: "postgres"` — это имя контейнера,
которое резолвится из подов kind через Docker DNS.

Готовый файл настроек: [`values-kind-linux.example.yaml`](values-kind-linux.example.yaml)

```bash
# из корня репозитория
docker build -t course-service:dev .
kind load docker-image course-service:dev --name course

cp helm/course-service/values-kind-linux.example.yaml helm/course-service/values-kind-linux.yaml
# при необходимости поправь пароль / имя контейнера

helm upgrade --install course-service ./helm/course-service \
  -n course --create-namespace \
  -f helm/course-service/values-kind-linux.yaml
```

**Альтернатива** — не подключать контейнер к сети `kind`, а ходить на Postgres через порт хоста (`-p 5432:5432`):

1. Создай кластер с `extraHosts` — см. [`kind-cluster.example.yaml`](kind-cluster.example.yaml).
2. В values включи `hostAliases` и `postgres.host: host.docker.internal` (пример в конце `values-kind-linux.example.yaml`).
3. IP в `hostAliases` подбери под свою машину: `ip -4 addr show docker0` (часто `172.17.0.1`, но не гарантировано).

### Вариант D. Реальный кластер (k3s, managed k8s, on-prem)

Получаешь от админа `kubeconfig` файл, кладёшь его в `~/.kube/config`
(или экспортируешь `KUBECONFIG=...`). Дальше `kubectl get nodes` должен работать.

---

## 6. Шаг 2. Готовим внешние Postgres и Kafka

Сервис ждёт, что:

- **PostgreSQL** уже создан и доступен по сети из подов кластера. БД, указанная
  в `config.postgres.db`, должна существовать (миграции её сами не создадут —
  они создают только таблицы внутри). Пользователь должен иметь права на CREATE/ALTER.

- **Kafka** опциональна. Если `config.kafka.brokers` пустой массив — приложение
  стартует без неё (см. `cmd/app/main.go`). Если указан — брокеры должны быть
  доступны до старта.

### Самый простой вариант для учёбы (Postgres на твоём ноуте через docker-compose)

В репозитории уже есть `testenv/docker-compose.yml`. Подними только Postgres:

```bash
cd testenv
echo "POSTGRES_USER=postgres"     >  .env
echo "POSTGRES_PASSWORD=postgres" >> .env
echo "POSTGRES_DB=courses"        >> .env
docker compose up -d postgres
```

Проверь:

```bash
docker exec -it postgres psql -U postgres -l
```

Должна быть БД `courses`. В `values.yaml` оставь `host: "host.docker.internal"`
(на Docker Desktop / Windows это адрес твоего хоста изнутри подов).

### Если Postgres стоит на отдельной машине

Просто пропиши её IP/DNS в `config.postgres.host` своего override-файла.
Убедись, что на сервере PG разрешены подключения с сети кластера
(`pg_hba.conf` + `listen_addresses` в `postgresql.conf`).

---

## 7. Шаг 3. Собираем и публикуем Docker-образ

Helm chart **сам не собирает образ** — он только запускает уже готовый.

### Если используешь Docker Desktop / minikube без приватного реестра

В Docker Desktop кластер видит образы локального docker-демона напрямую.
Достаточно собрать локально:

```bash
docker build -t course-service:dev .
```

Затем в `values.yaml` (или своём override-файле) указать:

```yaml
image:
  repository: course-service
  tag: dev
  pullPolicy: Never   # ВАЖНО: не пытаться тянуть из реестра
```

Для **minikube** нужно собрать образ внутри его docker-демона:

```bash
eval $(minikube docker-env)        # переключает докер-клиент на minikube
docker build -t course-service:dev .
```

…а в values:

```yaml
image: { repository: course-service, tag: dev, pullPolicy: Never }
```

### Если есть приватный/публичный реестр

```bash
docker build -t registry.example.com/myteam/course-service:1.0.3 .
docker push  registry.example.com/myteam/course-service:1.0.3
```

И в values:

```yaml
image:
  repository: registry.example.com/myteam/course-service
  tag: "1.0.3"
  pullPolicy: IfNotPresent
imagePullSecrets:
  - name: regcred   # см. шаг 4
```

---

## 8. Шаг 4. Создаём namespace и (опционально) секрет реестра

**Namespace** — это «папка» в кластере, в которой живут наши объекты.
Хорошая практика — для каждого сервиса/окружения свой:

```bash
kubectl create namespace course
```

(Можно не создавать заранее: при `helm install` укажем `--create-namespace`.)

Если реестр **приватный**, в этом namespace нужен секрет с credentials,
чтобы k8s мог тянуть образ:

```bash
kubectl create secret docker-registry regcred \
  --namespace=course \
  --docker-server=registry.example.com \
  --docker-username=USER \
  --docker-password=PASSWORD \
  --docker-email=you@example.com
```

И в `values.yaml` указать:

```yaml
imagePullSecrets:
  - name: regcred
```

Для публичных образов / локалки это не нужно.

---

## 9. Шаг 5. Готовим свой values-файл

Никогда не редактируй `values.yaml` напрямую — это «дефолты, годные для всех».
Создай **свой override-файл**, например `values-local.yaml` рядом
(или вне репо — особенно если там пароли):

```yaml
# values-local.yaml — для локальной разработки в Docker Desktop

image:
  repository: course-service
  tag: dev
  pullPolicy: Never           # образ собран локально

config:
  logLevel: debug
  postgres:
    host: "host.docker.internal"
    port: "5432"
    user: "postgres"
    db: "courses"
    sslmode: "disable"
  kafka:
    brokers: []               # без kafka

secrets:
  postgresPassword: "postgres"

service:
  type: NodePort              # удобно открыть наружу
```

Для production — посмотри `values-prod.example.yaml` рядом.

> **Безопасность.** Файлы с паролями **не коммить в git**. Либо храни их вне репо,
> либо используй `secrets.existingSecret: <name>` и создавай Secret отдельно
> (вручную, External Secrets Operator, sops, sealed-secrets).

---

## 10. Шаг 6. Установка чарта

Проверь, что шаблоны рендерятся без ошибок (это **dry-run**, в кластер ничего не пойдёт):

```bash
helm lint   ./helm/course-service
helm template course-service ./helm/course-service -f helm/course-service/values-local.yaml
```

Если в выводе нет ошибок — ставим:

```bash
helm upgrade --install course-service ./helm/course-service \
  --namespace course --create-namespace \
  -f helm/course-service/values-local.yaml
```

Что делает `upgrade --install`:
- если релиз `course-service` ещё не установлен — установит;
- если уже стоит — обновит.

Идемпотентно: можно гонять много раз, лишних объектов не наплодит.

После команды Helm выведет содержимое `NOTES.txt` — там подсказки, как достучаться.

---

## 11. Шаг 7. Проверяем, что всё работает

```bash
# Поды (должны быть в статусе Running, READY 1/1)
kubectl get pods -n course -l app.kubernetes.io/instance=course-service

# Логи приложения
kubectl logs -n course -l app.kubernetes.io/instance=course-service -f

# Описание пода (если падает — самое полезное)
kubectl describe pod -n course -l app.kubernetes.io/instance=course-service
```

В логах должны быть строчки:
```
successfully ran migrations
http server started at :10001
course service started
```

### Проверить /health

Если `service.type: ClusterIP` (по умолчанию), используй проброс порта:

```bash
kubectl port-forward -n course svc/course-service 10001:10001
# В другом терминале:
curl http://localhost:10001/health
# 200 OK = здоров. 503 = ping-и БД/Kafka валятся.
```

Если `service.type: NodePort`:

```bash
kubectl get svc -n course course-service
# Видишь колонку PORT(S) вида 10001:31234/TCP — 31234 это NodePort
curl http://<IP-узла>:31234/health
```

### Проверить, что .env правильный (без раскрытия пароля)

```bash
kubectl exec -n course deploy/course-service -- ls -la /app/.env
kubectl exec -n course deploy/course-service -- wc -l /app/.env
```

Чтобы посмотреть содержимое (осторожно, светит пароль):

```bash
kubectl exec -n course deploy/course-service -- cat /app/.env
```

---

## 12. Обновление, откат, удаление

### Обновление

После правки кода / values:

```bash
# Если поменялся код — пересобрать и запушить образ
docker build -t registry.example.com/myteam/course-service:1.0.4 .
docker push  registry.example.com/myteam/course-service:1.0.4

# Обновить релиз
helm upgrade course-service ./helm/course-service \
  -n course \
  -f helm/course-service/values-prod.yaml \
  --set image.tag=1.0.4
```

Helm сделает rolling update (по 1 поду): новый поднимет → дождётся readiness →
выключит старый. Без даунтайма (если `replicaCount > 1`).

### История и откат

```bash
helm history  course-service -n course
helm rollback course-service 2 -n course   # откатиться на ревизию №2
```

### Удалить релиз целиком

```bash
helm uninstall course-service -n course
# Удалить и namespace:
kubectl delete namespace course
```

---

## 13. Частые ошибки и как чинить

### `ImagePullBackOff` / `ErrImagePull`
k8s не может скачать образ.
- Неправильный `image.repository` или `tag` → проверь, что такой образ реально есть в реестре.
- Реестр приватный → создай `imagePullSecrets` (см. шаг 4).
- Локальный образ + minikube → собрал ли через `eval $(minikube docker-env)`?
- Локальный образ + Docker Desktop → `pullPolicy: Never`?

### `CrashLoopBackOff`
Контейнер падает сразу. Смотри логи: `kubectl logs ... --previous`.
Чаще всего:
- `failed to connect to postgres` — неправильный `postgres.host`/порт/пароль или БД недоступна из кластера.
- `failed to run migrations` — у пользователя нет прав / БД не существует.
- `load .env file` — Secret не примонтировался (редко; обычно из-за ручной правки шаблонов).

### `runAsNonRoot` и `non-numeric user (app)`
Сообщение: `container has runAsNonRoot and image has non-numeric user (app), cannot verify user is non-root`.

Kubernetes требует **числовой** UID в образе, если в Pod задан `runAsNonRoot: true`.
Пересобери образ из актуального `Dockerfile` (там `USER 1000:1000`) и обнови релиз:

```bash
docker build -t course-service:dev .
kind load docker-image course-service:dev --name course
helm upgrade course-service ./helm/course-service -n course -f values-kind-linux.yaml
```

### Postgres недоступен на Linux + kind
- `host.docker.internal` не резолвится → используй `docker network connect kind postgres` и `host: postgres`.
- Ручной IP из `docker inspect` после пересоздания контейнера меняется → лучше имя в сети `kind`, а не IP.
- Проверка из пода: `kubectl exec -n course deploy/course-service -- wget -qO- --timeout=2 host.docker.internal:5432` (или `nc -zv postgres 5432`, если есть в образе).

### Поды `Pending`
`kubectl describe pod ...` → ищи Events. Чаще всего:
- кластеру не хватает CPU/RAM → уменьшить `resources.requests`;
- нет узла, удовлетворяющего `nodeSelector`/`affinity`.

### Сервис не открывается снаружи
- `service.type: ClusterIP` — он и не должен; используй port-forward или Ingress/NodePort.
- `Ingress` есть, но не работает → в кластере нет ingress-controller'а или DNS-имя не направлено на кластер.

### Пароль в `.env` поменял, а сервис продолжает использовать старый
Проверь, что в Deployment есть аннотация `checksum/secret`. Она должна автоматически
пересоздавать поды при изменении Secret. Если её нет — `kubectl rollout restart deploy/course-service -n course`.

---

## 14. Полная справка по values

Тут краткая сводка. Полные комментарии — в `values.yaml`.

| Параметр                          | Тип             | Дефолт                  | Описание |
|-----------------------------------|-----------------|-------------------------|----------|
| `image.repository`                | string          | `course-service`        | Репо образа без тега |
| `image.tag`                       | string          | `""` (= `appVersion`)   | Тег образа |
| `image.pullPolicy`                | string          | `IfNotPresent`          | `Always` / `IfNotPresent` / `Never` |
| `imagePullSecrets`                | list            | `[]`                    | Секреты для приватного реестра |
| `replicaCount`                    | int             | `1`                     | Сколько копий пода |
| `strategy.type`                   | string          | `RollingUpdate`         | Стратегия обновления |
| `config.logLevel`                 | string          | `info`                  | Уровень логирования |
| `config.logFilePath`              | string          | `/app/logs`             | Папка для лог-файла |
| `config.httpServerAddr`           | string          | `:10001`                | Адрес HTTP-сервера |
| `config.postgres.host`            | string          | `host.docker.internal`  | Адрес PG (внешний!) |
| `config.postgres.port`            | string          | `5432`                  | Порт PG |
| `config.postgres.user`            | string          | `postgres`              | Пользователь PG |
| `config.postgres.db`              | string          | `courses`               | Имя БД |
| `config.postgres.sslmode`         | string          | `disable`               | `disable` / `require` / `verify-ca` / `verify-full` |
| `config.kafka.brokers`            | list[string]    | `[]`                    | Список брокеров (если пуст — kafka не используется) |
| `secrets.postgresPassword`        | string          | `postgres`              | Пароль PG. **Меняй и не коммить.** |
| `secrets.existingSecret`          | string          | `""`                    | Имя существующего Secret (вместо создания нового) |
| `service.type`                    | string          | `ClusterIP`             | `ClusterIP` / `NodePort` / `LoadBalancer` |
| `service.port`                    | int             | `10001`                 | Порт Service |
| `service.nodePort`                | int             | (auto)                  | Если type=NodePort |
| `ingress.enabled`                 | bool            | `false`                 | Включить Ingress |
| `ingress.className`               | string          | `nginx`                 | IngressClass |
| `ingress.hosts`                   | list            | пример                  | Список host'ов и path'ов |
| `ingress.tls`                     | list            | `[]`                    | TLS-секции |
| `resources.requests/limits`       | map             | small defaults          | CPU/RAM |
| `probes.startup/liveness/readiness` | map           | enabled                 | Проверки на `/health` |
| `podSecurityContext`              | map             | `runAsNonRoot: true`    | Безопасность пода |
| `containerSecurityContext`        | map             | drop ALL caps           | Безопасность контейнера |
| `serviceAccount.create`           | bool            | `true`                  | Создавать SA |
| `serviceAccount.name`             | string          | (auto)                  | Имя SA (если задано) |
| `nodeSelector` / `tolerations` / `affinity` | map/list | пусто                | Размещение подов |
| `hostAliases`                     | list            | `[]`                    | Записи в /etc/hosts пода (Linux + kind) |
| `extraEnv`                        | list            | `[]`                    | Доп. env-переменные в контейнер |

---

## Дальнейшие шаги, которые имеет смысл добавить позже

Это **не нужно** для учебного стенда, но полезно знать на будущее:

- **Логи в stdout вместо файла** — изменить `pkg/logger/logger.go`, чтобы при пустом
  `LOG_FILE_PATH` логи шли только в stdout. В k8s это «правильно», логи собирают
  agent'ы (fluent-bit, vector).
- **Опциональный `.env`** — поправить `internal/config/config.go`, чтобы
  отсутствие `.env` не было ошибкой. Тогда можно будет переключиться с Secret-as-file
  на классический `envFrom: secretRef`.
- **HorizontalPodAutoscaler** — автоскейлинг по CPU/RAM.
- **PodDisruptionBudget** — гарантия минимума работающих подов при обслуживании узлов.
- **NetworkPolicy** — ограничение, кому разрешено стучаться в наши поды.
- **External Secrets / Vault** — нормальное хранение паролей вне git.
- **Job для миграций** — отдельный k8s Job, который накатывает миграции до старта подов
  (сейчас миграции бегут на старте каждого пода, при `replicaCount > 1` могут гоняться).
