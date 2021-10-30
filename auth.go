package main

type GetMessage struct {
	Get string `json:"get"`
}

type ChallengeMessage struct {
	Challenge Challenge `json:"challenge"`
}

type Challenge struct {
	Txt string `json:"txt"`
	Sig string `json:"sig"`
}

type ChallengeResponse struct {
	PublicKey string    `json:"publicKey"`
	Challenge Challenge `json:"challenge"`
	Solution  string    `json:"solution"`
}
