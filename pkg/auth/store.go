package auth

import "log"

type authStoreCache map[string]*authInfo // token -> user scope

type authInfo struct {
	ExpiredTime string
	Username    string
	Scope       map[string]string // cluster/namespace/pod => action{exec,log,debug...}
}

var (
	storeCache authStoreCache
)

func init() {
	storeCache = make(map[string]*authInfo, 5) // init size 5
}

func (s authStoreCache) get(token string) *authInfo {
	//fmt.Println("get privilege from store")

	if _, ok := s[token]; !ok {
		return &authInfo{
			Scope: map[string]string{},
		}
	}
	return s[token]
}

func (s authStoreCache) save(token string, username string, scope map[string]string) {
	//fmt.Println("save privilege to store")

	if _, ok := s[token]; !ok {
		s[token] = &authInfo{
			Username: username,
			Scope:    scope,
		}
		return
	}

	s[token].Username = username
	for k, v := range scope {
		s[token].Scope[k] = v
	}
}

func  Get(token string) *authInfo {
	log.Println("all cache:",storeCache)
	return storeCache.get(token)
}

func Set(token string, username string) bool {

	if _, ok := storeCache[token]; !ok {
		return false
	}
	storeCache[token].Username = username
	return true
}



