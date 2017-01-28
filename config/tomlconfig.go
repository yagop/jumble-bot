package config

type TomlConfig struct {
	BotToken              string
	AdminId               int
	Degug                 bool
	TorrentDownloadPath   string
	ChatIdToKickUsersFrom int64
}
