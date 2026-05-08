package main

import (
	"fmt"
	"strings"
)

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

func (pm *PasswordManager) SearchAll(site string) []Credential {
	var results []Credential

	for _, c := range pm.credentials {
		if strings.EqualFold(c.Site, site) {
			results = append(results, c)
		}
	}

	return results
}

func (pm *PasswordManager) Update(site, username, newPassword string) bool {
	for i := range pm.credentials {
		if strings.EqualFold(pm.credentials[i].Site, site) &&
			strings.EqualFold(pm.credentials[i].Username, username) {

			pm.credentials[i].Password = newPassword
			return true
		}
	}
	return false
}

func (pm *PasswordManager) Total() int {
	return len(pm.credentials)
}

func (pm *PasswordManager) Exists(site, username string) bool {
	for _, c := range pm.credentials {
		if strings.EqualFold(c.Site, site) &&
			strings.EqualFold(c.Username, username) {
			return true
		}
	}
	return false
}

func (pm *PasswordManager) UpdateUsername(site, oldUsername, newUsername string) bool {
	for i := range pm.credentials {
		if strings.EqualFold(pm.credentials[i].Site, site) &&
			strings.EqualFold(pm.credentials[i].Username, oldUsername) {

			// prevent duplicates
			for _, c := range pm.credentials {
				if strings.EqualFold(c.Site, site) &&
					(c.Username == newUsername) {
					fmt.Println("Username already exists for this site")
					return false
				}
			}

			pm.credentials[i].Username = newUsername
			return true
		}
	}
	return false
}

func (pm *PasswordManager) GroupBySite() map[string][]Credential {
	grouped := make(map[string][]Credential)

	for _, c := range pm.credentials {
		grouped[c.Site] = append(grouped[c.Site], c)
	}

	return grouped
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

func (pm *PasswordManager) DeleteExact(site, username string) bool {
	for i := 0; i < len(pm.credentials); i++ {
		if strings.EqualFold(pm.credentials[i].Site, site) &&
			strings.EqualFold(pm.credentials[i].Username, username) {

			pm.credentials = append(pm.credentials[:i], pm.credentials[i+1:]...)
			return true
		}
	}
	return false
}

func maskPassword(password string) string {
	return strings.Repeat("*", len(password))
}
func success(msg string) {
	fmt.Println(Green + msg + Reset)
}

func errorMsg(msg string) {
	fmt.Println(Red + msg + Reset)
}

func info(msg string) {
	fmt.Println(Cyan + msg + Reset)
}