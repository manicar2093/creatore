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
	if err := createEntityFile(data); err != nil {
		panic(err)
	}
	if err := createRepositoryFile(CreateRepositoryInput{data}); err != nil {
		panic(err)
	}
}
