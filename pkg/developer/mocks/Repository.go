// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "wall/pkg/entity"
import mock "github.com/stretchr/testify/mock"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Find provides a mock function with given fields: name
func (_m *Repository) Find(name string) (*entity.Developer, error) {
	ret := _m.Called(name)

	var r0 *entity.Developer
	if rf, ok := ret.Get(0).(func(string) *entity.Developer); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Developer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
