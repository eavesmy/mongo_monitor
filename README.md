# mongodb oplog watcher

# Install
```go install github.com/eavesmy/mongo_monitor```

# Usage
1. start service.
2. run below scripts.
     
```javascript
const net = require("net");
                
class Watch {
                
    static initPath = "/tmp/mongo_watch"
                
    static connect(){
        let req = net.connect(this.initPath,() => {
            // bson.D{{"operationType", "update"}, {"updateDescription.updatedFields.cash", bson.D{{"$exists", true}}}}
            let match = { "operationType": "update", "updateDescription.updatedFields.cash": {$exists: true} }
            let db = "test"
            let collection = "a" 
                
            let data = { match, db, collection }
                
            req.write(JSON.stringify(data) + "\n");
        });  
                
        req.on("data",data => {
            console.log("receive: ",data.toString());
        });  
    }           
}               
                
module.exports = Watch;
                
                
Watch.connect()
```
