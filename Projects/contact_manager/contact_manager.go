package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var contactStorage string = "contact.json"
var trashStorage string = "trash.json"

type Contact struct {
	Name     string    `json:"name"`
	Phone    string    `json:"phone_number"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

type TrashEntry struct {
    Contact
    DeleteAt time.Time `json:"expire_at"`
}

type Manager struct {
	contacts map[string]Contact
	trash    map[string]TrashEntry
	stack    *Stack
}

func NewManager() *Manager {
	return &Manager{
		contacts: make(map[string]Contact),
		trash:    make(map[string]TrashEntry),
        stack: NewStack(),
	}
}

func validateContact(contact Contact) error {
    if contact.Phone == "" || contact.Email == "" {
        return errors.New("phone and email field are required")
    }
    return nil
}

func (m *Manager) AddNewContact(contact Contact) error {
    err := validateContact(contact)
    if err != nil {
        return err
    }
    err = m.AddNewContactImplement(contact)
    if err != nil {
        return err
    }
    m.stack.Push(Action{
        Undo: func() error {
            return  m.RemoveContactByEmailImplement(contact.Email)
        },
    })
    return nil
}

func (m *Manager) AddNewContactImplement(contact Contact) error {
	if _, ok := m.contacts[contact.Email]; ok {
		return errors.New("contact exists in database, " + contact.Email)
	}
	m.contacts[contact.Email] = contact
	return nil
}

func (m *Manager) RemoveContactByEmail(email string) error {
    err := m.RemoveContactByEmailImplement(email)
    if err != nil {
        return err
    }
    m.stack.Push(Action{
        Undo: func() error {
            return  m.RestoreContactByEmailImplement(email)
        },
    })
    return nil
 }

func (m *Manager) RemoveContactByEmailImplement(email string) error {
	if _, ok := m.contacts[email]; !ok {
		return errors.New("no contact found for " + email)
	}
	m.trash[email] = TrashEntry{
        Contact:  m.contacts[email],
        DeleteAt: time.Now(),
    }
	delete(m.contacts, email)
	return nil
}

func (m *Manager) GetContactByEmail(email string) (Contact, error) {
	if value, ok := m.contacts[email]; ok {
		return value, nil
	}
	return Contact{}, errors.New("no contact found for " + email)
}

func (m *Manager) GetContactByName(name string) []Contact {
	contacts := make([]Contact, 0)
	for _, value := range m.contacts {
		if strings.Contains(strings.ToLower(value.Name), strings.ToLower(name)) {
			contacts = append(contacts, value)
		}
	}
	return contacts
}

func (m *Manager) ListContacts() []Contact {
	contacts := make([]Contact, 0, len(m.contacts))
	for _, contact := range m.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}

func (m *Manager) ListTrash() []TrashEntry {
	trashContacts := make([]TrashEntry, 0, len(m.trash))
	for _, contact := range m.trash {
		trashContacts = append(trashContacts, contact)
	}
	return trashContacts
}

func (m *Manager) RestoreContactByEmail(email string) error {
    err := m.RestoreContactByEmailImplement(email)
    if err != nil {
        return err
    }
    m.stack.Push(Action{
        Undo: func() error {
            return  m.RemoveContactByEmailImplement(email)
        },
    })
    return nil
}

func (m *Manager) RestoreContactByEmailImplement(email string) error {
	if _, ok := m.trash[email]; !ok {
		return errors.New("no contact found for " + email)
	}
	m.contacts[email] = m.trash[email].Contact
	delete(m.trash, email)
	return nil
}

func (m *Manager) RemoveContactByEmailFromTrash(email string) error {
	if _, ok := m.trash[email]; !ok {
		return errors.New("no contact found for " + email)
	}
	delete(m.trash, email)
	return nil
}

func (m *Manager) EmptyTrash() {
	m.trash = make(map[string]TrashEntry)
}

func (m *Manager) DeleteExpiredContactFromTrash(duration time.Duration) {
    now := time.Now()
    for email, entry := range m.trash {
        if now.Sub(entry.DeleteAt) >= duration {
            delete(m.trash, email)
        }
    }
}

func (m *Manager) SortContactByName() []Contact {
	contact := m.ListContacts()

	sort.Slice(contact, func(i, j int) bool {
		return contact[i].Name < contact[j].Name
	})
	return contact
}

func (m *Manager) SortContactByEmail() []Contact {
	contact := m.ListContacts()

	sort.Slice(contact, func(i, j int) bool {
		return contact[i].Email < contact[j].Email
	})
	return contact
}

func (m *Manager) HasBirthdayToday() []Contact {
	today := time.Now().In(time.Local).Format("01-02") 
	contacts := make([]Contact, 0)
	for _, value := range m.contacts {
        birthday := value.Birthday.Format("01-02")
		if birthday == today {
			contacts = append(contacts, value)
		}
	}
	return contacts
}

func (m *Manager) HasBirthdayNextSevenDays() []Contact {
    contacts := make([]Contact, 0)
    nextSevenDays := make(map[string]bool)
    for i := 1; i <= 7; i++ {
        next := time.Now().AddDate(0, 0, i).Format("01-02")
        nextSevenDays[next] = true
    }

	for _, value := range m.contacts {
        birthday := value.Birthday.Format("01-02")
		if nextSevenDays[birthday] {
			contacts = append(contacts, value)
		}
	}
	return contacts
}

func (m *Manager) Undo() error {
   return m.stack.Undo()
}

func (m *Manager) ImportContact() error {
	contacts, err := ImportFromJSON[Contact](contactStorage)
	if err != nil {
		return err
	}
	for _, value := range contacts {
		m.contacts[value.Email] = value
	}
	return nil
}

func (m *Manager) ImportTrash() error {
	contacts, err := ImportFromJSON[TrashEntry](trashStorage)
	if err != nil {
		return err
	}
	for _, value := range contacts {
		m.trash[value.Email] = value
	}
	return nil
}

func ImportFromJSON[T any](storageName string) ([]T, error) {
	file, err := os.Open(storageName)
	if err != nil {
		return []T{}, fmt.Errorf("failed to open %s file : %w", storageName, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if len(data) == 0 {
		return []T{}, nil
	}
	if err != nil {
		return []T{}, fmt.Errorf("failed to read %s file : %w", storageName, err)
	}

	var contacts []T
	if err := json.Unmarshal(data, &contacts); err != nil {
		return []T{}, fmt.Errorf("failed to unmarshal contact from %s file : %w", storageName, err)
	}

	return contacts, nil
}

func (m *Manager) ExportContact() error {
	contacts := m.ListContacts()
	return ExportToJSON(contacts, contactStorage)
}

func (m *Manager) ExportTrash() error {
	contacts := m.ListTrash()
	return ExportToJSON(contacts, trashStorage)
}

func ExportToJSON[T any](contacts []T, storageName string) error {
	data, err := json.MarshalIndent(contacts, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal contact from %s file : %w", storageName, err)
	}

	file, err := os.OpenFile(storageName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file : %w", storageName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write to %s file: %w", storageName, err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush to write %s file : %w", storageName, err)
	}

	return nil
}

func main() {

	cm := NewManager()

	// Import
	if err := cm.ImportContact(); err != nil {
		log.Panic(err)
	}
	if err := cm.ImportTrash(); err != nil {
		log.Panic(err)
	}

	// AddContactTest(cm)
	// ListContactsTest(cm)
	// RemoveContactByEmailTest(cm)

	// ListContactsTest(cm)
	// fmt.Println()

	// GetContactByEmail(cm)
	// GetContactByName(cm)

	// ListTrashTest(cm)

	// RestoreContactByEmailTest(cm)
	// ListContactsTest(cm)

	// RemoveContactByEmailTest(cm)
	// RemoveContactByEmailFromTrashTest(cm)
	// ListTrashTest(cm)

	SortContactByEmailTest(cm)

	// if err := cm.Undo(); err != nil {
	// 	log.Println(err)
	// }

	// ListContactsTest(cm)


	// Export
	if err := cm.ExportContact(); err != nil {
		log.Panic(err)
	}

	if err := cm.ExportTrash(); err != nil {
		log.Panic(err)
	}
}
