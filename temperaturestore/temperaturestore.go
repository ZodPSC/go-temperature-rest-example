package temperaturestore

import (
	"fmt"
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

func New() *TemperatureStore {
	ts := &TemperatureStore{}
	ts.temperatures = make(map[int]Temperature)
	ts.nextId = 0
	return ts
}

func (ts *TemperatureStore) CreateTemperature(value int, city string, datetime time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	temperature := Temperature{
		Id:       ts.nextId,
		Value:    value,
		City:     city,
		Datetime: datetime,
	}
	ts.temperatures[ts.nextId] = temperature
	ts.nextId++
	return temperature.Id
}

func (ts *TemperatureStore) GetTemperature(id int) (Temperature, error) {
	ts.Lock()
	defer ts.Unlock()

	temperature, ok := ts.temperatures[id]
	if ok {
		return temperature, nil
	} else {
		return Temperature{}, fmt.Errorf("temperature with id=%id not found", id)
	}
}

func (ts *TemperatureStore) DeleteTemperature(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.temperatures[id]; !ok {
		return fmt.Errorf("temperature with id=%d not found")
	}

	delete(ts.temperatures, id)
	return nil
}

func (ts *TemperatureStore) DeleteAllTemperatures() error {
	ts.Lock()
	defer ts.Unlock()

	ts.temperatures = make(map[int]Temperature)
	return nil
}

func (ts *TemperatureStore) GetAllTemperatures() []Temperature {
	ts.Lock()
	defer ts.Unlock()

	allTemperatures := make([]Temperature, 0, len(ts.temperatures))
	for _, temperature := range ts.temperatures {
		allTemperatures = append(allTemperatures, temperature)
	}
	return allTemperatures
}

func (ts *TemperatureStore) GetTemperaturesByCity(city string) []Temperature {
	ts.Lock()
	defer ts.Unlock()

	var temperatures []Temperature

	for _, temperature := range ts.temperatures {
		if temperature.City == city {
			temperatures = append(temperatures, temperature)
		}
	}
	return temperatures
}

func (ts *TemperatureStore) GetTemperaturesByDatetime(year int, month time.Month, day int) []Temperature {
	ts.Lock()
	defer ts.Unlock()

	var temperatures []Temperature

	for _, temperature := range ts.temperatures {
		y, m, d := temperature.Datetime.Date()
		if y == year && m == month && d == day {
			temperatures = append(temperatures, temperature)
		}
	}
	return temperatures
}
