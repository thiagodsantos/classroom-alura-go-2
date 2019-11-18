package main

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Produto struct {
	Nome       string
	Descricao  string
	Preco      float64
	Quantidade int
}

var temp = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	produtos := listProdutos()
	temp.ExecuteTemplate(w, "Index", produtos)
}

func insertProduto(produto Produto) interface{} {
	db, ctx := connectDatabase()
	collection := db.Database("loja").Collection("produtos")

	res, err := collection.InsertOne(ctx, produto)

	if err != nil {
		panic(err)
	}

	db.Disconnect(ctx)

	return res.InsertedID
}

func listProdutos() []Produto {
	db, ctx := connectDatabase()
	collection := db.Database("loja").Collection("produtos")

	filter := bson.D{{}}
	findOptions := options.Find()

	p := Produto{}
	produtos := []Produto{}

	findProdutos, err := collection.Find(ctx, filter, findOptions)

	if err != nil {
		panic(err)
	}

	for findProdutos.Next(ctx) {
		var produto Produto
		err := findProdutos.Decode(&produto)
		if err != nil {
			panic(err)
		}

		p.Nome = produto.Nome
		p.Descricao = produto.Descricao
		p.Preco = produto.Preco
		p.Quantidade = produto.Quantidade

		produtos = append(produtos, p)
	}

	return produtos
}

func connectDatabase() (*mongo.Client, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	fmt.Println("MongoDB Conectado")

	return client, ctx
}
