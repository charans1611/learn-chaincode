
package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	

var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)

var customerOrdersStr = "_customerorders"    // This will  be the key, value will be a list of orders placed by customer - wil be called by Customer


type MilkContainer struct{

        ContainerID string `json:"containerid"`
        User string        `json:"user"`

        Litres int       `json:"litres"`

}

type Order struct{
        OrderID string `json:"orderid"`
       User string `json:"user"`
       Status string `json:"status"`
       Litres int   `json:"litres"`
}

type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}


type Asset struct{
	  User string        `json:"user"`
	containerIDs []string `json:"containerids"`
	LitresofMilk int `json:"litresofmilk"`
	Supplycoins int `json:"supplycoins"`
}




func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	
	var err error
	
	fmt.Println("Welcome to the Supply chain management Phase 1, Deployment has been started")
 
       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }
	
        /* Reset container index list - Making sure the value corresponding to containerIndexStr  is empty */

       var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(containerIndexStr, jsonAsBytes)                     //Resetting - Making milk container list as empty 
       if err != nil {
		return nil, err
        }  
	
	
	/* Resetting the order list - Making sure the value corresponding to openOrdersStr is empty */
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	// Resetting the Assets of Supplier,Market, Logistics, Customer
	
	var emptyasset Asset
	
	
	jsonAsBytes, _ = json.Marshal(emptyasset)                // this is the byte format format of empty Asset structure
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // key -Supplier assets and value is empty now --> Supplier has no assets
	err = stub.PutState("MarketAssets", jsonAsBytes)         // key -Market assets and value is empty now --> Market has no assets
	err = stub.PutState("LogisticsAssets", jsonAsBytes)      // key - Logistics assets and value is empty now --> Logistic has no assets
	err = stub.PutState("CustomerAssets", jsonAsBytes)      // key - Customer assets and value is empty now --> Customer has no assets
	
	fmt.Println("Successfully deployed the code and orders and assets are reset")
	
	 return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		
		return t.Create_milkcontainer(stub, args)
		
	}else if function == "Create_coins" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coins(stub, args)	
        }else if function == "BuyMilkfromRetailer" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.BuyMilkfromRetailer(stub, args)	
        }else if function == "View_orderbyMarket" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.View_orderbyMarket(stub, args)	
        }else if function == "shiptocustomer" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.shiptocustomer(stub, args)	
        }else if function == "Order_milktoSupplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Order_milktoSupplier(stub, args)	
        }else if function == " View_orderbySupplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t. View_orderbySupplier(stub, args)	
        }else if function == "shiptologistics" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.shiptologistics(stub, args)	
        }
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}






func (t *SimpleChaincode) Create_milkcontainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error

// "1x22" "supplier" 20 
// args[0] args[1] args[2] 
	if len(args) != 3{
		return nil, errors.New("Please enter all the details")
	}
	fmt.Println("Creating milkcontainer asset")
id := args[0]
user := args[1]
litres,err:=strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Litres argument must be a numeric string")
	}
	
// Checking if the container already exists in the network
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details of given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

//If not present, create it and Update ledger, containerIndexStr, Assets of Supplier
//Creation
res.ContainerID = id
res.User = user
res.Litres = litres
milkAsBytes, _ =json.Marshal(res)

stub.PutState(res.ContainerID,milkAsBytes)
	
	 fmt.Println("Container created successfully, details are")
	fmt.Printf("%+v\n", res)
	
//Update containerIndexStr	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string                                        //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(containerAsBytes, &containerIndex)	
	
	
	containerIndex = append(containerIndex, res.ContainerID)          //append the newly created container to the global container list									//add marble name to index list
	fmt.Println("! container index: ", containerIndex)
	jsonAsBytes, _ := json.Marshal(containerIndex)
        err = stub.PutState(containerIndexStr, jsonAsBytes)
	
// append the container ID to the existing assets of the Supplier
	
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	supplierasset.User = "Supplier"
	supplierasset.containerIDs = append(supplierasset.containerIDs, res.ContainerID)
	supplierasset.LitresofMilk += res.Litres
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	fmt.Println("Balance of Supplier")
       fmt.Printf("%+v\n", supplierasset)

	return nil,nil

}

func (t *SimpleChaincode) Create_coins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"Market/Logistics",                          "100"
//args[0]                                     args[1]
//targeted owner                         No of supplycoins     

	user:= args[0]
	userAssets := user +"Assets"
        assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
	asset.User = user
	asset.Supplycoins,_ = strconv.Atoi(args[1])
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	fmt.Println("Balance of " , user)
        fmt.Printf("%+v\n", asset)


return nil,nil
}



func (t *SimpleChaincode) BuyMilkfromRetailer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//args[0]      args[1]
//"cus123"       "10"
// customer asks for a qty, check if market has that much quantity, if there-create a container for customer with qty he asked, and subtract the same from Market
	fmt.Println("Hello customer, welcome ")

	
	Openorder := Order{}
        Openorder.User = "customer"
        Openorder.Status = "Order received by Market"
        Openorder.OrderID = args[0]
        Openorder.Litres, _ = strconv.Atoi(args[1])
	fmt.Println("Hello customer, your order has been generated successfully, you can track it with id in the following details")
	fmt.Println("%+v\n",,Openorder)
        orderAsBytes,_ := json.Marshal(Openorder)
	stub.PutState(Openorder.OrderID,orderAsBytes)
	
	customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended",  Openorder.OrderID,"to existing customer orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(customerOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}
	
	
	//t.View_orderbyMarket(stub)
	
	// I think there should be a break here. Market should login go to view orders this should show what all orders he have
	// Then he should select a order that should trigger everything
	
	return nil,nil
}

func(t *SimpleChaincode)  View_orderbyMarket(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET- think of UI-View orders- does he pass any parameter there...
// so here also no need of any arguments.
	
	fmt.Printf("Hello Market, these are the orders placed to  you by customer")
	
	
	ordersAsBytes, _ := stub.GetState(customerOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}

func (t *SimpleChaincode)   shiptocustomer(stub shim.ChaincodeStubInterface, args[]string) ([]byte, error)
{
	// In UI, beside each order one button to ship to customer
	// we will extract details of orderId
	//we will exract asset balance of Market
	// If balance is der direct delivery to customer, otherwise order  to supplier 
//Fetch order details
	OrderID := args[0]
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	quantity := ShipOrder.Litres
//fetching assets of market	
	marketassetAsBytes, _ := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset )
	
//checking if market has the stock	
	if (Marketasset.LitresofMilk >= quantity ){
		fmt.Println("Enough stock is available, delivering to customer in short time")
		
		delivertocustomer(stub,OrderID)
	}else{
	        fmt.Println("Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer")
		
	        ShipOrder.Status = "In transit" // No matter, where the order placed by market is , for customer we will show it is "in transit"
	        orderAsBytes,err = json.Marshal(ShipOrder)
                stub.PutState(OrderID,orderAsBytes)  
		
		//Now we should send details of updated order status to customer, should be done in UI
		
                //Giving order to Supplier/Manufacturer
		//a,c := Order_milktoSupplier(stub,"20")
		//fmt.Println(a)
		//fmt.Println(c)
		
		
		//delivertocustomer(stub,b)
	}
	
	
	return nil,nil

}


func delivertocustomer(stub shim.ChaincodeStubInterface ,args string) (error){

	//args[0] 
	//OrderID  
	OrderID := args
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	quantity := ShipOrder.Litres
	
	
	
        customerassetAsBytes,_ := stub.GetState("CustomerAssets")        // The same key which we used in Init function 
	customerasset := Asset{}
	json.Unmarshal( customerassetAsBytes, &customerasset)
	
           marketassetAsBytes, err := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset )
	
	customerasset.LitresofMilk += quantity
	Marketasset.LitresofMilk -= quantity
	
	customerassetAsBytes,_ = json.Marshal(customerasset)
	stub.PutState("CustomerAssets",customerassetAsBytes)
	
	marketassetAsBytes,_ = json.Marshal(Marketasset)
	stub.PutState("MarketAssets",marketassetAsBytes)
	
	ShipOrder.Status ="Delivered to Customer"
	fmt.Println("%v\n", ShipOrder)
	orderAsBytes,err = json.Marshal(ShipOrder)
        stub.PutState(OrderID,orderAsBytes)
	
	
	b := []stub{"30", "Customer", "Market"}
	transfer(stub,b)        //Transfer should be automated. So it can't be invoked from UI..Loop hole
	fmt.Println("FINALLLLLYYYY, END OF THE STORY")
         
    return nil

}


func(t *SimpleChaincode) Order_milktoSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//"abcd" "20"
//OrderID litres
var err error
Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "Order placed to Supplier "
Openorder.OrderID = args[0]
Openorder.Litres,err = strconv.Atoi(args[1])
orderAsBytes,_ := json.Marshal(Openorder)
	
err = stub.PutState(Openorder.OrderID,orderAsBytes)
	
	 fmt.Println("your Order has been generated successfully")
	fmt.Printf("%+v\n", Openorder)
	
if err != nil {
		return nil, err
}

//Add the new order to the orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended ",Openorder.OrderID,"to existing market orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}
	 //View_orderbySupplier(stub)
	
return nil,nil
}



func(t *SimpleChaincode)  View_orderbySupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// This will be invoked by Supplier in UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty-but just for syntax-pass something as parameter in angular js
      
	
/* fetching the Orders*/
	fmt.Printf("Hello Supplier, here are the orders placed to you by Market")
	
	
	ordersAsBytes, _ := stub.GetState(openOrdersStr)
	
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	fmt.Println(orders)
	return nil,nil
}

func (t *SimpleChaincode)  ShiptoLogistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
//args[0]
//OrderID
	
//fetching order details
        OrderId := args[0]
        orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetching supplier asset details
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)

//checking if the supplier has stock or not
	
if ( supplierasset.LitresofMilk >= OrderID.Litres ){
	
	
	for i,val := range supplierasset.ContainerIDs{
		
        // If ordered quantity and container quantity , then proceed and trigger logistics
		containerassetAsBytes, err := stub.GetState(supplierasset.ContainerIDs[i])
		res := MilkContainer{} 
		json.Unmarshal(containerassetAsBytes,&res)
	
	     if (res.Litres == supplierasset.ContainerIDs[i].Litres) {
		fmt.Println("Found a suitable container and about to ship to retailer/market")
		ShipOrder.Status = "Ready to be Shipped to Market"   // send the updated status to market
	        orderAsBytes,err = json.Marshal(ShipOrder)
                stub.PutState(OrderID,orderAsBytes)
		
	        a := []string{ShipOrder.OrderID,res.ContainerID}
		init_logistics(stub,a)
	        return nil,nil	
	        }
	}
		
} else {
	fmt.Println("Supplier doesn't have sufficient stock")
	return nil, nil
}


}

func init_logistics(stub shim.ChaincodeStubInterface, args []string) ( error) {
	
	
	
	//args[0] args[1]
	// OrderId, ContainerID
	fmt.Println("we are moving, item is on the way to market/retailer")
	fmt.Println("Right now with logistics")
	OrderID := args[0]
	//ContainerID := args[1]
	
	// fetch the order details and update status as "in transit"
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ShipOrder.Status = "In transit to market"
	 
	orderAsBytes,err = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,orderAsBytes)
	
	fmt.Printf("%+v\n", ShipOrder)
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	orders.OpenOrders[0].Status = "In transit to market"
	ordersAsBytes,_ = json.Marshal(orders)
	stub.PutState(openOrdersStr,ordersAsBytes)
	set_user(stub,args)
	
	
	
	
	
return nil
}


func  set_user(stub shim.ChaincodeStubInterface, args []string) ( error) {
	
// OrderId  ContainerID
//args[0] args[1]
	
//So here we will set the user name in container ID to the one in Order ID and Status to Delivered - Asset Transfer
	fmt.Println("Transferring asset owner ship")
	OrderID := args[0]
	ContainerID := args[1]
//fetch order details
       orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details	
	assetAsBytes,err := stub.GetState(ContainerID)
	container := MilkContainer{}
	json.Unmarshal(assetAsBytes, &container)

	if (container.User == "Supplier"){
	
	container.User = ShipOrder.User             //ASSET TRANSFER
	fmt.Printf("%+v\n", container)
	fmt.Println("pushing the updated container back to ledger")
	assetAsBytes,err = json.Marshal(container)
	stub.PutState(ContainerID, assetAsBytes)    //Pushing the updated container  back to the ledger
	
	fmt.Println("Updating Supplier assets..")
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
		
	userAssets := container.User +"Assets"
	fmt.Println("Updating ",userAssets)
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
		
	asset.LitresofMilk += container.Litres
	supplierasset.LitresofMilk -= container.Litres
		
	
		
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	// update the Order and push back to ledger
	ShipOrder.Status = "Delivered to market"
	orderAsBytes,err = json.Marshal(ShipOrder) 
	stub.PutState(OrderID,orderAsBytes)      
	fmt.Printf("%+v\n", ShipOrder)
	//Updating the orders list 
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	orders.OpenOrders[0].Status = ShipOrder.Status
	ordersAsBytes,_ = json.Marshal(orders)
        stub.PutState(openOrdersStr,ordersAsBytes)
	//check the product before transferring money 
	checktheproduct(stub,args)
		
	}else
        {
                stub.PutState("setuser",[]byte("failure in this function"))
                //t.read(stub,"setuser")
                return nil
        }


return nil
}


func  checktheproduct(stub shim.ChaincodeStubInterface, args []string) ( error) {

// args[0] args[1]
// OrderID, ContainerID
	fmt.Println("Let us check the product")
	OrderID := args[0]
	ContainerID := args[1]
//fetch order details
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details
       assetAsBytes,_ := stub.GetState(ContainerID)
	Deliveredcontainer := MilkContainer{}
	json.Unmarshal(assetAsBytes, &Deliveredcontainer)

//check and transfer coins
	if (Deliveredcontainer.User == "Market" && Deliveredcontainer.Litres == ShipOrder.Litres) {
		
		fmt.Println("Thanks, I got  the right product, transferring amount to Supplier/Manufacturer")
		stub.PutState("Market Response",[]byte("Product received"))
		var b [3]string
		b[0]= "50"
		b[1] = "Market"
		b[2] = "Supplier"
		
		err = transfer(stub,b)
		if err!=nil{
			return err
		}
	        b[0]= "25"
		b[1] = "Supplier"
		b[2] = "Logistics"
		err = transfer(stub,b)
		if err!=nil{
			return err
		}
		return nil
       }else{
                stub.PutState("checktheproduct",[]byte("failure"))
		fmt.Println("I didn't get the right product")
               // t.read(stub,"checktheproduct")
                return nil
        }

	
return nil


}


func transfer( stub shim.ChaincodeStubInterface, args [3]string) ( error) {
	
//args[0]             args[1]         args[2]
//No of supplycoin      Sender         Reciever   
	//lets keep it simple for now, just fetch the coin from ledger, change username to Supplier and End of Story
	transferamount,_ := strconv.Atoi(args[0])
	sender := args[1]                               // this thing should be given by us in UI background
	receiver := args[2]                            // this will be given by the user on web page
	
	fmt.Println( sender, "transferring", transferamount, "coins to", receiver)
	
        senderAssets := sender +"Assets"
        senderassetAsBytes,_ := stub.GetState(senderAssets)        // The same key which we used in Init function 
	senderasset := Asset{}
	json.Unmarshal( senderassetAsBytes, &senderasset)
	
	
	receiverAssets := receiver+"Assets"
        receiverassetAsBytes,_ := stub.GetState(receiverAssets)        // The same key which we used in Init function 
	receiverasset := Asset{}
	json.Unmarshal( receiverassetAsBytes, &receiverasset)
	
	if ( senderasset.Supplycoins >= transferamount){
		
	senderasset.Supplycoins -= transferamount
	receiverasset.Supplycoins += transferamount
	
		  senderassetAsBytes,_=  json.Marshal(	senderasset)
	stub.PutState(senderAssets,  senderassetAsBytes)
	fmt.Println("Balance of " , sender)
       fmt.Printf("%+v\n", senderasset)
		
		receiverassetAsBytes,_=  json.Marshal(receiverasset)
	stub.PutState( receiverAssets,receiverassetAsBytes)
	fmt.Println("Balance of " , receiver)
       fmt.Printf("%+v\n", receiverasset)
		return nil
	}else {
		str := "Failed to transfer amount from" + sender + "to" + receiver
		
		return  errors.New(str)
	}
	
	
return nil
	
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
