/*
Copyright (C) 2021 sleal (Stuart Leal Quesada)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package logs

import (
	"fmt"
	"log"
	"os"
)

type logHandler struct {
	Info             *log.Logger
	Debug            *log.Logger
	Warn             *log.Logger
	Error            *log.Logger
	activeSeverities map[string]bool
}

// logHandler singleton makes sure all loggers use the same Loggers
var sLogHandler *logHandler

// init initializes the sLogHandler singleton, which is used by all loggers
func initLogHandler() {

	// By default, the activeSeverities are set to enable all except debug
	activeSeverities := map[string]bool{
		"Info":  true,
		"Debug": false,
		"Warn":  true,
		"Error": true,
	}

	sLogHandler = &logHandler{activeSeverities: activeSeverities}

	file, err := os.OpenFile(
		LOG_FILENAME,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	sLogHandler.Info =
		log.New(file, INFO_HEADER, log.Ldate|log.Ltime|log.Lshortfile)
	sLogHandler.Debug =
		log.New(file, DEBUG_HEADER, log.Ldate|log.Ltime|log.Lshortfile)
	sLogHandler.Warn =
		log.New(file, WARNING_HEADER, log.Ldate|log.Ltime|log.Lshortfile)
	sLogHandler.Error =
		log.New(file, ERROR_HEADER, log.Ldate|log.Ltime|log.Lshortfile)
}

// logger has the complete information in order to log information to the
// log file. If no sLogHandler singleton is not initialized, then creating
// logs is meant to fail
type logger struct {
	logMsg     string
	logHandler *logHandler
}

// newLog creates a new logger given the message provided
func NewLog(msg string, args ...interface{}) logger {

	// Check that sLogHandler has been initialized, if not do so.
	if sLogHandler == nil {
		initLogHandler()
	}

	// Create the new logger and return it
	return logger{
		logMsg:     fmt.Sprintf(msg, args...),
		logHandler: sLogHandler,
	}
}

// isSeverityActivated returns true if the given severity is activated, false
// otherwise
func (l *logger) isSeverityActivated(tag string) bool {

	if activated, ok := l.logHandler.activeSeverities[tag]; ok {
		return activated
	}

	return false
}

// Info logs an info message using the sLogHandler which should be a
// singleton
func (l *logger) Info() {
	if !l.isSeverityActivated("Info") {
		return
	}

	l.logHandler.Info.Println(l.logMsg)
}

// Debug logs a debug message using the sLogHandler which should be a
// singleton
func (l *logger) Debug() {
	if !l.isSeverityActivated("Debug") {
		return
	}

	l.logHandler.Debug.Println(l.logMsg)
}

// Warn logs a warning message using the sLogHandler which should be a
// singleton
func (l *logger) Warn() {
	if !l.isSeverityActivated("Warn") {
		return
	}

	l.logHandler.Warn.Println(l.logMsg)
}

// Error logs an error message using the sLogHandler which should be a
// singleton
func (l *logger) Error() {
	if !l.isSeverityActivated("Error") {
		return
	}

	l.logHandler.Error.Println(l.logMsg)
}

// Msg returns the logMsg string
func (l *logger) Msg() string { return l.logMsg }
