package storage

import (
	"errors"
	"reflect"
)

type PathBinder interface {
	BindPath(p *Path)
}

func Load[T any](path *Path, target *T) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer to struct")
	}

	// for slice/maps just bind the path and return
	if binder, ok := interface{}(target).(PathBinder); ok {
		binder.BindPath(path)
		return nil
	}

	elem := val.Elem()
	typ := elem.Type()

	// only if we write custom accessor for the whole struct
	if accessor, ok := getAccessor(typ); ok {
		v, err := accessor.ReadAt(path.slot, path.info.Offset, path.st)
		if err != nil {
			return err
		}
		elem.Set(reflect.ValueOf(v))
		return nil
	}

	// struct fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldPath := path.Field(field.Name)
		if fieldPath.Error() != nil {
			return fieldPath.Error()
		}

		fieldValue := elem.Field(i)
		if fieldValue.CanAddr() { // if the value is addressable
			// for slice/maps just bind the path and continue
			if binder, ok := fieldValue.Addr().Interface().(PathBinder); ok {
				binder.BindPath(fieldPath)
				continue
			}
		}
		if accessor, ok := getAccessor(field.Type); ok {
			v, err := accessor.ReadAt(fieldPath.slot, fieldPath.info.Offset, fieldPath.st)
			if err != nil {
				return err
			}
			if v != nil {
				fieldValue.Set(reflect.ValueOf(v))
			}
			continue
		} else if field.Type.Kind() == reflect.Struct {
			if err := loadRecurse(fieldPath, elem.Field(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadRecurse(p *Path, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldPath := p.Field(field.Name)
		fieldValue := val.Field(i)
		if fieldValue.CanAddr() { // if the value is addressable
			// for slice/maps just bind the path and continue
			if binder, ok := fieldValue.Addr().Interface().(PathBinder); ok {
				binder.BindPath(fieldPath)
				continue
			}
		}
		if accessor, ok := getAccessor(field.Type); ok {
			v, err := accessor.ReadAt(fieldPath.slot, fieldPath.info.Offset, fieldPath.st)
			if err != nil {
				return err
			}
			if v != nil {
				val.Field(i).Set(reflect.ValueOf(v))
			}
		} else if field.Type.Kind() == reflect.Struct {
			if err := loadRecurse(fieldPath, val.Field(i)); err != nil {
				return err
			}
		}
		// skip dynamic types
	}
	return nil
}

func Save[T any](path *Path, source T) error {
	if _, ok := getAccessor(reflect.TypeOf(source)); ok {
		return Set(path, source)
	}

	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if de-referenced type has  accessor
	if _, ok := getAccessor(val.Type()); ok {
		return Set(path, source)
	}

	if val.Kind() != reflect.Struct {
		return errors.New("source must be a struct")
	}
	return saveRecurse(path, val)
}

func saveRecurse(p *Path, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldPath := p.Field(field.Name)
		fieldVal := val.Field(i)

		if accessor, ok := getAccessor(field.Type); ok {
			if err := accessor.WriteAt(fieldPath.slot, fieldPath.info.Offset, fieldVal.Interface(), fieldPath.st); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Struct {
			if err := saveRecurse(fieldPath, fieldVal); err != nil {
				return err
			}
		}
	}
	return nil
}
