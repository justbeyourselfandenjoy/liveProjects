package main

import (
	"log"
)

func validateJsonAgainstSchema(schema []byte, instance []byte) error {

	log.Println("validateSchema has been called")

	/*
		compiler := jsonschema.NewCompiler()
		sch, err := compiler.Compile(string(schema))
		if err != nil {
			return err
		}

		log.Println(sch.Location)
	*/
	/*
		var m interface{}
		if err := yaml.Unmarshal(schema, &m); err != nil {
			return err
		}
	*/
	/*
		petSchema, err := jsonschema.UnmarshalJSON(strings.NewReader(string(schema)))
		if err != nil {
			return err
		}

			inst, err := jsonschema.UnmarshalJSON(strings.NewReader(string(instance)))
			if err != nil {
				return err
			}

			compiler := jsonschema.NewCompiler()

			if err := compiler.AddResource("swagger.json", petSchema); err != nil {
				return err
			}
			log.Println("compiler.AddResource has been passed")

			sch, err := compiler.Compile("swagger.json")
			if err != nil {
				return err
			}
			log.Println("compiler.Compile has been passed")

			if err = sch.Validate(inst); err != nil {
				return err
			}
	*/
	return nil
}
