package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
)

// Кодирует строку по алгоритму AES
func Encrypt(v string, k string) (string, error) {
	// Приводим параметры к срезу байт
	value := []byte(v)
	key := []byte(k)

	// Длина ключа должена быть кратной размеру блока шифра (aes.BlockSize)
	// В данном случае 16-ти
	// Поэтому просто берем md5 сумму от ключа и получим срез из 32 элементов
	md5 := md5.Sum(key) // md5.Sum возвращает массив, а не срез
	key = md5[:]        // Создаем срез

	// Длина строки как и ключ должна быть кратной 16
	// Так что используем формат PKCS7
	// Просто добавляем недостающие байты
	// Смотри функцию pad ниже!
	value = pad(value)

	// Входной вектор для шифрования
	// Просто срез из случайных байт длинной блока шифра (aes.BlockSize)
	// Смотри функцию random ниже!
	iv := random(aes.BlockSize)

	// Конструктор aes
	block, err := aes.NewCipher(key)

	if err != nil {
		return "", err
	}

	// Создаем срез в который будет помещена наша закодированная строка
	ciphertext := make([]byte, len(value))

	// Магия
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, value)

	// Объединяем входной вектор и кодированную строку
	buf := bytes.NewBuffer(iv)
	buf.Write(ciphertext)
	result := buf.Bytes()

	// Возвращаем строку в представлении base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Декодирует строку, закодированную по алгоритму AES
func Decrypt(v string, k string) (string, error) {
	// Декодируем строку представленную в формате base64
	// Метод вернет срез байт
	value, err := base64.StdEncoding.DecodeString(v)

	if err != nil {
		return "", err
	}

	// Приводим параметры к срезу байт
	key := []byte(k)

	// Длина ключа должена быть кратной размеру блока шифра (aes.BlockSize)
	// В данном случае 16-ти
	// Поэтому просто возьмум md5 сумму от ключа и получим срез и 32 элементов
	md5 := md5.Sum(key) // md5.Sum возвращает массив, а не срез
	key = md5[:]        // Создаем срез

	// Извлекаем входной вектор для шифрования
	iv := value[:aes.BlockSize]

	// Извлекаем кодированный текст
	ciphertext := value[aes.BlockSize:]

	// Конструктор aes
	block, err := aes.NewCipher(key)

	if err != nil {
		return "", err
	}

	// Создаем срез в который будет помещена наша декодированная строка
	text := make([]byte, len(ciphertext))

	// Магия
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(text, ciphertext)

	// Получаем исходную строку
	return string(unpad(text)), nil
}

// Приводим к формату PKCS7
// Добавляем столько байт, сколько не хватает
func pad(value []byte) []byte {
	pdd := aes.BlockSize - (len(value) % aes.BlockSize)
	buf := bytes.NewBuffer(value)
	buf.Write(bytes.Repeat([]byte{byte(pdd)}, pdd))
	return buf.Bytes()
}

// Получаем исходное значение из строки в формате PKCS7
func unpad(value []byte) []byte {
	length := len(value)
	pdd := value[length-1:]
	before := length - int(pdd[0])

	if bytes.Equal(value[before:], bytes.Repeat(pdd, int(pdd[0]))) {
		return value[:before]
	}

	return value
}

// Генерирует случайную последостальность байт длинно size
func random(size int) []byte {
	r := make([]byte, size)
	rand.Read(r)
	return r
}
