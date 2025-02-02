package main

func main() {
	data := CreateEntityInput{
		EntityName: "User",
		IsUuid:     true,
		Fields: []EntityField{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name:       "LastNames",
				Type:       "string",
				IsOptional: true,
			},
			{
				Name:       "CreatedAt",
				Type:       "time",
				IsOptional: false,
			},
		},
	}
	if err := trigger(data); err != nil {
		panic(err)
	}
}

func trigger(input CreateEntityInput) error {
	data := createUsefulData(input)
	if err := createEntityFile(data); err != nil {
		return err
	}
	if err := createRepositoryFile(data); err != nil {
		return err
	}
	if err := createControllerFile(data); err != nil {
		return err
	}
	return nil
}
