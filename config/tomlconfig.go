package config

type TomlConfig struct {
	BotToken              string
	AdminId               int
	Debug                 bool
	TorrentDownloadPath   string
	ChatIdToKickUsersFrom int64
	BitBucketRepoUser     string
	BitBucketRepo         string
	BitBucketUser         string
	BitBucketPassword     string
	BitBucketTelegramChat int64
}
