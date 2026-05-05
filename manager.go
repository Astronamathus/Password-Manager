package main

import "strings"

type PasswordManager struct {
	credentials []Credential
}

func (pm *PasswordManager) Add(cred Credential) {
	pm.credentials = append(pm.credentials, cred)
}

func (pm *PasswordManager) GetAll() []Credential {
	return pm.credentials
}

func (pm *PasswordManager) Search(site string) *Credential {
	for i := range pm.credentials {
		if strings.EqualFold(pm.credentials[i].Site, site) {
			return &pm.credentials[i]
		}
	}
	return nil
}

func (pm *PasswordManager) Total() int {
	return len(pm.credentials)
}

func (pm *PasswordManager) Delete(site string) bool {
	for i:= range pm.credentials { 
		if strings.EqualFold(pm.credentials[i].Site, site) { 
			pm.credentials = append(pm.credentials[:i], pm.credentials[i+1:]...)
			return true
		}
	}
	return false

}