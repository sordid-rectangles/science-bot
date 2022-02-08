package auth

var Users struct {
	ADMINS map[string]string
	OWNER  map[string]string
}

func SetMaps(admins *map[string]string, owner *map[string]string) {
	Users.ADMINS = *admins
	Users.OWNER = *owner
}

func SetOwner(id string, name string) error {
	_, ok := Users.OWNER[id]
	if ok {
		return nil
	} else {
		Users.OWNER[id] = name
		//todo: also update the db
		return nil
	}
}

//TODO: func IsOwner (id string) (bool, error) {}
func IsOwner(id string) (bool, error) {
	_, ok := Users.OWNER[id]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

//TODO: func IsAdmin (id string) (bool, error) {}
func IsAdmin(id string) (bool, error) {
	_, ok := Users.ADMINS[id]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

//TODO: IsAuthed (id string) (bool, error) {}
func IsAuthed(id string) (bool, error) {
	own, _ := IsOwner(id)
	adm, _ := IsAdmin(id)
	if own || adm {
		return true, nil
	} else {
		return false, nil
	}
}

//TODO: RegisterAdmin (id string, name string) (error) {}
func RegisterAdmin(id string, name string) error {
	_, ok := Users.ADMINS[id]
	if ok {
		return nil
	} else {
		Users.ADMINS[id] = name
		//todo: also update the db
		return nil
	}
}

//TODO: RemoveAdmin (id string) (error) {}
func RemoveAdmin(id string) error {
	_, ok := Users.ADMINS[id]
	if ok {
		//todo: also update the db
		delete(Users.ADMINS, id)
		return nil
	} else {

		return nil
	}
}

//TODO: CurrentAdmins () (*map[string]string){}
func CurrentAdmins() (*map[string]string, error) {
	return &Users.ADMINS, nil
}

//DATABASE STORAGE OF CREDENTIALS

//TODO: storeCredential(id string, name string) (error)

//TODO: updateCredential(id string, name string) (error)

//TODO: removerCredential(id string) (error)
