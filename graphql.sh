#!/bin/bash
# prep for schema

#!/bin/bash

# 1. Install gqlgen
go get github.com/99designs/gqlgen

# 2. Initialize gqlgen (creates gqlgen.yml and graph/ directory)
go run github.com/99designs/gqlgen init

# 3. Move your schema to the default location
mkdir -p graph
mv internal/api/schema.graphqls graph/schema.graphqls

# 4. Generate code from your schema
go run github.com/99designs/gqlgen generate

echo "gqlgen setup complete! Edit graph/resolver.go to implement your resolvers."