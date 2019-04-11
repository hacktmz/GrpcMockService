package main

import (
	"fmt"
	//proto "github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	//"github.com/jhump/protoreflect/dynamic"
	"encoding/json"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/hacktmz/GrpcMockService/schema"
)

const (
	file = "./test/hello.proto"
)

func Proto2Json(protoFileBytes []byte) ([]byte, error) {
	err := schema.LoadFileProto(file, protoFileBytes)
	if err != nil {
		fmt.Errorf("Failed to open proto: %v", err)
		return nil, err
	}

	fd, err := schema.LoadFileDescriptor(file)
	if err != nil {
		fmt.Errorf("Not Found LoadFileDescriptor %s", file)
		return nil, err
	}
	json_map := make(map[string]interface{})
	sds := fd.GetServices()
	for _, sd := range sds {
		master := make(map[string]interface{})
		fmt.Printf("000===============(%s) \n", sd.GetFullyQualifiedName())
		mds := sd.GetMethods()
		fmt.Printf("111===============(%s) \n", sd.GetName())
		for _, methodd := range mds {
			method_name := methodd.GetName()
			msgd := methodd.GetOutputType()
			fieldsd := msgd.GetFields()
			outname := msgd.GetFullyQualifiedName()
			fmt.Printf("\n 333 ===============(%s) \n", outname)
			m_master := phrase(fieldsd)
			master[method_name] = m_master
			fmt.Printf("===========ALL 2===== %v \n", m_master)
		}
		json_map[sd.GetFullyQualifiedName()] = master
	}
	fmt.Printf("===========ALL===== %v, ===\n", json_map)

	body, err := json.Marshal(json_map)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("===========JSON===== %s\n", body)
	return body, nil
}

func phrase(fieldsd []*desc.FieldDescriptor) map[string]interface{} {
	m_master := make(map[string]interface{})
	for _, field := range fieldsd {
		fmt.Printf("444===============(%s) %v\n ", field, field.GetDefaultValue())
		if field.IsMap() {
			m_level1 := make(map[string]interface{})
			fieldd_key := field.GetMapKeyType()
			fieldd_val := field.GetMapValueType()
			fmt.Printf("5551===============(%s) ==== %s\n ", fieldd_key.GetJSONName(), fieldd_val.GetJSONName())
			if dpb.FieldDescriptorProto_TYPE_ENUM == fieldd_val.GetType() { //is map &is enum
				enumd := fieldd_val.GetEnumType()
				enum_valsd := enumd.GetValues()
				m_temp_enum := make(map[string]interface{})
				vals := make([]string, 0)
				for _, val := range enum_valsd {
					fmt.Printf("5552222===============(%s) \n ", val.GetName())
					vals = append(vals, val.GetName())
				}
				m_temp_enum[fieldd_val.GetType().String()] = vals
				m_level1[fieldd_key.GetType().String()] = m_temp_enum
				fmt.Printf("5552222333===============(%v) \n ", m_level1)

			} else if dpb.FieldDescriptorProto_TYPE_MESSAGE == fieldd_val.GetType() {
				msgd_temp := fieldd_val.GetMessageType()
				fieldsd_temp := msgd_temp.GetFields()
				m_level2 := make(map[string]interface{})
				m_level2 = phrase(fieldsd_temp)
				m_level1[fieldd_key.GetType().String()] = m_level2
			} else {
				fmt.Printf("666====== %s , === %s\n", fieldd_val.GetJSONName(), fieldd_val.GetType())
				m_level1[fieldd_key.GetType().String()] = fieldd_val.GetType().String()
			}
			m_master[field.GetJSONName()] = m_level1
		} else {
			fmt.Printf("555===============(%s) ====111 === %s \n ", field.GetType().String(), field.GetJSONName())
			if dpb.FieldDescriptorProto_TYPE_ENUM == field.GetType() { //not map &is enum
				enumd := field.GetEnumType()
				enum_valsd := enumd.GetValues()
				m_temp_enum := make(map[string]interface{})
				vals := make([]string, 0)
				for _, val := range enum_valsd {
					fmt.Printf("55566666===============(%s) \n ", val.GetName())
					vals = append(vals, val.GetName())
				}
				m_temp_enum[field.GetType().String()] = vals

				if field.IsRepeated() {
					enum_map_slice := make([]interface{}, 0)
					enum_map_slice = append(enum_map_slice, m_temp_enum)
					m_master[field.GetJSONName()] = enum_map_slice
				} else {
					m_master[field.GetJSONName()] = m_temp_enum
				}

			} else if dpb.FieldDescriptorProto_TYPE_MESSAGE == field.GetType() { //not map && is msg
				msgd_temp := field.GetMessageType()
				fieldsd_temp := msgd_temp.GetFields()
				m_level1 := make(map[string]interface{})
				m_level1 = phrase(fieldsd_temp)

				if field.IsRepeated() {
					fmt.Printf(" 888888881 %v \n", field.GetJSONName())
					val := make([]interface{}, 0)
					val = append(val, m_level1)
					m_master[field.GetJSONName()] = val
				} else {
					m_master[field.GetJSONName()] = m_level1
					fmt.Printf("===========ALL 1===== %v, ===\n", m_level1)
				}
			} else { //not map && not msg && not enum
				if field.IsRepeated() {
					fmt.Printf(" 888888882 %v \n", field.GetJSONName())
					val := make([]interface{}, 0)
					val = append(val, field.GetType().String())
					m_master[field.GetJSONName()] = val
				} else {
					m_master[field.GetJSONName()] = field.GetType().String()
				}
			}

		}
	}
	return m_master
}
