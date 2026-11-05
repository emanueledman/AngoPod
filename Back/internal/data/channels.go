package data

import "angopod-back/internal/youtube"

// ConfiguredChannels lista canais angolanos pré-configurados para sincronização
// Nota: Os ChannelIDs podem ser atualizados ou expandidos conforme necessário
var ConfiguredChannels = []youtube.ChannelConfig{
	{
		ChannelID: "UC_wO4S9WdJp098tZ-Xm3k7A", // Fly Podcast AO
		Slug:      "fly-podcast",
		Host:      "Flay Squad",
		Category:  "Sociedade",
		Province:  "Luanda",
		Language:  "Português",
		Color:     "coral",
	},
	{
		ChannelID: "UCjYtD4837593h9sd8239sdA", // Dinheiro Limpo
		Slug:      "dinheiro-limpo",
		Host:      "Euclides Francisco",
		Category:  "Negócios",
		Province:  "Luanda",
		Language:  "Português",
		Color:     "purple",
	},
	{
		ChannelID: "UC9sd82348df82347293847A", // Kota, Me Conta
		Slug:      "kota-me-conta",
		Host:      "Elisa Manuel",
		Category:  "Cultura",
		Province:  "Benguela",
		Language:  "Umbundu",
		Color:     "yellow",
	},
	{
		ChannelID: "UC823947239847293847239A", // Negócios na Banda
		Slug:      "negocios-na-banda",
		Host:      "Edson Mateus",
		Category:  "Negócios",
		Province:  "Luanda",
		Language:  "Português",
		Color:     "green",
	},
	{
		ChannelID: "UC238472938472938472394A", // Mwangolé Tech
		Slug:      "mwangole-tech",
		Host:      "Tânia Gaspar",
		Category:  "Tecnologia",
		Province:  "Luanda",
		Language:  "Português",
		Color:     "blue",
	},
	{
		ChannelID: "UC348729384729384729384A", // Batidas da Banda
		Slug:      "batidas-da-banda",
		Host:      "Lina Santos",
		Category:  "Música",
		Province:  "Luanda",
		Language:  "Português",
		Color:     "lime",
	},
	{
		ChannelID: "UC472938472938472938472A", // Ideias do Sul
		Slug:      "ideias-do-sul",
		Host:      "Mauro Gomes",
		Category:  "Tecnologia",
		Province:  "Huíla",
		Language:  "Português",
		Color:     "coral",
	},
}

