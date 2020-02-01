/*
 * Copyright (c) 2020.
 * Developer vsb
 * License Apache 2.0
 */

package mim

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	bufferSize = 1024 * 1024
)

var newLine = []byte("\n")
var seperator = []byte("(+)")
var errorString = []byte("Error")
var infoString = []byte("Info")

type MIM struct {
}

type logger struct {
	file *os.File
	lock sync.Mutex
}

func (l *logger) log(direction []byte, data []byte) {
	l.lock.Lock()
	defer l.lock.Unlock()
	write := func(d []byte) {
		_, err := l.file.Write(d)
		if err != nil {
			fmt.Println("Logging error:", err.Error())
			fmt.Println("Exiting application.")
			os.Exit(1)
		}
	}
	write([]byte(time.Now().String()))
	write(seperator)
	write(direction)
	write(seperator)
	write(data)
	write(newLine)
	if err := l.file.Sync(); err != nil {
		fmt.Println("Syncing error:", err.Error())
		fmt.Println("Exiting application.")
		os.Exit(1)
	}
}

func (l *logger) logInfo(data []byte) {
	l.log(infoString, data)
}

func (l *logger) logError(data []byte) {
	l.log(errorString, data)
}

func (m *MIM) Start(port string, destHostPort string) {

	logger := &logger{}
	file, err := os.Create("mim.dat")
	if err != nil {
		fmt.Println("Unable to open mim.dat. error:", err.Error())
		return
	}
	logger.file = file
	logger.logInfo([]byte("Starting listener on port: " + port))
	l, err := net.Listen("tcp", port)
	if err != nil {
		logger.logError([]byte("Error starting listener on port: " + port + " error: " + err.Error()))
		return
	}

	for con, err := l.Accept(); err == nil; {
		go handleAccept(logger, con, destHostPort)
	}
}

func handleAccept(logger *logger, sourceCon net.Conn, destHostPort string) {
	logger.logInfo([]byte("Connecting to destination: " + destHostPort))
	destCon, err := net.Dial("tcp", destHostPort)
	if err != nil {
		logger.logError([]byte("Error connecting " + destHostPort + " error: " + err.Error()))
		return
	}
	logger.logInfo([]byte("Connection established to destination: " + destHostPort))
	handler := func(quit chan string, source, dest net.Conn, direction string) {
		dirBytes := []byte(direction)
		buffer := make([]byte, bufferSize)
		for {
			n, err := source.Read(buffer)
			if err != nil {
				quit <- direction
				return
			}
			logger.log(dirBytes, buffer[0:n])
			_, err = dest.Write(buffer[0:n])
			if err != nil {
				quit <- direction
				return
			}
		}
	}

	quit := make(chan string, 2)

	go handler(quit, sourceCon, destCon, "S->D")
	go handler(quit, destCon, sourceCon, "D->S")

	for count := 0; count < 2; {
		logger.logInfo([]byte("Waiting for quite signal"))
		logger.logInfo([]byte("Got quite signal for " + <-quit))
		count++
	}
}
