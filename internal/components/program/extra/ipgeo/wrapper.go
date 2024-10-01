package ipgeo

type (
	Wrapper struct {
		Address string
	}
)

func (w *Wrapper) Country() string {
	return country(w.Address)
}

func (w *Wrapper) City() map[string]string {
	return city(w.Address)
}

func (w *Wrapper) CityEN() string {
	return cityEN(w.Address)
}
