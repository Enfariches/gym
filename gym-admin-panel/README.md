# Админ панель для Производственной Гимнастики.

Данный проект предназанчен для управления сервером производственной гимнастики. Его назанчение это загрузка и управление списоком видео для разминки, настройки расписания запуска видео на компьютерах сотрудников и сбора статистики по просмотру комплексов.

npm install -g protoc-gen-grpc-web

Stack: Quasar / Vue

## Install the dependencies
```bash
yarn
# or
npm install
```

### Start the app in development mode (hot-code reloading, error reporting, etc.)
```bash
quasar dev
```


### Lint the files
```bash
yarn lint
# or
npm run lint
```


### Format the files
```bash
yarn format
# or
npm run format
```


### Build the app for production
```bash
quasar build
```

### Test the app

```bash
# Run tests
npm run test

# Run unit tests
npm run test:unit

# Run tests with coverage
npm run test:coverage

# Watch for changes and run tests
npm run test:watch
```

## Test Coverage

The application includes unit tests for key components:

- Components with 100% coverage:
  - SmallCard
  - BigCard
  - MainPage (except for one branch)

- Components with partial coverage:
  - StatisticsChartComponent (36.3%)
  - statisticsStore (42.85%)

To improve test coverage, additional tests can be added for:
- Services
- Stores
- Other components
- Pages

### Customize the configuration
See [Configuring quasar.config.js](https://v2.quasar.dev/quasar-cli-vite/quasar-config-js).
