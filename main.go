package main

func main() {
	if err := createEntityFile(CreateEntityInput{
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
	}); err != nil {
		panic(err)
	}
}
