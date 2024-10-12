#someday_maybe

`Makefile`
```makefile
backup:
	gobackup perform
```

`gobackup.yaml`
```yaml
{
  web: {
    host: "0.0.0.0",
    port: 2703,
    username: "rprtr258",
    password: "fgjepr8tghje3p865t4u9r253h2",
  },
  models: {
    gtd: {
      schedule: {
        cron: "0 0 * * *",
      },
      compress_with: {
        type: "tgz",
      },
      storages: {
        s3: {
          type: "s3",
          bucket: "backup--gtd",
          region: "ru-central1",
          access_key_id: "YCAJE4HHa8pyWzSI2GJC97vEn",
          secret_access_key: "YCOrOR1SXioPH0tVXFkFM_Z7hZ_zF9YWs2rO0NCh",
          endpoint: "storage.yandexcloud.net",
          keep: 1 * 7, // 1 per day, 1 * 7 per week
        },
      },
      archive: {
        includes: ["."],
      },
      notifiers: {
        telegram: {
          type: "telegram",
          chat_id: "310506774",
          token: "5415014383:AAElsIuS3hqpF0PUayvhsIznbbDMSB4Ioh8",
        },
      },
    },
  },
}
```