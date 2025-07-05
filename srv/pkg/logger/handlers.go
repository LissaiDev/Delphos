package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileHandler escreve logs em arquivo
type FileHandler struct {
	file     *os.File
	filename string
	mutex    sync.Mutex
}

// NewFileHandler cria um novo handler de arquivo
func NewFileHandler(filename string) (*FileHandler, error) {
	// Cria diretório se não existir
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &FileHandler{
		file:     file,
		filename: filename,
	}, nil
}

// Handle escreve a mensagem no arquivo
func (h *FileHandler) Handle(level Level, message string, fields map[string]interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, err := fmt.Fprintln(h.file, message)
	return err
}

// Close fecha o arquivo
func (h *FileHandler) Close() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.file != nil {
		return h.file.Close()
	}
	return nil
}

// RotatingFileHandler escreve logs em arquivo com rotação
type RotatingFileHandler struct {
	baseFilename string
	maxSize      int64
	maxFiles     int
	currentFile  *os.File
	currentSize  int64
	mutex        sync.Mutex
}

// NewRotatingFileHandler cria um novo handler com rotação de arquivos
func NewRotatingFileHandler(baseFilename string, maxSize int64, maxFiles int) (*RotatingFileHandler, error) {
	handler := &RotatingFileHandler{
		baseFilename: baseFilename,
		maxSize:      maxSize,
		maxFiles:     maxFiles,
	}

	if err := handler.openFile(); err != nil {
		return nil, err
	}

	return handler, nil
}

// openFile abre o arquivo atual
func (h *RotatingFileHandler) openFile() error {
	file, err := os.OpenFile(h.baseFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	h.currentFile = file
	h.currentSize = info.Size()
	return nil
}

// rotate rotaciona os arquivos de log
func (h *RotatingFileHandler) rotate() error {
	if h.currentFile != nil {
		h.currentFile.Close()
	}

	// Move arquivos existentes
	for i := h.maxFiles - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d", h.baseFilename, i)
		newName := fmt.Sprintf("%s.%d", h.baseFilename, i+1)

		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	// Move arquivo atual
	if _, err := os.Stat(h.baseFilename); err == nil {
		os.Rename(h.baseFilename, h.baseFilename+".1")
	}

	return h.openFile()
}

// Handle escreve a mensagem no arquivo com rotação
func (h *RotatingFileHandler) Handle(level Level, message string, fields map[string]interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Verifica se precisa rotacionar
	if h.currentSize+h.currentSize > h.maxSize {
		if err := h.rotate(); err != nil {
			return err
		}
	}

	written, err := fmt.Fprintln(h.currentFile, message)
	if err != nil {
		return err
	}

	h.currentSize += int64(written)
	return nil
}

// Close fecha o arquivo atual
func (h *RotatingFileHandler) Close() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.currentFile != nil {
		return h.currentFile.Close()
	}
	return nil
}

// MultiHandler executa múltiplos handlers
type MultiHandler struct {
	handlers []Handler
}

// NewMultiHandler cria um novo handler múltiplo
func NewMultiHandler(handlers ...Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// Handle executa todos os handlers
func (h *MultiHandler) Handle(level Level, message string, fields map[string]interface{}) error {
	var lastErr error

	for _, handler := range h.handlers {
		if err := handler.Handle(level, message, fields); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// FilterHandler filtra logs baseado em critérios
type FilterHandler struct {
	handler Handler
	filter  func(level Level, message string, fields map[string]interface{}) bool
}

// NewFilterHandler cria um novo handler com filtro
func NewFilterHandler(handler Handler, filter func(level Level, message string, fields map[string]interface{}) bool) *FilterHandler {
	return &FilterHandler{
		handler: handler,
		filter:  filter,
	}
}

// Handle executa o handler apenas se passar no filtro
func (h *FilterHandler) Handle(level Level, message string, fields map[string]interface{}) error {
	if h.filter(level, message, fields) {
		return h.handler.Handle(level, message, fields)
	}
	return nil
}

// LevelFilterHandler filtra por nível mínimo
func LevelFilterHandler(handler Handler, minLevel Level) *FilterHandler {
	return NewFilterHandler(handler, func(level Level, message string, fields map[string]interface{}) bool {
		return level >= minLevel
	})
}
