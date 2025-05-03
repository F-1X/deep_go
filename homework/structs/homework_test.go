package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		var byteName [42]byte
		for i, b := range name {
			byteName[i] = byte(b)
		}
		person.name = byteName
	}
}


func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		leftMana := byte(mana & int(0b1111111100) >> 2)
		rightMana := byte(mana & int(0b11))
		
		person.manaHealth[0] = leftMana
		person.manaHealth[2] ^= (rightMana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		leftHealth := byte(health & int(0b1111111100) >> 2)
		rightHealth := byte(health & int(0b11))

		person.manaHealth[1] = leftHealth
		person.manaHealth[2] ^= (rightHealth << 2)
	}
}
func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectStrengthExpLvl ^= uint16((uint16(respect) << 12))
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectStrengthExpLvl ^= uint16((uint16(strength) << 8))
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectStrengthExpLvl ^= uint16((uint16(experience) << 4))
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectStrengthExpLvl ^= uint16((uint16(level)))
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.HouseGunFamilyRole ^= 0b100
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.HouseGunFamilyRole ^= 0b1000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.HouseGunFamilyRole ^= 0b10000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.HouseGunFamilyRole ^= byte(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)
type GamePerson struct {
	x                     int32
	y                     int32
	z                     int32
	gold                  uint32
	RespectStrengthExpLvl uint16
	name                  [42]byte
	manaHealth            [3]byte
	HouseGunFamilyRole    byte
}

func NewGamePerson(options ...Option) GamePerson {
	g := GamePerson{}
	for _, opt := range options {
		opt(&g)
	}
	return g
}


func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}
func (p *GamePerson) Mana() int {
	thirdByte := int(p.manaHealth[2])
	fmana := int(p.manaHealth[0]) << 2
	return int((fmana) ^ (thirdByte & 0b11))
}
func (p *GamePerson) Health() int {
	thirdByte := int(p.manaHealth[2]) >> 2
	fhealth := int(p.manaHealth[1]) << 2
	return int((fhealth) ^ (thirdByte & 0b11))
}

func (p *GamePerson) Respect() int {
	return int(p.RespectStrengthExpLvl&0b1111000000000000) >> 12
}

func (p *GamePerson) Strength() int {
	return int(p.RespectStrengthExpLvl&0b0000111100000000) >> 8
}

func (p *GamePerson) Experience() int {
	return int(p.RespectStrengthExpLvl&0b0000000011110000) >> 4
}


func (p *GamePerson) Level() int {
	return int(p.RespectStrengthExpLvl & 0b0000000000001111)
}

func (p *GamePerson) HasHouse() bool {
	return (p.HouseGunFamilyRole>>2)&1 == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.HouseGunFamilyRole>>3)&1 == 1
}

func (p *GamePerson) HasFamilty() bool {
	return (p.HouseGunFamilyRole>>4)&1 == 1
}

func (p *GamePerson) Type() int {
	return int(p.HouseGunFamilyRole & 0b11)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
