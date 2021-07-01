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
    "github.com/gorilla/mux"
)


type record struct {
    Key string `json:"key"`
    CreatedAt time.Time `json:"createdAt"`
    TotalCount int `json:"totalCount"`
}
type input struct {
    StartDate string `json:"startDate,omitempty"`
    EndDate string `json:"endDate,omitempty"`
    MinCount int `json:"minCount,omitempty"`
    MaxCount int `json:"maxCount,omitempty"`
}

type response struct {
    Code int `json:"code"`
    Msg string `json:"msg"`
    Records []record `json:"record,omitempty"`
}

type key_value struct {
    Key string `json:"key"`
    Value string `json:"value"`
}

var key_values []key_value

func handler_key_value(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","application/json")
/*    fmt.Println("here")*/
    params := r.URL.Query()
    var search_value string
    for _, str := range params["key"]{
        search_value=search_value+str
        if search_value == "active-tabs" {
            fmt.Println("here")
            break
        }
    }

    for _, key_value := range key_values {

        if key_value.Key == search_value {
            json.NewEncoder(w).Encode(key_value)
            return
        }
    }
}


func handler(w http.ResponseWriter, r *http.Request) {

    
    var requestData input
    switch r.Method {
    case "GET":
        w.Header().Set("Content-Type", "text/html")
        p , _:=ioutil.ReadFile("main.html")
        fmt.Fprintf(w, string(p))  
        /*http.ServeFile(w, r, "main.html")*/
        return
    case "POST":
        w.Header().Set("Content-Type","application/json")

       default:
        fmt.Println("Only POST or GET requests are allowed")
        fmt.Fprintf(w,"Only POST or GET requests are allowed")
    }

        body, _ := ioutil.ReadAll(r.Body)
        err :=json.Unmarshal(body, &requestData)
        if err != nil {
            fmt.Println(err)
        }


    if requestData==(input{}){

        var in_memory_input key_value
        _ =json.Unmarshal(body, &in_memory_input)
        key_values=append(key_values, in_memory_input)
        key_values_json, _:=json.Marshal(key_values)
        fmt.Fprintf(w, "%+v", string(key_values_json)) 
    } else {
        startDate :=requestData.StartDate
        startDateString := fmt.Sprintf("%v", startDate)
        endDate := requestData.EndDate
        endDateString := fmt.Sprintf("%v", endDate)
        minCount :=requestData.MinCount
        maxCount :=requestData.MaxCount


        client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"))
        if err != nil {
            log.Fatal(err)
        }
        ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
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



        var records []record
        if err = data.All(ctx, &records); err != nil {
            fmt.Println(err)
        }

        var result response
        if len(records)==0 {
            result.Code=-1
            result.Msg="No record between given dates..."
        } else {
            result.Code=0
            result.Msg="Success"
            result.Records=records        
        }

        result_json, _ :=json.Marshal(result)
        fmt.Fprintf(w, "%+v", string(result_json))  
    }

}



func main() {


    port, err := os.LookupEnv("PORT")
    if !err {
        port = "3000"
    } 

    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/in_memory", handler_key_value).Methods("GET")
    http.Handle("/", r)

    fmt.Println("Server is up on " + port)
    log.Fatal(http.ListenAndServe(":"+port,nil))
}