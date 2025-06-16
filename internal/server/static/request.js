async function fetchData(method, pathName, trackIdx, currentTime) {
    // The URL to which you want to send the JSON data
    const url = "/progress/"; // Replace with your server endpoint

    // The JSON data you want to send
    const dataToSend = {
        method: method,
        pathName: pathName,
        trackIdx: trackIdx,
        currentTime: currentTime,
    };

    // Convert the JavaScript object to a JSON string
    const jsonData = JSON.stringify(dataToSend);
    try {
        const response = await fetch(url, {
            method: "POST", // Specify the HTTP method (e.g., POST, PUT)
            headers: {
                "Content-Type": "application/json", // Set the Content-Type header to indicate JSON data
            },
            body: jsonData, // The JSON string to send in the request body
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json(); // Or response.text(), etc.
        console.log("Data fetched (async/await):", data);
    } catch (error) {
        console.error("Error fetching data (async/await):", error);
    }
}

function sendRequest(method, pathName, trackIdx, currentTime, callback) {
    // The URL to which you want to send the JSON data
    const url = "/progress/"; // Replace with your server endpoint

    // The JSON data you want to send
    const dataToSend = {
        method: method,
        pathName: pathName,
        trackIdx: trackIdx,
        currentTime: currentTime,
    };

    // Convert the JavaScript object to a JSON string
    const jsonData = JSON.stringify(dataToSend);

    // Use the fetch API to send the request
    fetch(url, {
        method: "POST", // Specify the HTTP method (e.g., POST, PUT)
        headers: {
            "Content-Type": "application/json", // Set the Content-Type header to indicate JSON data
        },
        body: jsonData, // The JSON string to send in the request body
    })
        .then((response) => {
            // Check if the request was successful (status code 2xx)
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            // Parse the response body as JSON
            return response.json();
        })
        .then((data) => {
            // Process the JSON response data
            console.log("Success:", data);
            callback(data);
        })
        .catch((error) => {
            // Handle any errors that occurred during the request
            console.error("Error:", error);
        });
}
