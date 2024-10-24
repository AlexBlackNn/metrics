package main

import (
	"fmt"
	"log"

	"github.com/AlexBlackNn/metrics/app/agent/encryption"
)

func main() {
	// Путь к файлам с ключами
	publicKeyPath := "/home/alex/Dev/GolandYandex/metrics/public_key.pem"   // Укажите путь к вашему публичному ключу
	privateKeyPath := "/home/alex/Dev/GolandYandex/metrics/private_key.pem" // Укажите путь к вашему закрытому ключу

	// Создание инкриптора
	encryptor, err := encryption.NewEncryptor(publicKeyPath)
	if err != nil {
		log.Fatalf("Ошибка при создании инкриптора: %v", err)
	}

	// Тестовое сообщение для шифрования
	testMessage := `{"id":"1", "type":"counter", "value": 10}`
	fmt.Println("Исходное сообщение:", testMessage)

	// Шифрование сообщения
	encryptedMessage, err := encryptor.EncryptMessage(testMessage)
	if err != nil {
		log.Fatalf("Ошибка при шифровании сообщения: %v", err)
	}
	fmt.Println("Зашифрованное сообщение:", encryptedMessage)

	// Создание декриптора
	decryptor, err := encryption.NewDecryptor(privateKeyPath)
	if err != nil {
		log.Fatalf("Ошибка при создании декриптора: %v", err)
	}

	// Расшифровка сообщения
	decryptedMessage, err := decryptor.DecryptMessage(encryptedMessage)
	if err != nil {
		log.Fatalf("Ошибка при расшифровке сообщения: %v", err)
	}
	fmt.Println("Расшифрованное сообщение:", decryptedMessage)
}
