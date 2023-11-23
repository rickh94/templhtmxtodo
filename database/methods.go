package database

import "github.com/go-webauthn/webauthn/webauthn"

func (user *User) WebAuthnID() []byte {
	return []byte(user.ID)
}

func (user *User) WebAuthnName() string {
	if user.Name != "" {
		return user.Name
	} else {
		return user.Email
	}
}

func (user *User) WebAuthnDisplayName() string {
	return string(user.WebAuthnName())
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	var credentials []webauthn.Credential
	for _, credential := range user.Credentials {
		credentials = append(credentials, webauthn.Credential{
			ID:              credential.CredentialID,
			PublicKey:       credential.PublicKey,
			AttestationType: credential.AttestationType,
			Transport:       credential.Transport,
			Flags:           credential.Flags,
			Authenticator:   credential.Authenticator,
		})
	}
	return credentials
}

func (user *User) WebAuthnIcon() string {
	return ""
}

func (user *User) AddCredential(newCredential *webauthn.Credential) error {
	var credential Credential
	credential.CredentialID = newCredential.ID
	credential.PublicKey = newCredential.PublicKey
	credential.Transport = newCredential.Transport
	credential.AttestationType = newCredential.AttestationType
	credential.Flags = newCredential.Flags
	credential.Authenticator = newCredential.Authenticator
	credential.UserID = user.ID

	result := DB.Create(&credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeletePasskeys(user *User) error {
	result := DB.Where("user_id = ?", user.ID).Delete(&Credential{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CountUserCredentials(user *User) (int64, error) {
	var count int64
	result := DB.Model(&Credential{}).Where("user_id = ?", user.ID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func GetUserTodos(user *User) ([]Todo, error) {
	var todos []Todo
	result := DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}

func CreateTodo(text string, user *User) (*Todo, error) {
	todo := Todo{Text: text, UserID: user.ID}
	result := DB.Create(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func DeleteTodo(todoID string, user *User) error {
	result := DB.Where("id = ? and user_id = ?", todoID, user.ID).Delete(&Todo{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CompleteTodo(todoID string, user *User) (*Todo, error) {
	var updatedTodo Todo
	result := DB.Model(&updatedTodo).Where("id = ? and user_id = ?", todoID, user.ID).Update("completed", true)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedTodo, nil
}

func UnCompleteTodo(todoID string, user *User) (*Todo, error) {
	var updatedTodo Todo
	updatedTodo.ID = todoID
	result := DB.Model(&updatedTodo).Where("id = ? and user_id = ?", todoID, user.ID).Update("completed", false)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedTodo, nil
}

func GetTodoByID(todoID string, user *User) (*Todo, error) {
	var todo Todo
	result := DB.Where("id = ? and user_id = ?", todoID, user.ID).First(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func UpdateTodoText(todoID string, newText string, user *User) (*Todo, error) {
	var updatedTodo Todo
	updatedTodo.ID = todoID
	result := DB.Model(&updatedTodo).Where("id = ? and user_id = ?", todoID, user.ID).Update("text", newText)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedTodo, nil
}
