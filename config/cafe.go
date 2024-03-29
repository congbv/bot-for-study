package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgryski/dgohash"
)

//CafeConfig ok
type CafeConfig struct {
	FirstOrderTime      CafeTime      `json:"first_order_time"`
	LastOrderTime       CafeTime      `json:"last_order_time"`
	TimeLocation        *CafeLocation `json:"time_location"`
	TimeSlotIntervalMin int           `json:"time_slot_interval_min"`

	// OrderChan is a public channel name in a form of @channelname
	// or a private channel id. See how to get private channel id:
	// https://github.com/python-telegram-bot/python-telegram-bot/issues/538
	OrderChan string `json:"orders_channel"`
	Menu      Menu   `json:"menu"`
}

// CafeTime exists only because we need to unmarshal string of type HH:MM into
// time.Time. It is possible only by having custom named type with its own
// json.Unmarshaler implementation
type CafeTime time.Time

//UnmarshalJSON ok
func (c *CafeTime) UnmarshalJSON(data []byte) error {
	dlen := len(data)
	if dlen < 2 {
		return errors.New("invalid time format")
	}

	t, err := time.Parse("15:04", string(data[1:dlen-1]))
	if err != nil {
		return err
	}

	ct := CafeTime(t)
	*c = ct

	return nil
}

//CafeLocation ok
type CafeLocation time.Location

//UnmarshalJSON ok
func (l *CafeLocation) UnmarshalJSON(data []byte) error {
	dlen := len(data)
	if dlen < 2 {
		return errors.New("invalid time location format")
	}

	loc, err := time.LoadLocation(string(data[1 : dlen-1]))
	if err != nil {
		return err
	}

	*l = CafeLocation(*loc)

	return nil
}

//Menu ok
type Menu struct {
	Map        map[string][]Meal
	Categories []string
}

//UnmarshalJSON ok
func (m *Menu) UnmarshalJSON(data []byte) error {
	dlen := len(data)
	if dlen < 2 {
		return errors.New("invalid menu data")
	}

	if err := json.Unmarshal(data, &m.Map); err != nil {
		return err
	}

	for cat := range m.Map {
		m.Categories = append(m.Categories, cat)
	}
	return nil
}

//Meal ok
type Meal struct {
	Val  string
	Hash string
}

//UnmarshalJSON ok
func (m *Meal) UnmarshalJSON(data []byte) error {
	dlen := len(data)
	if dlen < 2 {
		return errors.New("invalid meal data")
	}

	m.Val = string(data[1 : dlen-1])
	m.Hash = hashMeal(m.Val)

	return nil
}

//MealByHash ok
func (m *Menu) MealByHash(hash string) (string, bool) {
	if m == nil {
		return "", false
	}
	for _, cat := range m.Map {
		for _, meal := range cat {
			if meal.Hash == hash {
				return meal.Val, true
			}
		}
	}
	return "", false
}

func loadCafeConfig(f string) (conf CafeConfig, err error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &conf)
	return
}

func hashMeal(meal string) string {
	hash := dgohash.NewSDBM32()
	hash.Write([]byte(meal))
	return fmt.Sprintf("%d", hash.Sum32())
}
