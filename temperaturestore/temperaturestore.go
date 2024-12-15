package temperaturestore

import (
	"sync"
	"time"
)

type Temperature struct {
	Id       int       `json:"id"`
	Value    int       `json:"value"`
	City     string    `json:"city"`
	Datetime time.Time `json:"datetime"`
}

type TemperatureStore struct {
	sync.Mutex
	temperatures map[int]Temperature
	nextId       int
}

func New() *TemperatureStore

func (ts *TemperatureStore) CreateTemperature(value int, city string, datetime time.Time) int

func (ts *TemperatureStore) GetTemperature(id int) (Temperature, error)

func (ts *TemperatureStore) DeleteTemperature(id int) error

func (ts *TemperatureStore) DeleteAllTemperatures() error

func (ts *TemperatureStore) GetAllTemperatures() []Temperature

func (ts *TemperatureStore) GetTemperaturesByCity(city string) []Temperature

func (ts *TemperatureStore) GetTemperaturesByDatetime(year int, month time.Month, day int) []Temperature
