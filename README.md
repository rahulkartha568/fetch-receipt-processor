
# Fetch Receipt Processor

This repository contains the source code for a simple webservice based off the assesment description: https://github.com/fetch-rewards/receipt-processor-challenge. You can run this web service locally using Docker. Follow the steps below to get started.


## Usage

1. Clone this repository to your local machine:
   git clone https://github.com/rahulkartha568/fetch-receipt-processor.git

2. Navigate to the project directory:
    cd fetch-receipt-processor
3. Build the Docker image:
    docker build -t fetch-receipt-processor .
4. Run the docker container:
    docker run -p 3000:3000 fetch-receipt-processor
5. Access the web service in your web browser or through API calls:
Web service URL: http://localhost:3000

## API Endpoints
Endpoint: Process Receipts 
Path: /receipts/process
Method: POST
Payload: Receipt JSON
Response: JSON containing an id for the receipt.

Example Request & Response:
   Request URL: http://localhost:3000/receipts/process
   POST Body: 
   {
    "retailer": "Walgreens",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "2.65",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
             ]
   }
   
   Example Response:
   {
    "id": "f10819a6-5b59-4600-b16e-33902eb11662"
   }
   
Endpoint: Get Points
Path: /receipts/{id}/points
Method: GET
Response: A JSON object containing the number of points awarded.

Example Request & Response:
   GET Request URL: http://localhost:3000/receipts/f10819a6-5b59-4600-b16e-33902e11662/points

   
   Example Response:
   {
    "points": 15
   }
   
