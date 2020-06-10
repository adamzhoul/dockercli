package auth

import (
	"strings"
)

// return username, passOrNot
func CheckUser(token string, resource string, action string) (string, bool) {

	// 1. check if we have this  specific info in store cache
	a := storeCache.get(token)
	if a.Username != "" {
		if _, ok := a.Scope[resource]; ok {
			if strings.Contains(a.Scope[resource], action) {
				return a.Username, true
			} else {
				go syncLatestPrivilege(token, resource)
				return a.Username, false //
			}
		}
	}

	// 2. if don't,  ask server
	syncLatestPrivilege(token, resource)
	a = storeCache.get(token)

	// 3. return if can continue
	if !strings.Contains(a.Scope[resource], action) {
		return a.Username, false //
	}

	return a.Username, true
}

func syncLatestPrivilege(token string, resource string) {

	//fmt.Println("sync privilege to store")
	// 1. limit send window
	//  same token ask less than 15 times in one minute
	//  less than 200 times in one minute for total
	//  2 purpose: 1 is to avoid too goroutine created in short time; 2 for the good of server (can be done in sidecar)

	// 2. get from server
	username, actions := getAuth(token, resource)

	// 3. store user  privilege
	scope := map[string]string{
		resource: actions,
	}
	storeCache.save(token, username, scope)
}
