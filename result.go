package meekstv

import (
	"github.com/shawntoffel/election"
)

func (m *meekStv) Result() (*election.Result, error) {

	if m.Error != nil {
		return nil, m.Error
	}

	elected := m.Pool.Elected()

	result := election.Result{}
	result.Events = m.Events
	result.Candidates = elected.AsCandidates()

	return &result, nil
}
