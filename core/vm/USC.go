 package vm

 import (
	 "github.com/ethereum/go-ethereum/common"
	 "github.com/ethereum/go-ethereum/common/math"
	 "math/big"
	 "errors"
	 "encoding/hex"
	 "fmt"
 )
 /* This contract has no utility except saving all datas in its storage */
 var (
 	 // DataContract is the address of the contract that will save all USC data in its storage
	 DataContract = common.HexToAddress("0x1000")

	 //StotalSupply is the total amount in backed currency that is in the system
	 StotalSupply = common.BigToHash(new(big.Int).SetInt64(0))
	 //StotalClients is the number of Clients in the system
	 StotalClients = common.BigToHash(new(big.Int).SetInt64(1))
	 //SlastTokenId is the last token id created
	 SlastTokenId = new(big.Int).SetInt64(2)

	 //SbaseClientAddress is the base slot of an array containing Client addresses
	 SbaseClientAddress = new(big.Int).SetInt64(100)

	 /*
	 SbaseClientBalance is the base slot of an array containg the Client's tokens
	 SbaseClientBalance + Client_id * ClientsOffset => total balance
	 SbaseClientBalance + Client_id * ClientsOffset + 1 => token number
	 SbaseClientBalance + Client_id * ClientsOffset + x=> token id
	 */
	 SbaseClientBalance = new(big.Int).SetInt64(1000000)
	 /*
	 ClientsOffset is the number of slots between two Clients in the ClientBalance array
	 */
	 ClientsOffset = math.BigPow(10,14)// maximum number of tokens
 )

 var (
	 InvalidInputSize = errors.New("Invalid Input Size")
	 NotEnoughFunds = errors.New("Not Enough Funds")
	 UserNotInTheSystem = errors.New("User Not In The System")
 )

 var (
 	bigOne  = new(big.Int).SetUint64(1)
 	bigTwo  = new(big.Int).SetUint64(2)
 )

 func ReturnOrCreateParticipant(participant common.Address, evm *EVM) uint64 {
 	totalClients := evm.StateDB.GetState(DataContract, StotalClients).Big().Uint64()
 	found := false
 	i := uint64(0)
 	currentSlot := new(big.Int).Set(SbaseClientAddress)
 	for ; i < totalClients && !found; i+=1 {
		currentClient := common.BigToAddress(evm.StateDB.GetState(DataContract, common.BigToHash(currentSlot)).Big())
		found = currentClient == participant
		currentSlot.Add(currentSlot, bigOne)
		fmt.Println("current slot:")
		fmt.Println(hex.Dump(currentSlot.Bytes()))
	}
	if found {
		return i - 1
	}

	evm.StateDB.SetState(DataContract, common.BigToHash(currentSlot), participant.Hash())
 	newTotalClients := new(big.Int).Add(new(big.Int).SetUint64(totalClients), bigOne)
	evm.StateDB.SetState(DataContract, StotalClients, common.BigToHash(newTotalClients))

 	return i
 }

 type USC_Data struct{}
 func (c *USC_Data) RequiredGas(input []byte) uint64 {
 	return 0
 }

 func (c *USC_Data) Run(input []byte, evm *EVM, contract *Contract) ([]byte, error) {
 	//0: total amount
 	//1: total Clients
 	//[100, 100 + total Clients - 1] : Client's ethereum account address
 	//[1000
 	return nil, nil
 }

 /* This contract allow the central Client to issue new USC tokens */
 type USC_Fund struct{}
 type USC_Fund_args struct{
 	recipent common.Address
 	amount big.Int
 }

 func getBigIntState(evm *EVM, address *big.Int) *big.Int {
 	value := evm.StateDB.GetState(DataContract, common.BigToHash(address))
 	return value.Big()
 }

 func setBigIntState(evm *EVM, address *big.Int, value *big.Int){
	 evm.StateDB.SetState(DataContract, common.BigToHash(address), common.BigToHash(value))
 }

 func (c *USC_Fund) RequiredGas(input []byte) uint64 {
 	return 0
 }

 func (c *USC_Fund) ValidateInput(input []byte) (bool, USC_Fund_args){
 	inputSize := 20 + 32
 	args := USC_Fund_args{}
 	if len(input) != inputSize{
 		return false, args
	}
	args.recipent.SetBytes(input[0:20])
 	args.amount.SetBytes(input[20:52])
 	return true, args
 }


 func (c *USC_Fund) Run(input []byte, evm *EVM, contract *Contract) ([]byte, error) {

	fmt.Println("[DEBUG]USC FUND INVOKED")
	fmt.Println(hex.Dump(input))
	valid, args := c.ValidateInput(input)
	if !valid {
		fmt.Println("[DEBUG]Invalid Input Size")
		return nil,InvalidInputSize
	}

	fmt.Println("Participant:")
	fmt.Println(hex.Dump(args.recipent.Bytes()))

	clientId := ReturnOrCreateParticipant(args.recipent, evm)

	clientSlot := new(big.Int).Mul(ClientsOffset, new(big.Int).SetUint64(clientId))
	clientSlot.Add(clientSlot, SbaseClientBalance)

	tokenCounter := new(big.Int).SetUint64(0)
	clientTokens := getBigIntState(evm, new(big.Int).Add(clientSlot, bigOne))
	currentSlot := new(big.Int).Add(clientSlot,new(big.Int).Add(clientTokens,bigTwo))
	lastTokenId := getBigIntState(evm, SlastTokenId)
	for ; tokenCounter.Cmp(&args.amount) != 0; tokenCounter.Add(tokenCounter, bigOne){
		lastTokenId.Add(lastTokenId, bigOne)
		newToken := common.Hash{}.SetDenomId(8, lastTokenId)
		evm.StateDB.SetState(DataContract, common.BigToHash(currentSlot), newToken)
		currentSlot.Add(currentSlot, bigOne)
	}

	newClientBalance := new(big.Int).Add(getBigIntState(evm, clientSlot), tokenCounter)
	setBigIntState(evm, clientSlot, newClientBalance)

	newTokenNumber := new(big.Int).Add(clientTokens, tokenCounter)
	setBigIntState(evm, new(big.Int).Add(clientSlot, bigOne), newTokenNumber)

	newTotalSupply := new(big.Int).Add(getBigIntState(evm, StotalSupply.Big()), tokenCounter)
	setBigIntState(evm, StotalSupply.Big(), newTotalSupply)

	setBigIntState(evm, SlastTokenId, lastTokenId)

	return nil, nil
 }

 /*
 		USC Transfer Function
 */
 type USC_Transfer struct{}
 type USC_Transfer_args struct{
	 recipent common.Address
	 amount big.Int
 }
 func (c *USC_Transfer) RequiredGas(input []byte) uint64 {
	return 0
 }

 func (c *USC_Transfer) ValidateInput(input []byte) (bool, USC_Transfer_args){
	 inputSize := 20 + 32
	 args := USC_Transfer_args{}
	 if len(input) != inputSize{
		 return false, args
	 }
	 args.recipent.SetBytes(input[0:20])
	 args.amount.SetBytes(input[20:52])
	 return true, args
 }
 func (c *USC_Transfer) Run(input []byte, evm *EVM, contract *Contract) ([]byte, error) {

	 fmt.Println("[DEBUG]USC TRANSFER INVOKED")
	 fmt.Println(hex.Dump(input))
	 valid, args := c.ValidateInput(input)
	 if !valid {
		 fmt.Println("[DEBUG]Invalid Input Size")
		 return nil,InvalidInputSize
	 }

	 fmt.Println("Recipient:")
	 fmt.Println(hex.Dump(args.recipent.Bytes()))
	 fmt.Println("Sender:")
	 fmt.Println(hex.Dump(contract.caller.Address().Bytes()))

	 sender := contract.caller.Address()
	 senderId := ReturnOrCreateParticipant(sender, evm)

	 senderSlot := new(big.Int).Mul(ClientsOffset, new(big.Int).SetUint64(senderId))
	 senderSlot.Add(senderSlot, SbaseClientBalance)
	 senderTotalMoney := getBigIntState(evm, senderSlot)

	 if args.amount.Cmp(senderTotalMoney) == 1 {
		 fmt.Println("[DEBUG] Not Enough Funds")
		 return nil,NotEnoughFunds
	 }

	 receiverId := ReturnOrCreateParticipant(args.recipent, evm)
	 receiverSlot := new(big.Int).Mul(ClientsOffset, new(big.Int).SetUint64(receiverId))
	 receiverSlot.Add(receiverSlot, SbaseClientBalance)

	 senderTotalTokensSlot := new(big.Int).Add(senderSlot,bigOne)
	 senderTotalTokens := getBigIntState(evm, senderTotalTokensSlot)
	 senderLastTokenSlot := new(big.Int).Add(senderSlot,senderTotalTokens)
	 senderLastTokenSlot.Add(senderLastTokenSlot,bigOne)

	 receiverTotalTokensSlot := new(big.Int).Add(receiverSlot,bigOne)
	 receiverTotalTokens := getBigIntState(evm, receiverTotalTokensSlot)
	 receiverTokenFreeSlot := new(big.Int).Add(senderSlot,receiverTotalTokens)
	 receiverTokenFreeSlot.Add(receiverTokenFreeSlot,bigTwo)

	 //tokenCounter := new(big.Int).SetUint64(0)
	 /*
	 for ; tokenCounter.Cmp(&args.amount) != 0; tokenCounter.Add(tokenCounter, bigOne){
		 lastTokenId.Add(lastTokenId, bigOne)
		 newToken := common.Hash{}.SetDenomId(8, lastTokenId)
		 evm.StateDB.SetState(DataContract, common.BigToHash(currentSlot), newToken)
		 currentSlot.Add(currentSlot, bigOne)
	 }

	 comment
	 clientTokens := getBigIntState(evm, new(big.Int).Add(clientSlot, bigOne))
	 currentSlot := new(big.Int).Add(clientSlot,new(big.Int).Add(clientTokens,bigTwo))
	 lastTokenId := getBigIntState(evm, SlastTokenId)


	 newClientBalance := new(big.Int).Add(getBigIntState(evm, clientSlot), tokenCounter)
	 setBigIntState(evm, clientSlot, newClientBalance)

	 newTokenNumber := new(big.Int).Add(clientTokens, tokenCounter)
	 setBigIntState(evm, new(big.Int).Add(clientSlot, bigOne), newTokenNumber)

	 newTotalSupply := new(big.Int).Add(getBigIntState(evm, StotalSupply.Big()), tokenCounter)
	 setBigIntState(evm, StotalSupply.Big(), newTotalSupply)

	 setBigIntState(evm, SlastTokenId, lastTokenId)
	*/
	 return nil, nil
 }

type USC_Defund struct{}
type USC_Defund_args struct{
	 recipent common.Address
	 amount big.Int
}

func (c *USC_Defund) RequiredGas(input []byte) uint64 {
	return 0
}

func (c *USC_Defund) ValidateInput(input []byte) (bool, USC_Defund_args){
	inputSize := 20 + 32
	args := USC_Defund_args{}
	if len(input) != inputSize{
		return false, args
	}
	args.recipent.SetBytes(input[0:20])
	args.amount.SetBytes(input[20:52])
	return true, args
}


func (c *USC_Defund) Run(input []byte, evm *EVM, contract *Contract) ([]byte, error) {

	//input validation
	fmt.Println("[DEBUG]USC DEFUND INVOKED")
	fmt.Println(hex.Dump(input))
	valid, args := c.ValidateInput(input)
	if !valid {
		fmt.Println("[DEBUG]Invalid Input Size")
		return nil,InvalidInputSize
	}

	fmt.Println("Participant:")
	fmt.Println(hex.Dump(args.recipent.Bytes()))

	//check user is in the system
	totalClients := evm.StateDB.GetState(DataContract, StotalClients).Big().Uint64()
	found := false
	i := uint64(0)
	currentSlot := new(big.Int).Set(SbaseClientAddress)
	for ; i < totalClients && !found; i+=1 {
		currentClient := common.BigToAddress(evm.StateDB.GetState(DataContract, common.BigToHash(currentSlot)).Big())
		found = currentClient == args.recipent
		currentSlot.Add(currentSlot, bigOne)
		fmt.Println("current slot:")
		fmt.Println(hex.Dump(currentSlot.Bytes()))
	}

	if !found {
		return nil, UserNotInTheSystem
	}

	fmt.Println("USER FOUND ...")
	clientId := i
	fmt.Printf("CLIENT ID: %d\n", clientId)
	clientSlot := new(big.Int).Mul(ClientsOffset, new(big.Int).SetUint64(clientId-1))
	clientSlot.Add(clientSlot, SbaseClientBalance)
	fmt.Println("CLIENT SLOT:"+clientSlot.String())
	clientTokens := getBigIntState(evm, new(big.Int).Add(clientSlot, bigOne))
	fmt.Println("CLIENT TOKENS:"+clientTokens.String())

	
	//check balance is enough
	if clientTokens.Cmp(&args.amount) == -1 {
		return nil, NotEnoughFunds
	}
	fmt.Printf("ENOUGH FUNDS\n")
	//remove tokens and update related variables
	tokenCounter := args.amount.Uint64()
	fmt.Println("Number of tokens to eliminate:",tokenCounter)
	currentSlot = new(big.Int).Add(clientSlot,new(big.Int).Add(clientTokens,bigTwo))// now I am on the last token
	fmt.Println("LAST TOKEN SLOT:"+currentSlot.String())

	zeroToken := common.Hash{}.SetDenomId(0, new(big.Int).SetUint64(0))
	//zeroToken:= make([]byte, 32)
	
	for ; tokenCounter > 0; tokenCounter-- {
		//fmt.Println("Token value before:")		
		//fmt.Println(hex.Dump(evm.StateDB.GetState(DataContract, common.BigToHash(currentSlot)).Bytes()))
		//evm.StateDB.SetState(DataContract, common.BigToHash(currentSlot), common.BytesToHash(zeroToken))
		evm.StateDB.SetState(DataContract, common.BigToHash(currentSlot), zeroToken)
		//fmt.Println("Token value after:")
		//fmt.Println(hex.Dump(evm.StateDB.GetState(DataContract, common.BigToHash(currentSlot)).Bytes()))
		currentSlot.Sub(currentSlot, bigOne)
	}

	newClientBalance := new(big.Int).Sub(clientTokens, &args.amount)
	setBigIntState(evm, clientSlot, newClientBalance)

	newTokenNumber := new(big.Int).Sub(clientTokens, &args.amount)
	setBigIntState(evm, new(big.Int).Add(clientSlot, bigOne), newTokenNumber)

	return nil, nil
	}
