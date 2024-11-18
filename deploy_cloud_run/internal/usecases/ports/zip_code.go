package ports

type ZipCodeResponse struct {
	Cep          string `json:"cep"`
	Street       string `json:"logradouro"`
	Complement   string `json:"complemento"`
	Unit         string `json:"unidade"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	StateName    string `json:"estado"`
	Region       string `json:"regiao"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	AreaCode     string `json:"ddd"`
	Siafi        string `json:"siafi"`
	Error        string `json:"erro"`
}

type ZipCodeService interface {
	Get(code string) (*ZipCodeResponse, error)
}
