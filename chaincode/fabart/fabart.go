package main

import (
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim" // import for Chaincode Interface
    pb "github.com/hyperledger/fabric/protos/peer"      // import for peer response
    "encoding/json"
    "unicode/utf8"
)

// Defined to implement chaincode interface
type Fabart struct {
}

// Define our struct to store art in Blockchain
type Art struct {
    Id string  // This one will be our key
    Price string   // The current price of the art
    Organisation string // The organisation through which this asset is managed
    OwnerName string // The name of the current owner
    IsActive bool // Whether this piece is active or not
}

// Initialise
func (c *Fabart) Init(stub shim.ChaincodeStubInterface) pb.Response {
    creatorByte, err := stub.GetCreator()
    if err != nil {
        return shim.Error("GetCreator err")
    }

    // Prevent errors if string is not UFT-8 encoded
    s := string(creatorByte);
    if !utf8.ValidString(s) {
        v := make([]rune, 0, len(s))
        for i, r := range s {
            if r == utf8.RuneError {
                _, size := utf8.DecodeRuneInString(s[i:])
                if size == 1 {
                    continue
                }
            }
            v = append(v, r)
        }
        s = string(v)
    }

    stub.PutState(s, []byte("producer"))
    return shim.Success(nil)
}

// Invoke 
func (c *Fabart) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    chaincodeAction, args := stub.GetFunctionAndParameters()

    switch chaincodeAction {
    case "register":
        return c.register(stub, args)
    case "transfer":
        return c.transfer(stub, args)
    case "queryDetails":
        return c.queryDetails(stub, args)
    default:
        return shim.Error("Not available")
    }
}

// Register a piece of art with the application
func (c *Fabart) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 4 {
        return shim.Error("An id, price and type must be provided to register an asset")
    }

    art := Art{args[0], args[1], args[2], args[3], true}
    encodedArt, err := json.Marshal(art)

    if err != nil {
        return shim.Error(err.Error())
    }

    err = stub.PutState(art.Id, encodedArt)

    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

// Query a piece of art for price and owner etc
func (c *Fabart) queryDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    value, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("ID number " + args[0] + " not found")
    }

    var art Art
    // Decode value
    json.Unmarshal(value, &art)

    fmt.Print(art)
    // Response info
    return shim.Success([]byte(" ID: " + art.Id + " Price: " + art.Price + " OwnerName: " + art.OwnerName))
}

// Transfer ownership of art
func (c *Fabart) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error("An id and an owner must be supplied")
    }

    v, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("ID " + args[0] + " not found ")
    }

    var art Art
    json.Unmarshal(v, &art)
    art.OwnerName = args[1]
    encodedArt, err := json.Marshal(art)

    err = stub.PutState(art.Id, encodedArt)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(nil)
}

func main() {
    err := shim.Start(new(Fabart))
    if err != nil {
        fmt.Printf("Error starting chaincode sample: %s", err)
    }
}
