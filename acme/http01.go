package acme

var HTTP01Provider = HTTP01ChallengesProvider{challenges: make(map[string]HTTP01ChallengeData)}

type HTTP01ChallengeData struct {
	Domain  string
	Token   string
	KeyAuth string
}

type HTTP01ChallengesProvider struct {
	challenges map[string]HTTP01ChallengeData
}

func (p HTTP01ChallengesProvider) Present(domain, token, keyAuth string) error {
	p.challenges[token] = HTTP01ChallengeData{
		Domain:  domain,
		Token:   token,
		KeyAuth: keyAuth,
	}

	return nil
}

func (p HTTP01ChallengesProvider) CleanUp(domain, token, keyAuth string) error {
	delete(p.challenges, token)
	return nil
}

func (p HTTP01ChallengesProvider) GetChallenge(token string) (*HTTP01ChallengeData, bool) {
	chal, exists := p.challenges[token]
	return &chal, exists
}
