package requests

import "errors"

type RemoteGameForm struct {
	GameName   *string `json:"gameName"` 
	SpyNumber  int     `json:"spyNumber"`
	MaxNumbers int     `json:"maxNumbers"`
	IsPrivate  bool    `json:"isPrivate"`
}

func (f RemoteGameForm) Validate() error {
if f.SpyNumber <= 0 || f.SpyNumber > 2 {
		return errors.New("spyNumber must be at least 1 and no bigger than 2")
	}
	//nolint
	if f.MaxNumbers < 2 {
		return errors.New("maxNumbers must be at least 2")
	}
	if f.SpyNumber > f.MaxNumbers {
		return errors.New("spyNumber cannot exceed maxNumbers")
	}
	if f.GameName != nil && len(*f.GameName) > 127 {
		return errors.New("gameName too long (max 127)")
	}
	return nil
}
