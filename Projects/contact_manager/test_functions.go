package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func AddContactTest(cm *Manager) {
	birthday, err := time.Parse("2006-01-02", "2000-10-18")
	if err != nil {
		log.Println("bidrthday is not valid")
	}

	for i := 0; i <= 10; i++ {
		err = cm.AddNewContact(Contact{
			Name:     "Labib-faisal",
			Phone:    "019xx-xxx-xxx",
			Email:    "labibfasisaltest" + strconv.Itoa(i) + "@gmail.com",
			Birthday: birthday,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func ListContactsTest(cm *Manager) {
	contacts := cm.ListContacts()
	for _, value := range contacts {
		fmt.Println(value.Email)
	}
}

func RemoveContactByEmailTest(cm *Manager) {
	err := cm.RemoveContactByEmail("labibtest1@gmail.com")
	if err != nil {
		log.Println(err)
	}
}

func GetContactByEmail(cm *Manager) {
	contact, err := cm.GetContactByEmail("labibtest2@gmail.com")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(contact)
}

func GetContactByName(cm *Manager) {
	contact := cm.GetContactByName("faisal")
	for _, value := range contact {
		fmt.Println(value.Email)
	}
}

func ListTrashTest(cm *Manager) {
	contacts := cm.ListTrash()
	for _, value := range contacts {
		fmt.Println(value.Email)
	}
}

func RestoreContactByEmailTest(cm *Manager) {
	if err := cm.RestoreContactByEmail("labibtest1@gmail.com"); err != nil {
		log.Println(err)
	}
}

func RemoveContactByEmailFromTrashTest(cm *Manager) {
	if err := cm.RemoveContactByEmailFromTrash("labibtest1@gmail.com"); err != nil {
		log.Println(err)
	}
}

func SortContactByEmailTest(cm *Manager) {
	contacts := cm.SortContactByEmail()
	for _, value := range contacts {
		fmt.Println(value.Email)
	}
}