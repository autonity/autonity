// Define the contract address and the event signature
var contractAddress = ""; // Replace with your contract address

// Compute the event signature hash (example: Transfer event)
var eventSignature = web3.sha3("Transfer(address,address,uint256)"); // Replace with your event signature
console.log("Event Signature Hash:", eventSignature);

// Define the log query filter
var filter = {
    fromBlock: "105960", // Start block (use "latest" for real-time queries)
    toBlock: "latest", // End block
    address: contractAddress, // Contract address to filter logs
    topics: [eventSignature] // Filter by the event signature hash
};

// Fetch the logs
var logs = eth.getLogs(filter);

// Decode and display logs
logs.forEach(function (log) {
    console.log("Raw Log:", log);

    // Decode the indexed parameters (e.g., `from` and `to` in Transfer event)
    var from = "0x" + log.topics[1].slice(26); // Indexed parameter 1 (address from)
    var to = "0x" + log.topics[2].slice(26); // Indexed parameter 2 (address to)

    // Decode the non-indexed parameter (e.g., `value` in Transfer event)
    var value = web3.toBigNumber(log.data).toString(); // Non-indexed parameter (uint256 value)

    // Display the decoded data
    console.log("Decoded Log:");
    console.log("  From:", from);
    console.log("  To:", to);
    console.log("  Value:", value);
});
