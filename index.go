package main

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "context"
    //"strconv"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)


// type reqPayload struct {
//     startDate json.RawMessage
//     endDate string
//     minCount int
//     maxCount int
// }

// func handler(w http.ResponseWriter, r *http.Request) {

//     switch r.Method {
//     case "POST":
//         body, err := ioutil.ReadAll(r.Body)
//         if err != nil {
//             panic(err)
//         }
//         var t reqPayload
//         dec :=json.NewDecoder(strings.NewReader(string(body)))
//         // dec :=string(r.Body)
//         dec.Decode(&t)
//         // err = json.Unmarshal(body, &t)
//         if err !=nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }
//         fmt.Fprintf(w,"Request: %+v",string(body))
//         // log.Println(t)
//         // fmt.Fprintf(w, "Request: %+v", t)       
//        default:
//         fmt.Println("Only POST requests are allowed")
//         fmt.Fprintf(w,"Only POST requests are allowed")


//     }


// }

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
        // var t reqPayload
        // out, _ :=json.Marshal(requestString)
        json.Unmarshal([]byte(requestString),&t)
        // fmt.Fprintf(w,"Request: %+v",t["startDate"])
        // log.Println(t)
        // fmt.Fprintf(w, "Request: %+v", t)       
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
    // client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir- case-study?retryWrites=true"))
    // if err != nil {
    //     log.Fatal(err)
    // }
    // ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    // err = client.Connect(ctx)
    // if err != nil {
    //     log.Fatal(err)
    // }


    collection := client.Database("getir-case-study").Collection("records")
    // o1 := bson.M{{"$project": bson.M{{
    //             _id: 0,
    //             key: 1,
    //             createdAt: 1,
    //             totalCount: bson.M{
    //                 "$reduce": bson.M{
    //                     input: "$counts",
    //                     initialValue: 0,
    //                     in: bson.M{ "$add": bson.M["$$value" , "$$this"]}
    //                 }
    //             }
    //         }}
    // }}

    // o2 := bson.M{ "$match": bson.M{totalCount: bson.M{"$gt": minCount, "$lt": maxCount}}}

    // o3 :=bson.M{{ "$match", bson.D{{"createdAt", bson.E{{"$gt", startDate,"$lt", endDate}}}}}}



    // o2 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$lt", endDate}}}}}}

    // o3 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$gt", startDate}}}}}}




// operations := []bson.M{o3}

fmt.Println("Succesfully connected to the database...")


// databases, err := client.ListDatabaseNames(ctx, bson.M{})
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(databases)


// pipe := collection.mongo.Pipeline([]bson.M{{o3}})

// // Run the queries and capture the results
// results := []bson.M{}
// err1 := pipe.All(&results)

// if err1 != nil {
//     fmt.Printf("ERROR : %s\n", err1.Error())
//     return
// }

// fmt.Printf(results)

// fmt.Println("Succesfully connected to the database...")


// var result []bson.M
// tm := time.Date(int(endDateString[0:4]),time.Month(endDateString[5:7]), int(endDateString[8:10]), 0, 0, 0, 0, Z)
tm1, err := time.Parse("2006-01-02", endDateString)
// tm := tm1.UTC().UnixNano() / int64(time.Millisecond)
tm2, err := time.Parse("2006-01-02", startDateString)
// tm3 := tm2.UTC().UnixNano() / int64(time.Millisecond)
//a1 , _ :=strconv.Atoi(endDateString[0:4])
//tm := 3.154*1010*a1
// fmt.Println(tm)
// endDate1 = time.Date(endDateString[0:4],endDateString[5:7], endDateString[8:10])
// fmt.Println(int(endDate))
// addFieldsStage := bson.D{
//     {"$addFields", bson.D{
//         {"newdate", bson.D{
//             {"$dateFromString", bson.D{
//                 {"dateString", "$createdAt"}, 
//                 {"format", "%Y-%m-%d"},
//             }},
//         }},
//     }},
// }
o1 := bson.D{{"$project", bson.D{{"_id", 0}, {"key", 1}, {"createdAt", 1}, {"totalCount", bson.D{{"$reduce", bson.D{{"input", "$counts"},{"initialValue",0},{"in", bson.D{{"$add", bson.A{"$$value","$$this"}}}}}}}}} }}
o2 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$gt", tm2}}}}}}
o3 := bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$lt", tm1}}}}}}
o4 := bson.D{{"$match", bson.D{{"totalCount", bson.D{{"$lt", maxCount}, {"$gt", minCount}}}}}}
pipeline :=mongo.Pipeline{o1,o2,o3, o4}

data , err := collection.Aggregate(ctx, pipeline)
// fmt.Println(data)
if err != nil {
    panic(err)
}
var showsWithInfo []bson.M
if err = data.All(ctx, &showsWithInfo); err != nil {
    panic(err)
}

fmt.Fprintf(w, "%+v", showsWithInfo)  


// fmt.Println(showsWithInfo)
// fmt.Println(tm3)
// fmt.Println(tm)
// fmt.Println(minCount)
// fmt.Println(maxCount)

// var showsWithInfo []bson.M
// data.All(ctx, &showsWithInfo)
// fmt.Println(showsWithInfo)

// cursor ,err := collection.Find(context.Background(), bson.D{{}})
// fmt.Println("resuktsfound")
// var episodes []bson.M
// if err = cursor.All(ctx, &episodes); err != nil {
//     log.Fatal(err)
// }
// fmt.Fprintf(w, "Request: %+v", episodes)
// fmt.Println(episodes)

// showInfoCursor, err := collection.Aggregate(ctx, mongo.Pipeline{o2,o3})

// fmt.Println("Succesfully connected to the database...")
// if err != nil {
//     panic(err)
// }
// var showsWithInfo []bson.M
// if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
//     panic(err)
// }
// fmt.Println(showsWithInfo)





// pipe := collection.Pipe(operations)

// // Run the queries and capture the results
// results := []bson.M{}
// err1 := pipe.All(&results)

// if err1 != nil {
//     fmt.Printf("ERROR : %s\n", err1.Error())
//     return
// }

// fmt.Printf(w, results)



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