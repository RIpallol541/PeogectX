const dgram = require("dgram");
const readline = require("readline");

const client = dgram.createSocket("udp4");
const serverPort = 8081;
const serverAddress = "127.0.0.1";

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
});

// Helper function to send a UDP request
function sendUDPRequest(data, callback) {
    const message = Buffer.from(JSON.stringify(data));
    client.send(message, serverPort, serverAddress, (err) => {
        if (err) {
            console.error("Error sending UDP message:", err);
            callback({ error: "Failed to send message" });
        } else {
            console.log("Request sent:", data);
        }
    });

    client.on("message", (msg) => {
        const response = JSON.parse(msg.toString());
        callback(response);
    });
}

// Main menu
function showMenu() {
    console.log("\nChoose an option:");
    console.log("1. Authenticate");
    console.log("2. Start matrix multiplication task");
    console.log("3. Exit");
    rl.question("Enter your choice: ", (choice) => {
        switch (choice.trim()) {
            case "1":
                authenticateUser();
                break;
            case "2":
                startMatrixTask();
                break;
            case "3":
                console.log("Goodbye!");
                client.close();
                rl.close();
                break;
            default:
                console.log("Invalid choice. Try again.");
                showMenu();
        }
    });
}

// Authentication
function authenticateUser() {
    rl.question("Enter username: ", (username) => {
        rl.question("Enter password: ", (password) => {
            const requestData = {
                action: "auth",
                username,
                password,
            };
            sendUDPRequest(requestData, (response) => {
                console.log("Server response:", response);
                showMenu();
            });
        });
    });
}

// Start matrix multiplication task
function startMatrixTask() {
    rl.question("Enter sizeMatrix (default 10): ", (sizeMatrixInput) => {
        rl.question("Enter maxDimension (default 9): ", (maxDimensionInput) => {
            const sizeMatrix = parseInt(sizeMatrixInput) || 10;
            const maxDimension = parseInt(maxDimensionInput) || 9;

            const requestData = {
                action: "matrix_task",
                sizeMatrix,
                maxDimension,
            };

            console.log("Starting matrix task...");
            sendUDPRequest(requestData, (response) => {
                if (Array.isArray(response)) {
                    response.forEach((result, index) => {
                        console.log(`Experiment #${index + 1}:`);
                        console.log(`  Lambda: ${result.Lambda}`);
                        console.log(`  Mu: ${result.Mu}`);
                        console.log(`  Dimension: ${result.Dimension}`);
                        console.log(`  Duration: ${result.Duration}`);
                        console.log(`  Result: ${result.Result}`);
                    });
                } else {
                    console.log("Server response:", response);
                }
                showMenu();
            });
        });
    });
}

// Start application
console.log("UDP Client Started.");
showMenu();
