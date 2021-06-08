package main

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)



func handler(w http.ResponseWriter, r *http.Request) {

    var t map[string]interface{}
    switch r.Method {
    case "GET":
        http.ServeFile(w, r, "main.html")
        return
    case "POST":
        requestData, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Fatal(err)
        }

        requestString := string(requestData)

        json.Unmarshal([]byte(requestString),&t)
    
       default:
        fmt.Println("Only POST or GET requests are allowed")
        fmt.Fprintf(w,"Only POST or GET requests are allowed")
    }


    startDate :=t["startDate"]
    startDateString := fmt.Sprintf("%v", startDate)
    endDate := t["endDate"]
    endDateString := fmt.Sprintf("%v", endDate)
    minCount :=t["minCount"]
    maxCount :=t["maxCount"]

    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }



    defer client.Disconnect(ctx)    


    collection := client.Database("getir-case-study").Collection("records")
    



fmt.Println("Succesfully connected to the database...")
tm1, err := time.Parse("2006-01-02", endDateString)
tm2, err := time.Parse("2006-01-02", startDateString)
o1 := bson.D{{"$project", bson.D{{"_id", 0}, {"key", 1}, {"createdAt", 1}, {"totalCount", bson.D{{"$reduce", bson.D{{"input", "$counts"},{"initialValue",0},{"in", bson.D{{"$add", bson.A{"$$value","$$this"}}}}}}}}} }}
o2 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$gt", tm2}}}}}}
o3 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$lt", tm1}}}}}}
o4 := bson.D{{"$match", bson.D{{"totalCount", bson.D{{"$lt", maxCount}, {"$gt", minCount}}}}}}
pipeline :=mongo.Pipeline{o1,o2,o3, o4}

data , err := collection.Aggregate(ctx, pipeline)
if err != nil {
    panic(err)
}
var showsWithInfo []bson.M
if err = data.All(ctx, &showsWithInfo); err != nil {
    panic(err)
}

fmt.Fprintf(w, "%+v", showsWithInfo)  




}



func main() {


    port, err := os.LookupEnv("PORT")
    // fmt.Println(err)
    if !err {
        port = "3000"
    } 


    http.HandleFunc("/",handler)
    fmt.Println("Server is up on " + port)
    log.Fatal(http.ListenAndServe(":"+port,nil))
}