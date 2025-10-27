package domain

type Info struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Subject  string   `json:"subject"`
	Issuer   string   `json:"issuer"`
	Audience []string `json:"audience"`
}

func (inf *Info) FromMap(m map[string]any) {
	inf.Username = m["username"].(string)
	inf.Password = m["password"].(string)
	inf.Subject = m["subject"].(string)
	inf.Issuer = m["issuer"].(string)
	inf.Audience = make([]string, len(m["audience"].([]any)))
	for i, aud := range m["audience"].([]any) {
		inf.Audience[i] = aud.(string)
	}
}
