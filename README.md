
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
    docker run -p 3000:3000 my-fetch-receipt-app

