package tables

import (
	"fmt"
	"jhgambling/backend/core/utils"
	"sync"
)

// TableRegistry manages all registered tables
type TableRegistry struct {
	tables map[string]Table
	mutex  sync.RWMutex
}

// NewTableRegistry creates a new registry
func NewTableRegistry() *TableRegistry {
	return &TableRegistry{
		tables: make(map[string]Table),
	}
}

// Register adds a table to the registry
func (r *TableRegistry) Register(table Table) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tableID := table.GetID()
	if tableID == "" {
		return fmt.Errorf("table ID cannot be empty")
	}

	if _, exists := r.tables[tableID]; exists {
		return fmt.Errorf("table with ID '%s' is already registered", tableID)
	}

	r.tables[tableID] = table
	utils.Log("ok", "casino::data::tables", fmt.Sprintf("registered table: %s", tableID))
	return nil
}

// Get retrieves a registered table by ID
func (r *TableRegistry) Get(tableID string) (Table, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	table, exists := r.tables[tableID]
	if !exists {
		return nil, fmt.Errorf("table with ID '%s' not found", tableID)
	}

	return table, nil
}

// GetAll returns all registered tables
func (r *TableRegistry) GetAll() []Table {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tables := make([]Table, 0, len(r.tables))
	for _, table := range r.tables {
		tables = append(tables, table)
	}

	return tables
}

// Remove unregisters a table by ID
func (r *TableRegistry) Remove(tableID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tables[tableID]; !exists {
		return fmt.Errorf("table with ID '%s' not found", tableID)
	}

	delete(r.tables, tableID)
	utils.Log("ok", "casino::data::tables", fmt.Sprintf("unregistered table: %s", tableID))
	return nil
}
