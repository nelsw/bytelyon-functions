package model

type Work struct {
	Job   Job   `json:"job"`
	Items Items `json:"items"`
}

func MakeWork(j Job, items Items) Work {
	return Work{
		Job:   j,
		Items: items,
	}
}

func (w Work) IsEmpty() bool {
	return len(w.Items) == 0
}
